package models

import "time"

// Span represents a distributed trace span
type Span struct {
	TraceID    string            `json:"trace_id"`
	SpanID     string            `json:"span_id"`
	ParentID   string            `json:"parent_id,omitempty"`
	Name       string            `json:"name"`
	Kind       string            `json:"kind"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
	Attributes map[string]string `json:"attributes"`
	Events     []SpanEvent       `json:"events"`
	Status     SpanStatus        `json:"status"`
}

// SpanEvent represents an event within a span
type SpanEvent struct {
	Name       string            `json:"name"`
	Timestamp  time.Time         `json:"timestamp"`
	Attributes map[string]string `json:"attributes"`
}

// SpanStatus represents the status of a span
type SpanStatus struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

// LogRecord represents a log entry
type LogRecord struct {
	Timestamp   time.Time         `json:"timestamp"`
	TraceID     string            `json:"trace_id,omitempty"`
	SpanID      string            `json:"span_id,omitempty"`
	Severity    string            `json:"severity"`
	Message     string            `json:"message"`
	Attributes  map[string]string `json:"attributes"`
	ServiceName string            `json:"service_name"`
}

// Metric represents a metric data point
type Metric struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"` // counter, gauge, histogram
	Value       float64           `json:"value"`
	Timestamp   time.Time         `json:"timestamp"`
	Attributes  map[string]string `json:"attributes"`
	ServiceName string            `json:"service_name"`
}

// ExceptionRecord represents an exception/error
type ExceptionRecord struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Message     string            `json:"message"`
	Severity    string            `json:"severity"`
	Timestamp   time.Time         `json:"timestamp"`
	TraceID     string            `json:"trace_id,omitempty"`
	SpanID      string            `json:"span_id,omitempty"`
	StackTrace  string            `json:"stack_trace"`
	Tags        map[string]string `json:"tags"`
	ServiceName string            `json:"service_name"`
}

// TelemetryBatch represents a batch of telemetry data
type TelemetryBatch struct {
	Spans      []Span            `json:"spans"`
	Logs       []LogRecord       `json:"logs"`
	Metrics    []Metric          `json:"metrics"`
	Exceptions []ExceptionRecord `json:"exceptions"`
}

