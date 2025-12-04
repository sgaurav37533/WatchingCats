package alerts

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewManager(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	if manager == nil {
		t.Fatal("Manager is nil")
	}

	if len(manager.alerts) != 0 {
		t.Errorf("Expected 0 alerts, got %d", len(manager.alerts))
	}
}

func TestRegisterAlert(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	alert := Alert{
		Name:      "test_alert",
		Metric:    "test_metric",
		Threshold: 100.0,
		Severity:  SeverityWarning,
	}

	manager.RegisterAlert(alert)

	if len(manager.alerts) != 1 {
		t.Errorf("Expected 1 alert, got %d", len(manager.alerts))
	}
}

func TestUpdateMetric(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	manager.UpdateMetric("test_metric", 42.5)

	value, exists := manager.GetMetric("test_metric")
	if !exists {
		t.Fatal("Metric should exist")
	}

	if value != 42.5 {
		t.Errorf("Expected metric value 42.5, got %f", value)
	}
}

func TestIncrementCounter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	manager.IncrementCounter("test_counter", 5)
	manager.IncrementCounter("test_counter", 3)

	value, exists := manager.GetMetric("test_counter")
	if !exists {
		t.Fatal("Counter should exist")
	}

	if value != 8 {
		t.Errorf("Expected counter value 8, got %f", value)
	}
}

func TestAlertEvaluation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := NewManager(logger)

	alert := Alert{
		Name:      "high_value",
		Metric:    "test_metric",
		Threshold: 100.0,
		Severity:  SeverityCritical,
	}

	manager.RegisterAlert(alert)
	manager.UpdateMetric("test_metric", 50.0)

	ctx := context.Background()
	manager.Evaluate(ctx)

	// Should not fire because value < threshold
	time.Sleep(100 * time.Millisecond)
}

