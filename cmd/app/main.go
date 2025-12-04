package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gaurav/otel-observability/internal/alerts"
	"github.com/gaurav/otel-observability/internal/config"
	"github.com/gaurav/otel-observability/internal/exceptions"
	"github.com/gaurav/otel-observability/internal/logging"
	"github.com/gaurav/otel-observability/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging
	logger, err := logging.NewLogger(logging.LoggerConfig{
		Level:              cfg.Logging.Level,
		Format:             cfg.Logging.Format,
		ServiceName:        cfg.Service.Name,
		CorrelationEnabled: cfg.Logging.CorrelationEnabled,
		IncludeCaller:      cfg.Logging.IncludeCaller,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting application",
		zap.String("service", cfg.Service.Name),
		zap.String("version", cfg.Service.Version),
		zap.String("environment", cfg.Service.Environment),
	)

	// Initialize tracing
	if cfg.Tracing.Enabled {
		tp, err := tracing.InitTracer(tracing.TracerConfig{
			ServiceName:    cfg.Service.Name,
			ServiceVersion: cfg.Service.Version,
			Environment:    cfg.Service.Environment,
			Endpoint:       cfg.Tracing.Endpoint,
			Insecure:       cfg.Tracing.Insecure,
			SamplingRate:   cfg.Tracing.SamplingRate,
		})
		if err != nil {
			logger.Fatal("Failed to initialize tracer", zap.Error(err))
		}
		defer tp.Shutdown(context.Background())
		logger.Info("Tracing initialized", zap.String("endpoint", cfg.Tracing.Endpoint))
	}

	// Initialize exception tracking
	exceptionTracker := exceptions.NewTracker(logger.Logger)
	exceptionTracker.SetMaxStackDepth(cfg.Exceptions.MaxStackDepth)
	exceptionTracker.SetCaptureStack(cfg.Exceptions.CaptureStackTrace)
	for _, pattern := range cfg.Exceptions.IgnorePatterns {
		exceptionTracker.AddIgnorePattern(pattern)
	}
	exceptions.InitGlobalTracker(logger.Logger)
	logger.Info("Exception tracking initialized")

	// Initialize alert manager
	alertManager := alerts.NewManager(logger.Logger)
	alertManager.RegisterHandler(alerts.NewConsoleHandler(logger.Logger))

	// Register sample alerts based on config
	for _, rule := range cfg.Alerts.Rules {
		alert := alerts.Alert{
			Name:        rule.Name,
			Description: rule.Description,
			Metric:      rule.Metric,
			Threshold:   rule.Threshold,
			Severity:    alerts.Severity(rule.Severity),
		}
		alertManager.RegisterAlert(alert)
	}
	logger.Info("Alert manager initialized", zap.Int("rules", len(cfg.Alerts.Rules)))

	// Start alert evaluation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if cfg.Alerts.Enabled {
		go alertManager.Start(ctx, cfg.Alerts.GetEvaluationInterval())
	}

	// Run sample workload
	logger.Info("Starting sample workload...")
	go runSampleWorkload(ctx, logger, alertManager, exceptionTracker)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down gracefully...")
	cancel()
	time.Sleep(2 * time.Second)
	logger.Info("Application stopped")
}

// runSampleWorkload simulates a workload that generates telemetry data
func runSampleWorkload(ctx context.Context, logger *logging.Logger, alertMgr *alerts.Manager, exTracker *exceptions.Tracker) {
	tracer := tracing.GetTracer("sample-app")
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	requestCount := 0
	errorCount := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			requestCount++
			
			// Simulate processing a request with distributed tracing
			reqCtx, span := tracing.StartSpan(ctx, tracer, "process_request")
			
			// Add span attributes
			tracing.AddSpanAttributes(span,
				attribute.String("request.id", fmt.Sprintf("req-%d", requestCount)),
				attribute.String("user.id", fmt.Sprintf("user-%d", rand.Intn(100))),
				attribute.String("method", "GET"),
			)

			// Log with trace context
			logger.InfoContext(reqCtx, "Processing request",
				zap.Int("request_number", requestCount),
			)

			// Simulate some work with child spans
			processRequest(reqCtx, tracer, logger, requestCount)

			// Randomly generate errors to test exception tracking and alerts
			if rand.Float64() < 0.15 { // 15% error rate
				errorCount++
				err := errors.New("simulated database connection error")
				
				// Record exception with context
				exTracker.RecordException(reqCtx, err, exceptions.Options{
					Severity:    exceptions.SeverityError,
					CaptureStack: true,
					Tags: map[string]string{
						"component": "database",
						"operation": "query",
					},
				})

				logger.ErrorContext(reqCtx, "Request failed", zap.Error(err))
				tracing.RecordError(span, err)
			} else {
				logger.InfoContext(reqCtx, "Request completed successfully")
			}

			span.End()

			// Update metrics for alerting
			errorRate := float64(errorCount) / float64(requestCount)
			alertMgr.UpdateMetric("error_rate", errorRate)
			alertMgr.UpdateMetric("request_count", float64(requestCount))
			alertMgr.UpdateMetric("error_count", float64(errorCount))

			// Log metrics periodically
			if requestCount%10 == 0 {
				logger.Info("Metrics snapshot",
					zap.Int("total_requests", requestCount),
					zap.Int("total_errors", errorCount),
					zap.Float64("error_rate", errorRate),
				)
			}
		}
	}
}

// processRequest simulates processing a request with multiple operations
func processRequest(ctx context.Context, tracer trace.Tracer, logger *logging.Logger, requestNum int) {
	// Simulate database query
	dbCtx, dbSpan := tracing.StartSpan(ctx, tracer, "database_query")
	tracing.AddSpanAttributes(dbSpan,
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "SELECT"),
	)
	
	time.Sleep(time.Duration(10+rand.Intn(50)) * time.Millisecond)
	logger.DebugContext(dbCtx, "Database query executed")
	dbSpan.End()

	// Simulate external API call
	apiCtx, apiSpan := tracing.StartSpan(ctx, tracer, "external_api_call")
	tracing.AddSpanAttributes(apiSpan,
		attribute.String("http.method", "GET"),
		attribute.String("http.url", "https://api.example.com/data"),
	)
	
	time.Sleep(time.Duration(20+rand.Intn(80)) * time.Millisecond)
	logger.DebugContext(apiCtx, "External API called")
	apiSpan.End()

	// Simulate cache check
	cacheCtx, cacheSpan := tracing.StartSpan(ctx, tracer, "cache_lookup")
	tracing.AddSpanAttributes(cacheSpan,
		attribute.String("cache.key", fmt.Sprintf("user-data-%d", requestNum)),
		attribute.Bool("cache.hit", rand.Float64() > 0.3),
	)
	
	time.Sleep(time.Duration(5+rand.Intn(15)) * time.Millisecond)
	logger.DebugContext(cacheCtx, "Cache checked")
	cacheSpan.End()
}

