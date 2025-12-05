package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// JaegerDAO handles interactions with Jaeger backend
type JaegerDAO struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

// Trace represents a complete trace with all spans
type Trace struct {
	TraceID   string            `json:"traceID"`
	Spans     []Span            `json:"spans"`
	Processes map[string]Process `json:"processes"`
}

// Span represents a single span in a trace
type Span struct {
	TraceID       string      `json:"traceID"`
	SpanID        string      `json:"spanID"`
	OperationName string      `json:"operationName"`
	References    []Reference `json:"references"`
	StartTime     int64       `json:"startTime"`     // microseconds since epoch
	Duration      int64       `json:"duration"`      // microseconds
	Tags          []Tag       `json:"tags"`
	Logs          []Log       `json:"logs"`
	ProcessID     string      `json:"processID"`
	Warnings      []string    `json:"warnings,omitempty"`
}

// Process represents service information
type Process struct {
	ServiceName string `json:"serviceName"`
	Tags        []Tag  `json:"tags"`
}

// Reference represents span relationships
type Reference struct {
	RefType string `json:"refType"` // CHILD_OF, FOLLOWS_FROM
	TraceID string `json:"traceID"`
	SpanID  string `json:"spanID"`
}

// Tag represents a key-value tag
type Tag struct {
	Key   string      `json:"key"`
	Type  string      `json:"type"` // string, int64, float64, bool
	Value interface{} `json:"value"`
}

// Log represents a log entry in a span
type Log struct {
	Timestamp int64 `json:"timestamp"` // microseconds since epoch
	Fields    []Tag `json:"fields"`
}

// ServiceInfo represents service metadata
type ServiceInfo struct {
	Name       string   `json:"name"`
	Operations []string `json:"operations"`
}

// NewJaegerDAO creates a new Jaeger DAO
func NewJaegerDAO(baseURL string, logger *zap.Logger) *JaegerDAO {
	return &JaegerDAO{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Ping checks if Jaeger is accessible
func (j *JaegerDAO) Ping(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/services", j.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := j.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Jaeger: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jaeger returned status %d", resp.StatusCode)
	}

	return nil
}

// GetTrace retrieves a single trace by ID
func (j *JaegerDAO) GetTrace(ctx context.Context, traceID string) (*Trace, error) {
	url := fmt.Sprintf("%s/api/traces/%s", j.baseURL, traceID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	j.logger.Debug("Fetching trace from Jaeger",
		zap.String("url", url),
		zap.String("trace_id", traceID),
	)

	resp, err := j.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trace: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jaeger returned status %d", resp.StatusCode)
	}

	var result struct {
		Data []Trace `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("trace not found")
	}

	j.logger.Debug("Successfully fetched trace",
		zap.String("trace_id", traceID),
		zap.Int("spans", len(result.Data[0].Spans)),
	)

	return &result.Data[0], nil
}

// SearchTraces searches for traces based on parameters
func (j *JaegerDAO) SearchTraces(ctx context.Context, params SearchParams) ([]Trace, error) {
	url := fmt.Sprintf("%s/api/traces?service=%s&limit=%d",
		j.baseURL, params.ServiceName, params.Limit)

	if params.Operation != "" {
		url += fmt.Sprintf("&operation=%s", params.Operation)
	}
	if params.MinDuration != "" {
		url += fmt.Sprintf("&minDuration=%s", params.MinDuration)
	}
	if params.MaxDuration != "" {
		url += fmt.Sprintf("&maxDuration=%s", params.MaxDuration)
	}
	if params.Start != 0 {
		url += fmt.Sprintf("&start=%d", params.Start)
	}
	if params.End != 0 {
		url += fmt.Sprintf("&end=%d", params.End)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	j.logger.Debug("Searching traces in Jaeger",
		zap.String("url", url),
		zap.String("service", params.ServiceName),
	)

	resp, err := j.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search traces: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jaeger returned status %d", resp.StatusCode)
	}

	var result struct {
		Data []Trace `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	j.logger.Debug("Successfully searched traces",
		zap.String("service", params.ServiceName),
		zap.Int("count", len(result.Data)),
	)

	return result.Data, nil
}

// GetServices retrieves list of all services
func (j *JaegerDAO) GetServices(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/api/services", j.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := j.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch services: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jaeger returned status %d", resp.StatusCode)
	}

	var result struct {
		Data []string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	j.logger.Debug("Successfully fetched services",
		zap.Int("count", len(result.Data)),
	)

	return result.Data, nil
}

// GetOperations retrieves operations for a specific service
func (j *JaegerDAO) GetOperations(ctx context.Context, serviceName string) ([]string, error) {
	url := fmt.Sprintf("%s/api/services/%s/operations", j.baseURL, serviceName)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := j.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch operations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jaeger returned status %d", resp.StatusCode)
	}

	var result struct {
		Data []struct {
			Name     string `json:"name"`
			SpanKind string `json:"spanKind"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	operations := make([]string, len(result.Data))
	for i, op := range result.Data {
		operations[i] = op.Name
	}

	j.logger.Debug("Successfully fetched operations",
		zap.String("service", serviceName),
		zap.Int("count", len(operations)),
	)

	return operations, nil
}

// SearchParams defines parameters for trace search
type SearchParams struct {
	ServiceName string
	Operation   string
	MinDuration string // e.g., "100ms"
	MaxDuration string // e.g., "1s"
	Limit       int
	Start       int64 // microseconds since epoch
	End         int64 // microseconds since epoch
	Tags        map[string]string
}

