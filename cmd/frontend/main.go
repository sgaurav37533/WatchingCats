package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gaurav/otel-observability/internal/config"
	"github.com/gaurav/otel-observability/internal/logging"
	"github.com/gaurav/otel-observability/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const serviceName = "frontend"
const servicePort = "8080"

type Frontend struct {
	logger *logging.Logger
	tracer trace.Tracer
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

	frontend := &Frontend{
		logger: logger,
		tracer: tracing.GetTracer(serviceName),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", frontend.handleHome)
	mux.HandleFunc("/product/", frontend.handleProduct)
	mux.HandleFunc("/cart", frontend.handleCart)
	mux.HandleFunc("/checkout", frontend.handleCheckout)
	mux.HandleFunc("/health", frontend.handleHealth)

	server := &http.Server{
		Addr:    ":" + servicePort,
		Handler: mux,
	}

	go func() {
		logger.Info("Frontend service starting", zap.String("port", servicePort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down frontend service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func (f *Frontend) handleHome(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), f.tracer, "GET /")
	defer span.End()

	tracing.AddSpanAttributes(span,
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.Path),
	)

	f.logger.InfoContext(ctx, "Home page requested")

	response := map[string]interface{}{
		"service": serviceName,
		"message": "Welcome to OpenTelemetry Demo",
		"endpoints": []string{
			"/product/{id}",
			"/cart",
			"/checkout",
			"/health",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (f *Frontend) handleProduct(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), f.tracer, "GET /product")
	defer span.End()

	productID := r.URL.Path[len("/product/"):]
	tracing.AddSpanAttributes(span,
		attribute.String("product.id", productID),
	)

	f.logger.InfoContext(ctx, "Product page requested", zap.String("product_id", productID))

	// Simulate calling product catalog service
	time.Sleep(50 * time.Millisecond)

	response := map[string]interface{}{
		"product_id": productID,
		"name":       "Sample Product",
		"price":      99.99,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (f *Frontend) handleCart(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), f.tracer, "GET /cart")
	defer span.End()

	f.logger.InfoContext(ctx, "Cart page requested")

	response := map[string]interface{}{
		"items": []map[string]interface{}{
			{"product_id": "123", "quantity": 2},
		},
		"total": 199.98,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (f *Frontend) handleCheckout(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), f.tracer, "POST /checkout")
	defer span.End()

	f.logger.InfoContext(ctx, "Checkout requested")

	// Simulate checkout process
	time.Sleep(100 * time.Millisecond)

	response := map[string]interface{}{
		"order_id": "ORD-12345",
		"status":   "confirmed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (f *Frontend) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

