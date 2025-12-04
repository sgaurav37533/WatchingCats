package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gaurav/otel-observability/internal/config"
	"github.com/gaurav/otel-observability/internal/logging"
	"github.com/gaurav/otel-observability/pkg/models"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Collector receives and processes telemetry data
type Collector struct {
	config   *config.Config
	logger   *logging.Logger
	server   *grpc.Server
	
	// Storage
	spans      []models.Span
	logs       []models.LogRecord
	metrics    []models.Metric
	exceptions []models.ExceptionRecord
	mu         sync.RWMutex
}

// NewCollector creates a new telemetry collector
func NewCollector(cfg *config.Config, logger *logging.Logger) *Collector {
	return &Collector{
		config:     cfg,
		logger:     logger,
		spans:      make([]models.Span, 0),
		logs:       make([]models.LogRecord, 0),
		metrics:    make([]models.Metric, 0),
		exceptions: make([]models.ExceptionRecord, 0),
	}
}

// Start starts the collector service
func (c *Collector) Start(ctx context.Context) error {
	// Start gRPC server
	listener, err := net.Listen("tcp", c.config.Collector.GRPC.Endpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	c.server = grpc.NewServer()
	
	// Note: In a real implementation, you would register OTLP service handlers here
	// For this skeleton, we're showing the structure
	
	c.logger.Info("Collector gRPC server starting",
		zap.String("endpoint", c.config.Collector.GRPC.Endpoint),
	)

	// Start background processors
	go c.processIncomingData(ctx)
	go c.exportData(ctx)
	go c.printStats(ctx)

	// Start server
	go func() {
		if err := c.server.Serve(listener); err != nil {
			c.logger.Error("gRPC server error", zap.Error(err))
		}
	}()

	return nil
}

// Stop stops the collector service
func (c *Collector) Stop() {
	if c.server != nil {
		c.logger.Info("Stopping collector server...")
		c.server.GracefulStop()
	}
}

// ReceiveSpan receives a trace span
func (c *Collector) ReceiveSpan(span models.Span) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.spans = append(c.spans, span)
	
	c.logger.Debug("Span received",
		zap.String("trace_id", span.TraceID),
		zap.String("span_id", span.SpanID),
		zap.String("name", span.Name),
	)
}

// ReceiveLog receives a log record
func (c *Collector) ReceiveLog(log models.LogRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logs = append(c.logs, log)
	
	c.logger.Debug("Log received",
		zap.String("severity", log.Severity),
		zap.String("message", log.Message),
	)
}

// ReceiveMetric receives a metric
func (c *Collector) ReceiveMetric(metric models.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = append(c.metrics, metric)
	
	c.logger.Debug("Metric received",
		zap.String("name", metric.Name),
		zap.Float64("value", metric.Value),
	)
}

// ReceiveException receives an exception record
func (c *Collector) ReceiveException(exception models.ExceptionRecord) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.exceptions = append(c.exceptions, exception)
	
	c.logger.Debug("Exception received",
		zap.String("id", exception.ID),
		zap.String("type", exception.Type),
		zap.String("message", exception.Message),
	)
}

// processIncomingData processes incoming telemetry data
func (c *Collector) processIncomingData(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.RLock()
			spanCount := len(c.spans)
			logCount := len(c.logs)
			metricCount := len(c.metrics)
			exceptionCount := len(c.exceptions)
			c.mu.RUnlock()

			if spanCount > 0 || logCount > 0 || metricCount > 0 || exceptionCount > 0 {
				c.logger.Debug("Processing telemetry data",
					zap.Int("spans", spanCount),
					zap.Int("logs", logCount),
					zap.Int("metrics", metricCount),
					zap.Int("exceptions", exceptionCount),
				)
			}

			// In a real implementation, here you would:
			// 1. Process and transform data
			// 2. Apply sampling/filtering rules
			// 3. Enrich with additional context
			// 4. Detect anomalies
		}
	}
}

// exportData exports telemetry data to configured backends
func (c *Collector) exportData(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.Lock()
			batch := models.TelemetryBatch{
				Spans:      c.spans,
				Logs:       c.logs,
				Metrics:    c.metrics,
				Exceptions: c.exceptions,
			}
			
			// Clear after copying
			c.spans = make([]models.Span, 0)
			c.logs = make([]models.LogRecord, 0)
			c.metrics = make([]models.Metric, 0)
			c.exceptions = make([]models.ExceptionRecord, 0)
			c.mu.Unlock()

			if len(batch.Spans) > 0 || len(batch.Logs) > 0 || 
			   len(batch.Metrics) > 0 || len(batch.Exceptions) > 0 {
				c.exportBatch(batch)
			}
		}
	}
}

// exportBatch exports a batch of telemetry data
func (c *Collector) exportBatch(batch models.TelemetryBatch) {
	// Export to Jaeger (if enabled)
	if cfg, ok := c.config.Collector.Exporters["jaeger"]; ok && cfg.Enabled {
		c.logger.Debug("Exporting to Jaeger",
			zap.Int("spans", len(batch.Spans)),
		)
		// In real implementation: send to Jaeger
	}

	// Export to Prometheus (if enabled)
	if cfg, ok := c.config.Collector.Exporters["prometheus"]; ok && cfg.Enabled {
		c.logger.Debug("Exporting to Prometheus",
			zap.Int("metrics", len(batch.Metrics)),
		)
		// In real implementation: send to Prometheus
	}

	// Export to Elasticsearch (if enabled)
	if cfg, ok := c.config.Collector.Exporters["elasticsearch"]; ok && cfg.Enabled {
		c.logger.Debug("Exporting to Elasticsearch",
			zap.Int("logs", len(batch.Logs)),
			zap.Int("exceptions", len(batch.Exceptions)),
		)
		// In real implementation: send to Elasticsearch
	}

	c.logger.Info("Batch exported",
		zap.Int("spans", len(batch.Spans)),
		zap.Int("logs", len(batch.Logs)),
		zap.Int("metrics", len(batch.Metrics)),
		zap.Int("exceptions", len(batch.Exceptions)),
	)
}

// printStats prints periodic statistics
func (c *Collector) printStats(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.RLock()
			spanCount := len(c.spans)
			logCount := len(c.logs)
			metricCount := len(c.metrics)
			exceptionCount := len(c.exceptions)
			c.mu.RUnlock()

			uptime := time.Since(startTime)

			c.logger.Info("Collector statistics",
				zap.Duration("uptime", uptime),
				zap.Int("spans_buffered", spanCount),
				zap.Int("logs_buffered", logCount),
				zap.Int("metrics_buffered", metricCount),
				zap.Int("exceptions_buffered", exceptionCount),
			)
		}
	}
}

func main() {
	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logging
	logger, err := logging.NewLogger(logging.LoggerConfig{
		Level:         cfg.Logging.Level,
		Format:        cfg.Logging.Format,
		ServiceName:   "otel-collector",
		IncludeCaller: true,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting OpenTelemetry Collector",
		zap.String("version", cfg.Service.Version),
		zap.String("environment", cfg.Service.Environment),
	)

	// Create collector (pass the underlying zap logger)
	zapLogger := logger.Logger
	_ = zapLogger // We'll use the logger wrapper instead
	collector := NewCollector(cfg, logger)

	// Start collector
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := collector.Start(ctx); err != nil {
		logger.Fatal("Failed to start collector", zap.Error(err))
	}

	logger.Info("Collector started successfully")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down collector...")
	collector.Stop()
	cancel()
	
	time.Sleep(2 * time.Second)
	logger.Info("Collector stopped")
}

