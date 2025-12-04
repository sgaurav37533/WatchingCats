package alerts

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Severity represents the severity level of an alert
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeverityCritical Severity = "critical"
)

// Alert represents an alert definition
type Alert struct {
	Name        string
	Description string
	Metric      string
	Threshold   float64
	Window      time.Duration
	Severity    Severity
	Condition   func(value float64) bool
}

// AlertEvent represents a triggered alert
type AlertEvent struct {
	Alert     Alert
	Value     float64
	Timestamp time.Time
	Message   string
}

// AlertHandler handles triggered alerts
type AlertHandler interface {
	Handle(event AlertEvent) error
}

// ConsoleHandler logs alerts to console
type ConsoleHandler struct {
	logger *zap.Logger
}

// NewConsoleHandler creates a new console alert handler
func NewConsoleHandler(logger *zap.Logger) *ConsoleHandler {
	return &ConsoleHandler{
		logger: logger,
	}
}

// Handle handles an alert by logging to console
func (h *ConsoleHandler) Handle(event AlertEvent) error {
	h.logger.Warn("Alert triggered",
		zap.String("alert_name", event.Alert.Name),
		zap.String("description", event.Alert.Description),
		zap.String("severity", string(event.Alert.Severity)),
		zap.Float64("value", event.Value),
		zap.Float64("threshold", event.Alert.Threshold),
		zap.Time("timestamp", event.Timestamp),
		zap.String("message", event.Message),
	)
	return nil
}

// WebhookHandler sends alerts to a webhook
type WebhookHandler struct {
	url    string
	logger *zap.Logger
}

// NewWebhookHandler creates a new webhook alert handler
func NewWebhookHandler(url string, logger *zap.Logger) *WebhookHandler {
	return &WebhookHandler{
		url:    url,
		logger: logger,
	}
}

// Handle handles an alert by sending to webhook
func (h *WebhookHandler) Handle(event AlertEvent) error {
	// In a real implementation, this would send HTTP POST to webhook
	h.logger.Info("Sending alert to webhook",
		zap.String("url", h.url),
		zap.String("alert", event.Alert.Name),
	)
	return nil
}

// Manager manages alerts and their evaluation
type Manager struct {
	alerts   []Alert
	handlers []AlertHandler
	metrics  map[string]float64
	mu       sync.RWMutex
	logger   *zap.Logger
	
	// For deduplication
	firedAlerts map[string]time.Time
	cooldown    time.Duration
}

// NewManager creates a new alert manager
func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		alerts:      make([]Alert, 0),
		handlers:    make([]AlertHandler, 0),
		metrics:     make(map[string]float64),
		logger:      logger,
		firedAlerts: make(map[string]time.Time),
		cooldown:    5 * time.Minute, // Don't repeat same alert within 5 minutes
	}
}

// RegisterAlert registers a new alert
func (m *Manager) RegisterAlert(alert Alert) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alerts = append(m.alerts, alert)
	m.logger.Info("Alert registered",
		zap.String("name", alert.Name),
		zap.String("metric", alert.Metric),
		zap.Float64("threshold", alert.Threshold),
	)
}

// RegisterHandler registers an alert handler
func (m *Manager) RegisterHandler(handler AlertHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers = append(m.handlers, handler)
}

// UpdateMetric updates a metric value
func (m *Manager) UpdateMetric(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics[name] = value
}

// GetMetric retrieves a metric value
func (m *Manager) GetMetric(name string) (float64, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.metrics[name]
	return value, exists
}

// Evaluate evaluates all alerts against current metrics
func (m *Manager) Evaluate(ctx context.Context) {
	m.mu.RLock()
	alerts := make([]Alert, len(m.alerts))
	copy(alerts, m.alerts)
	metrics := make(map[string]float64)
	for k, v := range m.metrics {
		metrics[k] = v
	}
	m.mu.RUnlock()

	for _, alert := range alerts {
		value, exists := metrics[alert.Metric]
		if !exists {
			continue
		}

		// Check if alert condition is met
		shouldFire := false
		if alert.Condition != nil {
			shouldFire = alert.Condition(value)
		} else {
			// Default condition: value exceeds threshold
			shouldFire = value > alert.Threshold
		}

		if shouldFire {
			m.fireAlert(alert, value)
		}
	}
}

// fireAlert fires an alert if not in cooldown
func (m *Manager) fireAlert(alert Alert, value float64) {
	m.mu.Lock()
	
	// Check cooldown
	if lastFired, exists := m.firedAlerts[alert.Name]; exists {
		if time.Since(lastFired) < m.cooldown {
			m.mu.Unlock()
			return
		}
	}
	
	m.firedAlerts[alert.Name] = time.Now()
	handlers := make([]AlertHandler, len(m.handlers))
	copy(handlers, m.handlers)
	m.mu.Unlock()

	event := AlertEvent{
		Alert:     alert,
		Value:     value,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("%s: %s (value: %.2f, threshold: %.2f)", 
			alert.Name, alert.Description, value, alert.Threshold),
	}

	// Notify all handlers
	for _, handler := range handlers {
		if err := handler.Handle(event); err != nil {
			m.logger.Error("Failed to handle alert",
				zap.String("alert", alert.Name),
				zap.Error(err),
			)
		}
	}
}

// Start starts the alert evaluation loop
func (m *Manager) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	m.logger.Info("Alert manager started",
		zap.Duration("interval", interval),
	)

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("Alert manager stopped")
			return
		case <-ticker.C:
			m.Evaluate(ctx)
		}
	}
}

// IncrementCounter increments a counter metric
func (m *Manager) IncrementCounter(name string, delta float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics[name] += delta
}

// RecordHistogram records a histogram value (simplified)
func (m *Manager) RecordHistogram(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// For simplicity, just update the latest value
	// In production, you'd calculate percentiles, etc.
	m.metrics[name] = value
}

// Reset resets all metrics
func (m *Manager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metrics = make(map[string]float64)
}

