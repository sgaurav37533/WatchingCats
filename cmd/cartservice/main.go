package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gaurav/otel-observability/internal/config"
	"github.com/gaurav/otel-observability/internal/logging"
	"github.com/gaurav/otel-observability/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const serviceName = "cartservice"
const servicePort = "8081"

type CartItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartService struct {
	logger *logging.Logger
	tracer trace.Tracer
	carts  map[string]*Cart
	mu     sync.RWMutex
}

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := logging.NewLogger(logging.LoggerConfig{
		Level:       cfg.Logging.Level,
		Format:      cfg.Logging.Format,
		ServiceName: serviceName,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	if cfg.Tracing.Enabled {
		tp, err := tracing.InitTracer(tracing.TracerConfig{
			ServiceName:  serviceName,
			Endpoint:     cfg.Tracing.Endpoint,
			Insecure:     cfg.Tracing.Insecure,
			SamplingRate: cfg.Tracing.SamplingRate,
		})
		if err != nil {
			logger.Error("Failed to initialize tracer", zap.Error(err))
		} else {
			defer tp.Shutdown(context.Background())
		}
	}

	cartService := &CartService{
		logger: logger,
		tracer: tracing.GetTracer(serviceName),
		carts:  make(map[string]*Cart),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cart/", cartService.handleGetCart)
	mux.HandleFunc("/cart/add", cartService.handleAddItem)
	mux.HandleFunc("/cart/remove", cartService.handleRemoveItem)
	mux.HandleFunc("/health", cartService.handleHealth)

	server := &http.Server{
		Addr:    ":" + servicePort,
		Handler: mux,
	}

	go func() {
		logger.Info("Cart service starting", zap.String("port", servicePort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down cart service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func (cs *CartService) handleGetCart(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), cs.tracer, "GetCart")
	defer span.End()

	userID := r.URL.Query().Get("user_id")
	tracing.AddSpanAttributes(span, attribute.String("user.id", userID))

	cs.logger.InfoContext(ctx, "Get cart", zap.String("user_id", userID))

	cs.mu.RLock()
	cart, exists := cs.carts[userID]
	cs.mu.RUnlock()

	if !exists {
		cart = &Cart{UserID: userID, Items: []CartItem{}}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (cs *CartService) handleAddItem(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), cs.tracer, "AddItem")
	defer span.End()

	var req struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tracing.AddSpanAttributes(span,
		attribute.String("user.id", req.UserID),
		attribute.String("product.id", req.ProductID),
		attribute.Int("quantity", req.Quantity),
	)

	cs.logger.InfoContext(ctx, "Add item to cart",
		zap.String("user_id", req.UserID),
		zap.String("product_id", req.ProductID),
		zap.Int("quantity", req.Quantity),
	)

	cs.mu.Lock()
	cart, exists := cs.carts[req.UserID]
	if !exists {
		cart = &Cart{UserID: req.UserID, Items: []CartItem{}}
		cs.carts[req.UserID] = cart
	}
	cart.Items = append(cart.Items, CartItem{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})
	cs.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (cs *CartService) handleRemoveItem(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), cs.tracer, "RemoveItem")
	defer span.End()

	userID := r.URL.Query().Get("user_id")
	productID := r.URL.Query().Get("product_id")

	cs.logger.InfoContext(ctx, "Remove item from cart",
		zap.String("user_id", userID),
		zap.String("product_id", productID),
	)

	cs.mu.Lock()
	if cart, exists := cs.carts[userID]; exists {
		newItems := []CartItem{}
		for _, item := range cart.Items {
			if item.ProductID != productID {
				newItems = append(newItems, item)
			}
		}
		cart.Items = newItems
	}
	cs.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func (cs *CartService) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

