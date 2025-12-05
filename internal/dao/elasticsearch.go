package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

// ElasticsearchDAO handles interactions with Elasticsearch
type ElasticsearchDAO struct {
	client *elasticsearch.Client
	index  string
	logger *zap.Logger
}

// LogEntry represents a log entry
type LogEntry struct {
	Timestamp   string                 `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Service     string                 `json:"service"`
	TraceID     string                 `json:"trace_id,omitempty"`
	SpanID      string                 `json:"span_id,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// SearchResult represents search results from Elasticsearch
type SearchResult struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			ID     string   `json:"_id"`
			Source LogEntry `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

// NewElasticsearchDAO creates a new Elasticsearch DAO
func NewElasticsearchDAO(url string, logger *zap.Logger) *ElasticsearchDAO {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Error("Failed to create Elasticsearch client", zap.Error(err))
		return &ElasticsearchDAO{
			client: nil,
			index:  "logs-*",
			logger: logger,
		}
	}

	return &ElasticsearchDAO{
		client: client,
		index:  "logs-*",
		logger: logger,
	}
}

// Ping checks if Elasticsearch is accessible
func (e *ElasticsearchDAO) Ping(ctx context.Context) error {
	if e.client == nil {
		return fmt.Errorf("elasticsearch client not initialized")
	}

	res, err := e.client.Ping(e.client.Ping.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to ping elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch returned error: %s", res.Status())
	}

	return nil
}

// Search performs a search query
func (e *ElasticsearchDAO) Search(ctx context.Context, query map[string]interface{}, from, size int) (*SearchResult, error) {
	if e.client == nil {
		return nil, fmt.Errorf("elasticsearch client not initialized")
	}

	// Build search request
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	e.logger.Debug("Executing Elasticsearch search",
		zap.String("index", e.index),
		zap.Int("from", from),
		zap.Int("size", size),
	)

	// Execute search
	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex(e.index),
		e.client.Search.WithBody(strings.NewReader(buf.String())),
		e.client.Search.WithFrom(from),
		e.client.Search.WithSize(size),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search returned error: %s", res.Status())
	}

	var result SearchResult
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	e.logger.Debug("Successfully executed search",
		zap.Int("total", result.Hits.Total.Value),
		zap.Int("returned", len(result.Hits.Hits)),
	)

	return &result, nil
}

// SearchLogs searches for logs with filters
func (e *ElasticsearchDAO) SearchLogs(ctx context.Context, params LogSearchParams) (*SearchResult, error) {
	// Build query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{},
			},
		},
		"sort": []interface{}{
			map[string]interface{}{
				"timestamp": map[string]string{
					"order": "desc",
				},
			},
		},
	}

	must := query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{})

	// Add time range
	if !params.StartTime.IsZero() || !params.EndTime.IsZero() {
		timeRange := map[string]interface{}{
			"range": map[string]interface{}{
				"timestamp": map[string]interface{}{},
			},
		}
		rangeMap := timeRange["range"].(map[string]interface{})["timestamp"].(map[string]interface{})
		
		if !params.StartTime.IsZero() {
			rangeMap["gte"] = params.StartTime.Format("2006-01-02T15:04:05.000Z")
		}
		if !params.EndTime.IsZero() {
			rangeMap["lte"] = params.EndTime.Format("2006-01-02T15:04:05.000Z")
		}
		
		must = append(must, timeRange)
	}

	// Add service filter
	if params.Service != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"service": params.Service,
			},
		})
	}

	// Add level filter
	if params.Level != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"level": params.Level,
			},
		})
	}

	// Add trace ID filter
	if params.TraceID != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"trace_id": params.TraceID,
			},
		})
	}

	// Add full-text search
	if params.Query != "" {
		must = append(must, map[string]interface{}{
			"query_string": map[string]interface{}{
				"query": params.Query,
			},
		})
	}

	query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = must

	return e.Search(ctx, query, params.From, params.Size)
}

// GetLogByTraceID retrieves all logs for a specific trace
func (e *ElasticsearchDAO) GetLogsByTraceID(ctx context.Context, traceID string) ([]LogEntry, error) {
	result, err := e.SearchLogs(ctx, LogSearchParams{
		TraceID: traceID,
		Size:    1000,
	})
	if err != nil {
		return nil, err
	}

	logs := make([]LogEntry, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		logs[i] = hit.Source
	}

	return logs, nil
}

// LogSearchParams defines parameters for log search
type LogSearchParams struct {
	Query     string
	Service   string
	Level     string
	TraceID   string
	StartTime time.Time
	EndTime   time.Time
	From      int
	Size      int
}

