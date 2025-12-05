# Code Examples

Practical examples for using the OpenTelemetry Observability Platform components.

## Table of Contents

1. [Distributed Tracing](#distributed-tracing)
2. [Structured Logging](#structured-logging)
3. [Exception Tracking](#exception-tracking)
4. [Alerting](#alerting)
5. [Full Integration](#full-integration)

---

## Distributed Tracing

### Basic Span Creation

```go
package main

import (
    "context"
    "github.com/gaurav/otel-observability/internal/tracing"
    "go.opentelemetry.io/otel/attribute"
)

func processOrder(ctx context.Context, orderID string) error {
    // Get tracer
    tracer := tracing.GetTracer("order-service")
    
    // Create span
    ctx, span := tracing.StartSpan(ctx, tracer, "process_order")
    defer span.End()
    
    // Add attributes
    tracing.AddSpanAttributes(span,
        attribute.String("order.id", orderID),
        attribute.String("order.status", "processing"),
    )
    
    // Your business logic here...
    
    return nil
}
```

### Nested Spans (Parent-Child)

```go
func processOrder(ctx context.Context, orderID string) error {
    tracer := tracing.GetTracer("order-service")
    
    // Parent span
    ctx, parentSpan := tracing.StartSpan(ctx, tracer, "process_order")
    defer parentSpan.End()
    
    // Child span 1: Validate
    if err := validateOrder(ctx, orderID); err != nil {
        tracing.RecordError(parentSpan, err)
        return err
    }
    
    // Child span 2: Payment
    if err := processPayment(ctx, orderID); err != nil {
        tracing.RecordError(parentSpan, err)
        return err
    }
    
    return nil
}

func validateOrder(ctx context.Context, orderID string) error {
    tracer := tracing.GetTracer("order-service")
    ctx, span := tracing.StartSpan(ctx, tracer, "validate_order")
    defer span.End()
    
    tracing.AddSpanAttributes(span,
        attribute.String("order.id", orderID),
    )
    
    // Validation logic...
    return nil
}
```

### Recording Errors

```go
func fetchUserData(ctx context.Context, userID string) error {
    tracer := tracing.GetTracer("user-service")
    ctx, span := tracing.StartSpan(ctx, tracer, "fetch_user_data")
    defer span.End()
    
    user, err := db.GetUser(userID)
    if err != nil {
        // Record error on span
        tracing.RecordError(span, err)
        return err
    }
    
    tracing.AddSpanAttributes(span,
        attribute.String("user.id", user.ID),
        attribute.String("user.email", user.Email),
    )
    
    return nil
}
```

### Adding Span Events

```go
func processLargeFile(ctx context.Context, filename string) error {
    tracer := tracing.GetTracer("file-processor")
    ctx, span := tracing.StartSpan(ctx, tracer, "process_large_file")
    defer span.End()
    
    // Add event: Started reading
    tracing.AddSpanEvent(span, "file.read.start",
        attribute.String("filename", filename),
    )
    
    // Read file...
    
    // Add event: Finished reading
    tracing.AddSpanEvent(span, "file.read.complete",
        attribute.Int("bytes_read", 1024000),
    )
    
    return nil
}
```

---

## Structured Logging

### Basic Logging

```go
package main

import (
    "github.com/gaurav/otel-observability/internal/logging"
    "go.uber.org/zap"
)

func main() {
    logger, _ := logging.NewLogger(logging.LoggerConfig{
        Level:       "info",
        Format:      "json",
        ServiceName: "my-service",
    })
    defer logger.Sync()
    
    logger.Info("Application started")
    logger.Error("Something went wrong", zap.Error(err))
}
```

### Context-Aware Logging (with Trace Correlation)

```go
func handleRequest(ctx context.Context, logger *logging.Logger) {
    // Logs will automatically include trace_id and span_id
    logger.InfoContext(ctx, "Processing request",
        zap.String("user_id", "user-123"),
        zap.String("method", "POST"),
    )
    
    // Your logic...
    
    if err != nil {
        logger.ErrorContext(ctx, "Request failed",
            zap.Error(err),
            zap.String("error_code", "E1001"),
        )
    }
}
```

### Logging with Additional Fields

```go
func processPayment(ctx context.Context, logger *logging.Logger, amount float64) {
    // Create logger with permanent fields
    paymentLogger := logger.WithFields(
        zap.String("component", "payment"),
        zap.Float64("amount", amount),
    )
    
    paymentLogger.InfoContext(ctx, "Payment processing started")
    
    // All logs from this logger will include component and amount
    paymentLogger.InfoContext(ctx, "Validating payment method")
    paymentLogger.InfoContext(ctx, "Payment completed successfully")
}
```

### Different Log Levels

```go
func debugExample(ctx context.Context, logger *logging.Logger) {
    // Debug - for development
    logger.DebugContext(ctx, "Debug information",
        zap.Any("config", config),
    )
    
    // Info - general information
    logger.InfoContext(ctx, "User logged in",
        zap.String("user_id", "123"),
    )
    
    // Warn - warning but not error
    logger.WarnContext(ctx, "API rate limit approaching",
        zap.Int("current_rate", 450),
        zap.Int("limit", 500),
    )
    
    // Error - errors that need attention
    logger.ErrorContext(ctx, "Database connection failed",
        zap.Error(err),
        zap.String("database", "postgres"),
    )
}
```

---

## Exception Tracking

### Basic Exception Recording

```go
package main

import (
    "context"
    "errors"
    "github.com/gaurav/otel-observability/internal/exceptions"
)

func processData(ctx context.Context, tracker *exceptions.Tracker) error {
    data, err := fetchData()
    if err != nil {
        // Record exception with full context
        tracker.RecordException(ctx, err, exceptions.Options{
            Severity:     exceptions.SeverityError,
            CaptureStack: true,
            Tags: map[string]string{
                "component": "data-processor",
                "operation": "fetch",
            },
        })
        return err
    }
    return nil
}
```

### Exception with Additional Data

```go
func processOrder(ctx context.Context, tracker *exceptions.Tracker, orderID string) error {
    err := validateOrder(orderID)
    if err != nil {
        tracker.RecordException(ctx, err, exceptions.Options{
            Severity:     exceptions.SeverityCritical,
            CaptureStack: true,
            Tags: map[string]string{
                "order_id": orderID,
                "stage":    "validation",
            },
            AdditionalData: map[string]interface{}{
                "order_amount":  1500.00,
                "payment_method": "credit_card",
                "user_id":       "user-789",
            },
        })
        return err
    }
    return nil
}
```

### Using Global Exception Tracker

```go
import "github.com/gaurav/otel-observability/internal/exceptions"

func init() {
    // Initialize global tracker
    logger, _ := zap.NewProduction()
    exceptions.InitGlobalTracker(logger)
}

func myFunction(ctx context.Context) error {
    err := doSomething()
    if err != nil {
        // Use global tracker
        exceptions.RecordException(ctx, err, exceptions.Options{
            Severity: exceptions.SeverityError,
        })
        return err
    }
    return nil
}
```

---

## Alerting

### Registering Alerts

```go
package main

import (
    "github.com/gaurav/otel-observability/internal/alerts"
    "go.uber.org/zap"
)

func setupAlerts(logger *zap.Logger) *alerts.Manager {
    manager := alerts.NewManager(logger)
    
    // Register console handler
    manager.RegisterHandler(alerts.NewConsoleHandler(logger))
    
    // Register error rate alert
    manager.RegisterAlert(alerts.Alert{
        Name:        "high_error_rate",
        Description: "Error rate exceeded threshold",
        Metric:      "error_rate",
        Threshold:   0.05, // 5%
        Severity:    alerts.SeverityCritical,
    })
    
    // Register latency alert
    manager.RegisterAlert(alerts.Alert{
        Name:        "high_latency",
        Description: "P95 latency exceeded threshold",
        Metric:      "latency_p95",
        Threshold:   1000, // 1000ms
        Severity:    alerts.SeverityWarning,
    })
    
    return manager
}
```

### Updating Metrics

```go
func handleRequest(manager *alerts.Manager) {
    startTime := time.Now()
    
    // Process request...
    
    duration := time.Since(startTime).Milliseconds()
    
    // Update metrics
    manager.IncrementCounter("request_count", 1)
    manager.RecordHistogram("latency_p95", float64(duration))
    
    if err != nil {
        manager.IncrementCounter("error_count", 1)
        
        // Calculate and update error rate
        reqCount, _ := manager.GetMetric("request_count")
        errCount, _ := manager.GetMetric("error_count")
        errorRate := errCount / reqCount
        manager.UpdateMetric("error_rate", errorRate)
    }
}
```

### Custom Alert Conditions

```go
func setupCustomAlert(manager *alerts.Manager) {
    manager.RegisterAlert(alerts.Alert{
        Name:        "complex_condition",
        Description: "Custom alert logic",
        Metric:      "custom_metric",
        Severity:    alerts.SeverityError,
        Condition: func(value float64) bool {
            // Custom logic
            return value > 100 && value < 200
        },
    })
}
```

---

## Full Integration

### Complete Example: HTTP Handler

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/gaurav/otel-observability/internal/alerts"
    "github.com/gaurav/otel-observability/internal/exceptions"
    "github.com/gaurav/otel-observability/internal/logging"
    "github.com/gaurav/otel-observability/internal/tracing"
    "go.opentelemetry.io/otel/attribute"
    "go.uber.org/zap"
)

type Server struct {
    logger           *logging.Logger
    alertManager     *alerts.Manager
    exceptionTracker *exceptions.Tracker
}

func (s *Server) handleOrder(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    tracer := tracing.GetTracer("order-api")
    
    // Start span
    ctx, span := tracing.StartSpan(ctx, tracer, "handle_order")
    defer span.End()
    
    // Extract order ID
    orderID := r.URL.Query().Get("order_id")
    
    // Add span attributes
    tracing.AddSpanAttributes(span,
        attribute.String("http.method", r.Method),
        attribute.String("http.url", r.URL.Path),
        attribute.String("order.id", orderID),
    )
    
    // Log request
    s.logger.InfoContext(ctx, "Processing order request",
        zap.String("order_id", orderID),
        zap.String("client_ip", r.RemoteAddr),
    )
    
    // Process order
    err := s.processOrder(ctx, orderID)
    if err != nil {
        // Record error on span
        tracing.RecordError(span, err)
        
        // Log error
        s.logger.ErrorContext(ctx, "Order processing failed",
            zap.Error(err),
            zap.String("order_id", orderID),
        )
        
        // Track exception
        s.exceptionTracker.RecordException(ctx, err, exceptions.Options{
            Severity:     exceptions.SeverityError,
            CaptureStack: true,
            Tags: map[string]string{
                "handler":  "order",
                "order_id": orderID,
            },
        })
        
        // Update metrics
        s.alertManager.IncrementCounter("order_errors", 1)
        
        http.Error(w, "Order processing failed", http.StatusInternalServerError)
        return
    }
    
    // Update success metrics
    s.alertManager.IncrementCounter("order_success", 1)
    
    // Log success
    s.logger.InfoContext(ctx, "Order processed successfully",
        zap.String("order_id", orderID),
    )
    
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Order processed"))
}

func (s *Server) processOrder(ctx context.Context, orderID string) error {
    tracer := tracing.GetTracer("order-service")
    
    // Validate order
    ctx, validateSpan := tracing.StartSpan(ctx, tracer, "validate_order")
    s.logger.DebugContext(ctx, "Validating order")
    // ... validation logic
    validateSpan.End()
    
    // Process payment
    ctx, paymentSpan := tracing.StartSpan(ctx, tracer, "process_payment")
    s.logger.DebugContext(ctx, "Processing payment")
    // ... payment logic
    paymentSpan.End()
    
    // Update inventory
    ctx, inventorySpan := tracing.StartSpan(ctx, tracer, "update_inventory")
    s.logger.DebugContext(ctx, "Updating inventory")
    // ... inventory logic
    inventorySpan.End()
    
    return nil
}
```

### Complete Application Setup

```go
package main

import (
    "context"
    "os"
    
    "github.com/gaurav/otel-observability/internal/alerts"
    "github.com/gaurav/otel-observability/internal/config"
    "github.com/gaurav/otel-observability/internal/exceptions"
    "github.com/gaurav/otel-observability/internal/logging"
    "github.com/gaurav/otel-observability/internal/tracing"
)

func main() {
    // Load configuration
    cfg, err := config.Load("configs/config.yaml")
    if err != nil {
        panic(err)
    }
    
    // Initialize logging
    logger, err := logging.NewLogger(logging.LoggerConfig{
        Level:       cfg.Logging.Level,
        Format:      cfg.Logging.Format,
        ServiceName: cfg.Service.Name,
    })
    if err != nil {
        panic(err)
    }
    defer logger.Sync()
    
    // Initialize tracing
    tp, err := tracing.InitTracer(tracing.TracerConfig{
        ServiceName:  cfg.Service.Name,
        Endpoint:     cfg.Tracing.Endpoint,
        Insecure:     cfg.Tracing.Insecure,
        SamplingRate: cfg.Tracing.SamplingRate,
    })
    if err != nil {
        logger.Fatal("Failed to initialize tracer", zap.Error(err))
    }
    defer tp.Shutdown(context.Background())
    
    // Initialize exception tracking
    exceptionTracker := exceptions.NewTracker(logger.Logger)
    
    // Initialize alerting
    alertManager := alerts.NewManager(logger.Logger)
    alertManager.RegisterHandler(alerts.NewConsoleHandler(logger.Logger))
    
    // Start alert evaluation
    ctx := context.Background()
    go alertManager.Start(ctx, 30*time.Second)
    
    // Create server
    server := &Server{
        logger:           logger,
        alertManager:     alertManager,
        exceptionTracker: exceptionTracker,
    }
    
    // Setup HTTP routes
    http.HandleFunc("/order", server.handleOrder)
    
    // Start server
    logger.Info("Server starting", zap.Int("port", 8080))
    http.ListenAndServe(":8080", nil)
}
```

---

## Best Practices

### 1. Always Close Resources

```go
defer span.End()
defer logger.Sync()
defer tp.Shutdown(ctx)
```

### 2. Use Context for Correlation

```go
// Always pass context through the call chain
func handleRequest(ctx context.Context) {
    processData(ctx)  // Pass context
}
```

### 3. Add Meaningful Attributes

```go
tracing.AddSpanAttributes(span,
    attribute.String("user.id", userID),
    attribute.Int("batch.size", len(items)),
    attribute.Bool("cache.hit", true),
)
```

### 4. Log Structured Data

```go
logger.InfoContext(ctx, "Order created",
    zap.String("order_id", orderID),
    zap.Float64("amount", 99.99),
    zap.Time("created_at", time.Now()),
)
```

### 5. Track Important Errors

```go
if err != nil {
    tracing.RecordError(span, err)
    logger.ErrorContext(ctx, "Operation failed", zap.Error(err))
    exceptionTracker.RecordException(ctx, err, exceptions.Options{...})
}
```

---

For more examples, see `cmd/app/main.go` in the repository!

