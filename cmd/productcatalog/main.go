package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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

const serviceName = "productcatalogservice"
const servicePort = "8082"

type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Categories  []string `json:"categories"`
}

type ProductCatalog struct {
	logger   *logging.Logger
	tracer   trace.Tracer
	products map[string]Product
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

	catalog := &ProductCatalog{
		logger: logger,
		tracer: tracing.GetTracer(serviceName),
		products: map[string]Product{
			"OLJCESPC7Z": {ID: "OLJCESPC7Z", Name: "Vintage Typewriter", Description: "Classic mechanical typewriter", Price: 67.99, Categories: []string{"vintage"}},
			"66VCHSJNUP": {ID: "66VCHSJNUP", Name: "Vintage Camera Lens", Description: "Telephoto camera lens", Price: 12.49, Categories: []string{"photography", "vintage"}},
			"1YMWWN1N4O": {ID: "1YMWWN1N4O", Name: "Home Barista Kit", Description: "Complete coffee making set", Price: 124.99, Categories: []string{"kitchen"}},
			"L9ECAV7KIM": {ID: "L9ECAV7KIM", Name: "Terrarium", Description: "Glass terrarium for plants", Price: 36.45, Categories: []string{"gardening"}},
			"2ZYFJ3GM2N": {ID: "2ZYFJ3GM2N", Name: "Film Camera", Description: "Vintage 35mm film camera", Price: 2244.99, Categories: []string{"photography", "vintage"}},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/products", catalog.handleListProducts)
	mux.HandleFunc("/product/", catalog.handleGetProduct)
	mux.HandleFunc("/search", catalog.handleSearchProducts)
	mux.HandleFunc("/health", catalog.handleHealth)

	server := &http.Server{
		Addr:    ":" + servicePort,
		Handler: mux,
	}

	go func() {
		logger.Info("Product catalog service starting", zap.String("port", servicePort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down product catalog service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func (pc *ProductCatalog) handleListProducts(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), pc.tracer, "ListProducts")
	defer span.End()

	pc.logger.InfoContext(ctx, "List products requested")

	// Simulate database query
	time.Sleep(time.Duration(10+rand.Intn(30)) * time.Millisecond)

	products := make([]Product, 0, len(pc.products))
	for _, p := range pc.products {
		products = append(products, p)
	}

	tracing.AddSpanAttributes(span, attribute.Int("product.count", len(products)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (pc *ProductCatalog) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), pc.tracer, "GetProduct")
	defer span.End()

	productID := r.URL.Path[len("/product/"):]
	tracing.AddSpanAttributes(span, attribute.String("product.id", productID))

	pc.logger.InfoContext(ctx, "Get product", zap.String("product_id", productID))

	// Simulate database query
	time.Sleep(time.Duration(5+rand.Intn(20)) * time.Millisecond)

	product, exists := pc.products[productID]
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (pc *ProductCatalog) handleSearchProducts(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), pc.tracer, "SearchProducts")
	defer span.End()

	query := r.URL.Query().Get("q")
	tracing.AddSpanAttributes(span, attribute.String("search.query", query))

	pc.logger.InfoContext(ctx, "Search products", zap.String("query", query))

	// Simulate search
	time.Sleep(time.Duration(15+rand.Intn(35)) * time.Millisecond)

	results := []Product{}
	for _, p := range pc.products {
		// Simple search implementation
		results = append(results, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (pc *ProductCatalog) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

