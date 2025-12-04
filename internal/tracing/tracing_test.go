package tracing

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestTraceIDFromContext(t *testing.T) {
	ctx := context.Background()
	
	traceID := TraceIDFromContext(ctx)
	if traceID != "" {
		t.Errorf("Expected empty trace ID for context without span, got %s", traceID)
	}
}

func TestSpanIDFromContext(t *testing.T) {
	ctx := context.Background()
	
	spanID := SpanIDFromContext(ctx)
	if spanID != "" {
		t.Errorf("Expected empty span ID for context without span, got %s", spanID)
	}
}

func TestAddSpanAttributes(t *testing.T) {
	ctx := context.Background()
	span := SpanFromContext(ctx)
	
	// Should not panic when adding attributes to a non-recording span
	AddSpanAttributes(span, 
		attribute.String("test.key", "test.value"),
		attribute.Int("test.number", 42),
	)
}

