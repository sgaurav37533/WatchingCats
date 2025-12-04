package exceptions

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Severity represents the severity level of an exception
type Severity string

const (
	SeverityDebug    Severity = "debug"
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityError    Severity = "error"
	SeverityCritical Severity = "critical"
)

// Exception represents a captured exception
type Exception struct {
	ID         string
	Message    string
	Type       string
	Severity   Severity
	Timestamp  time.Time
	StackTrace []StackFrame
	Tags       map[string]string
	TraceID    string
	SpanID     string
	Context    map[string]interface{}
}

// StackFrame represents a single stack frame
type StackFrame struct {
	Function string
	File     string
	Line     int
}

// Options for recording exceptions
type Options struct {
	Severity       Severity
	Tags           map[string]string
	CaptureStack   bool
	MaxStackDepth  int
	AdditionalData map[string]interface{}
}

// Tracker tracks and reports exceptions
type Tracker struct {
	exceptions      []Exception
	mu              sync.RWMutex
	logger          *zap.Logger
	captureStack    bool
	maxStackDepth   int
	groupByMessage  bool
	ignorePatterns  []string
	
	// Exception grouping
	groups map[string]int // fingerprint -> count
}

// NewTracker creates a new exception tracker
func NewTracker(logger *zap.Logger) *Tracker {
	return &Tracker{
		exceptions:     make([]Exception, 0),
		logger:         logger,
		captureStack:   true,
		maxStackDepth:  50,
		groupByMessage: true,
		ignorePatterns: make([]string, 0),
		groups:         make(map[string]int),
	}
}

// SetCaptureStack enables or disables stack trace capture
func (t *Tracker) SetCaptureStack(capture bool) {
	t.captureStack = capture
}

// SetMaxStackDepth sets the maximum stack depth to capture
func (t *Tracker) SetMaxStackDepth(depth int) {
	t.maxStackDepth = depth
}

// AddIgnorePattern adds a pattern to ignore
func (t *Tracker) AddIgnorePattern(pattern string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ignorePatterns = append(t.ignorePatterns, pattern)
}

// shouldIgnore checks if an error should be ignored
func (t *Tracker) shouldIgnore(err error) bool {
	errMsg := err.Error()
	for _, pattern := range t.ignorePatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}
	return false
}

// RecordException records an exception with context
func (t *Tracker) RecordException(ctx context.Context, err error, opts Options) {
	if err == nil {
		return
	}

	// Check if should ignore
	if t.shouldIgnore(err) {
		return
	}

	// Extract trace context
	span := trace.SpanFromContext(ctx)
	spanCtx := span.SpanContext()

	var traceID, spanID string
	if spanCtx.HasTraceID() {
		traceID = spanCtx.TraceID().String()
	}
	if spanCtx.HasSpanID() {
		spanID = spanCtx.SpanID().String()
	}

	// Capture stack trace
	var stackTrace []StackFrame
	if opts.CaptureStack || (opts.CaptureStack == false && t.captureStack) {
		stackTrace = captureStackTrace(t.maxStackDepth)
	}

	// Create exception
	exception := Exception{
		ID:         generateID(),
		Message:    err.Error(),
		Type:       fmt.Sprintf("%T", err),
		Severity:   opts.Severity,
		Timestamp:  time.Now(),
		StackTrace: stackTrace,
		Tags:       opts.Tags,
		TraceID:    traceID,
		SpanID:     spanID,
		Context:    opts.AdditionalData,
	}

	// Default severity if not provided
	if exception.Severity == "" {
		exception.Severity = SeverityError
	}

	// Store exception
	t.mu.Lock()
	t.exceptions = append(t.exceptions, exception)
	
	// Update grouping
	fingerprint := t.getFingerprint(exception)
	t.groups[fingerprint]++
	t.mu.Unlock()

	// Log the exception
	t.logger.Error("Exception recorded",
		zap.String("exception_id", exception.ID),
		zap.String("message", exception.Message),
		zap.String("type", exception.Type),
		zap.String("severity", string(exception.Severity)),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	)

	// Record on span
	if span.IsRecording() {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		
		attrs := []attribute.KeyValue{
			attribute.String("exception.id", exception.ID),
			attribute.String("exception.type", exception.Type),
			attribute.String("exception.severity", string(exception.Severity)),
		}
		
		for key, value := range opts.Tags {
			attrs = append(attrs, attribute.String(fmt.Sprintf("exception.tag.%s", key), value))
		}
		
		span.SetAttributes(attrs...)
	}
}

// GetExceptions returns all recorded exceptions
func (t *Tracker) GetExceptions() []Exception {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	exceptions := make([]Exception, len(t.exceptions))
	copy(exceptions, t.exceptions)
	return exceptions
}

// GetExceptionGroups returns exception counts by fingerprint
func (t *Tracker) GetExceptionGroups() map[string]int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	
	groups := make(map[string]int)
	for k, v := range t.groups {
		groups[k] = v
	}
	return groups
}

// Clear clears all recorded exceptions
func (t *Tracker) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	t.exceptions = make([]Exception, 0)
	t.groups = make(map[string]int)
}

// getFingerprint generates a fingerprint for exception grouping
func (t *Tracker) getFingerprint(e Exception) string {
	if t.groupByMessage {
		return fmt.Sprintf("%s:%s", e.Type, e.Message)
	}
	
	// Group by type and first stack frame
	if len(e.StackTrace) > 0 {
		frame := e.StackTrace[0]
		return fmt.Sprintf("%s:%s:%d", e.Type, frame.Function, frame.Line)
	}
	
	return e.Type
}

// captureStackTrace captures the current stack trace
func captureStackTrace(maxDepth int) []StackFrame {
	pcs := make([]uintptr, maxDepth)
	n := runtime.Callers(3, pcs) // Skip 3 frames (Callers, captureStackTrace, RecordException)
	
	frames := make([]StackFrame, 0, n)
	callersFrames := runtime.CallersFrames(pcs[:n])
	
	for {
		frame, more := callersFrames.Next()
		
		frames = append(frames, StackFrame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		
		if !more {
			break
		}
	}
	
	return frames
}

// generateID generates a unique ID for an exception
func generateID() string {
	return fmt.Sprintf("exc_%d", time.Now().UnixNano())
}

// Global exception tracker
var globalTracker *Tracker

// InitGlobalTracker initializes the global exception tracker
func InitGlobalTracker(logger *zap.Logger) {
	globalTracker = NewTracker(logger)
}

// GetGlobalTracker returns the global exception tracker
func GetGlobalTracker() *Tracker {
	return globalTracker
}

// RecordException records an exception using the global tracker
func RecordException(ctx context.Context, err error, opts Options) {
	if globalTracker != nil {
		globalTracker.RecordException(ctx, err, opts)
	}
}

