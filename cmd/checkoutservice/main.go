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
	"github.com/gaurav/otel-observability/internal/exceptions"
	"github.com/gaurav/otel-observability/internal/logging"
	"github.com/gaurav/otel-observability/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const serviceName = "checkoutservice"
const servicePort = "8083"

type CheckoutService struct {
	logger           *logging.Logger
	tracer           trace.Tracer
	exceptionTracker *exceptions.Tracker
}

type CheckoutRequest struct {
	UserID    string  `json:"user_id"`
	Email     string  `json:"email"`
	Address   Address `json:"address"`
	CreditCard CreditCard `json:"credit_card"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

type CreditCard struct {
	Number string `json:"number"`
	CVV    string `json:"cvv"`
	Expiry string `json:"expiry"`
}

type CheckoutResponse struct {
	OrderID       string    `json:"order_id"`
	ShippingCost  float64   `json:"shipping_cost"`
	Total         float64   `json:"total"`
	TransactionID string    `json:"transaction_id"`
	Timestamp     time.Time `json:"timestamp"`
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

	exceptionTracker := exceptions.NewTracker(logger.Logger)

	checkoutSvc := &CheckoutService{
		logger:           logger,
		tracer:           tracing.GetTracer(serviceName),
		exceptionTracker: exceptionTracker,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/checkout", checkoutSvc.handleCheckout)
	mux.HandleFunc("/health", checkoutSvc.handleHealth)

	server := &http.Server{
		Addr:    ":" + servicePort,
		Handler: mux,
	}

	go func() {
		logger.Info("Checkout service starting", zap.String("port", servicePort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down checkout service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func (cs *CheckoutService) handleCheckout(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracing.StartSpan(r.Context(), cs.tracer, "PlaceOrder")
	defer span.End()

	var req CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tracing.AddSpanAttributes(span,
		attribute.String("user.id", req.UserID),
		attribute.String("user.email", req.Email),
	)

	cs.logger.InfoContext(ctx, "Processing checkout", zap.String("user_id", req.UserID))

	// Simulate payment processing with occasional failures
	if rand.Float64() < 0.1 { // 10% failure rate
		err := fmt.Errorf("payment processing failed")
		cs.exceptionTracker.RecordException(ctx, err, exceptions.Options{
			Severity: exceptions.SeverityError,
			Tags: map[string]string{
				"user_id": req.UserID,
				"stage":   "payment",
			},
		})
		tracing.RecordError(span, err)
		cs.logger.ErrorContext(ctx, "Checkout failed", zap.Error(err))
		http.Error(w, "Payment processing failed", http.StatusInternalServerError)
		return
	}

	// Simulate processing time
	time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)

	orderID := fmt.Sprintf("ORD-%d", time.Now().Unix())
	transactionID := fmt.Sprintf("TXN-%d", time.Now().UnixNano())

	response := CheckoutResponse{
		OrderID:       orderID,
		ShippingCost:  5.99,
		Total:         105.99,
		TransactionID: transactionID,
		Timestamp:     time.Now(),
	}

	tracing.AddSpanAttributes(span,
		attribute.String("order.id", orderID),
		attribute.String("transaction.id", transactionID),
	)

	cs.logger.InfoContext(ctx, "Checkout completed",
		zap.String("order_id", orderID),
		zap.String("user_id", req.UserID),
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (cs *CheckoutService) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

