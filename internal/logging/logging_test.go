package logging

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	cfg := LoggerConfig{
		Level:         "info",
		Format:        "json",
		ServiceName:   "test-service",
		IncludeCaller: true,
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Logger is nil")
	}

	if logger.serviceName != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", logger.serviceName)
	}
}

func TestLoggerWithContext(t *testing.T) {
	cfg := LoggerConfig{
		Level:       "info",
		Format:      "json",
		ServiceName: "test-service",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	ctx := context.Background()
	contextLogger := logger.WithContext(ctx)

	if contextLogger == nil {
		t.Fatal("Context logger is nil")
	}
}

func TestLoggerWithFields(t *testing.T) {
	cfg := LoggerConfig{
		Level:       "info",
		Format:      "json",
		ServiceName: "test-service",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	fieldLogger := logger.WithFields(
		zap.String("key1", "value1"),
		zap.Int("key2", 42),
	)

	if fieldLogger == nil {
		t.Fatal("Field logger is nil")
	}

	if fieldLogger.serviceName != logger.serviceName {
		t.Error("Service name should be preserved")
	}
}

