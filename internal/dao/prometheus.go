package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// PrometheusDAO handles interactions with Prometheus
type PrometheusDAO struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

// QueryResult represents a Prometheus query result
type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string         `json:"resultType"`
		Result     []MetricResult `json:"result"`
	} `json:"data"`
}

// MetricResult represents a single metric result
type MetricResult struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value,omitempty"`
	Values [][]interface{}   `json:"values,omitempty"`
}

// NewPrometheusDAO creates a new Prometheus DAO
func NewPrometheusDAO(baseURL string, logger *zap.Logger) *PrometheusDAO {
	return &PrometheusDAO{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Ping checks if Prometheus is accessible
func (p *PrometheusDAO) Ping(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/v1/query?query=up", p.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Prometheus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("prometheus returned status %d", resp.StatusCode)
	}

	return nil
}

// Query executes an instant query
func (p *PrometheusDAO) Query(ctx context.Context, query string, timestamp time.Time) (*QueryResult, error) {
	params := url.Values{}
	params.Add("query", query)
	if !timestamp.IsZero() {
		params.Add("time", fmt.Sprintf("%d", timestamp.Unix()))
	}

	url := fmt.Sprintf("%s/api/v1/query?%s", p.baseURL, params.Encode())
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	p.logger.Debug("Executing Prometheus query",
		zap.String("query", query),
	)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus returned status %d", resp.StatusCode)
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("query failed with status: %s", result.Status)
	}

	p.logger.Debug("Successfully executed query",
		zap.String("query", query),
		zap.Int("results", len(result.Data.Result)),
	)

	return &result, nil
}

// QueryRange executes a range query
func (p *PrometheusDAO) QueryRange(ctx context.Context, query string, start, end time.Time, step time.Duration) (*QueryResult, error) {
	params := url.Values{}
	params.Add("query", query)
	params.Add("start", fmt.Sprintf("%d", start.Unix()))
	params.Add("end", fmt.Sprintf("%d", end.Unix()))
	params.Add("step", fmt.Sprintf("%ds", int(step.Seconds())))

	url := fmt.Sprintf("%s/api/v1/query_range?%s", p.baseURL, params.Encode())
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	p.logger.Debug("Executing Prometheus range query",
		zap.String("query", query),
		zap.Time("start", start),
		zap.Time("end", end),
	)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus returned status %d", resp.StatusCode)
	}

	var result QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("query failed with status: %s", result.Status)
	}

	p.logger.Debug("Successfully executed range query",
		zap.String("query", query),
		zap.Int("results", len(result.Data.Result)),
	)

	return &result, nil
}

// GetLabels retrieves all label names
func (p *PrometheusDAO) GetLabels(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/labels", p.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch labels: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus returned status %d", resp.StatusCode)
	}

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// GetLabelValues retrieves values for a specific label
func (p *PrometheusDAO) GetLabelValues(ctx context.Context, label string) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/label/%s/values", p.baseURL, label)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch label values: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus returned status %d", resp.StatusCode)
	}

	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

// GetSeries retrieves time series matching label selectors
func (p *PrometheusDAO) GetSeries(ctx context.Context, matches []string, start, end time.Time) ([]map[string]string, error) {
	params := url.Values{}
	for _, match := range matches {
		params.Add("match[]", match)
	}
	if !start.IsZero() {
		params.Add("start", fmt.Sprintf("%d", start.Unix()))
	}
	if !end.IsZero() {
		params.Add("end", fmt.Sprintf("%d", end.Unix()))
	}

	url := fmt.Sprintf("%s/api/v1/series?%s", p.baseURL, params.Encode())
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch series: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("prometheus returned status %d", resp.StatusCode)
	}

	var result struct {
		Status string              `json:"status"`
		Data   []map[string]string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data, nil
}

