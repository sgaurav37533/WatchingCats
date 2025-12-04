package exceptions

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
)

func TestNewTracker(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tracker := NewTracker(logger)

	if tracker == nil {
		t.Fatal("Tracker is nil")
	}

	if len(tracker.exceptions) != 0 {
		t.Errorf("Expected 0 exceptions, got %d", len(tracker.exceptions))
	}
}

func TestRecordException(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tracker := NewTracker(logger)

	ctx := context.Background()
	err := errors.New("test error")

	tracker.RecordException(ctx, err, Options{
		Severity:     SeverityError,
		CaptureStack: true,
		Tags: map[string]string{
			"component": "test",
		},
	})

	exceptions := tracker.GetExceptions()
	if len(exceptions) != 1 {
		t.Fatalf("Expected 1 exception, got %d", len(exceptions))
	}

	exc := exceptions[0]
	if exc.Message != "test error" {
		t.Errorf("Expected message 'test error', got '%s'", exc.Message)
	}

	if exc.Severity != SeverityError {
		t.Errorf("Expected severity error, got %s", exc.Severity)
	}

	if len(exc.StackTrace) == 0 {
		t.Error("Expected stack trace to be captured")
	}

	if exc.Tags["component"] != "test" {
		t.Error("Expected tag 'component' to be 'test'")
	}
}

func TestIgnorePatterns(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tracker := NewTracker(logger)

	tracker.AddIgnorePattern("context canceled")

	ctx := context.Background()
	err := errors.New("context canceled")

	tracker.RecordException(ctx, err, Options{
		Severity: SeverityError,
	})

	exceptions := tracker.GetExceptions()
	if len(exceptions) != 0 {
		t.Errorf("Expected 0 exceptions (should be ignored), got %d", len(exceptions))
	}
}

func TestExceptionGrouping(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tracker := NewTracker(logger)

	ctx := context.Background()
	err := errors.New("database error")

	// Record the same error twice
	tracker.RecordException(ctx, err, Options{Severity: SeverityError})
	tracker.RecordException(ctx, err, Options{Severity: SeverityError})

	groups := tracker.GetExceptionGroups()
	if len(groups) != 1 {
		t.Errorf("Expected 1 exception group, got %d", len(groups))
	}

	// Should have count of 2
	for _, count := range groups {
		if count != 2 {
			t.Errorf("Expected group count 2, got %d", count)
		}
	}
}

func TestClear(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	tracker := NewTracker(logger)

	ctx := context.Background()
	err := errors.New("test error")

	tracker.RecordException(ctx, err, Options{Severity: SeverityError})

	if len(tracker.GetExceptions()) != 1 {
		t.Fatal("Should have 1 exception before clear")
	}

	tracker.Clear()

	if len(tracker.GetExceptions()) != 0 {
		t.Error("Should have 0 exceptions after clear")
	}

	if len(tracker.GetExceptionGroups()) != 0 {
		t.Error("Should have 0 groups after clear")
	}
}

