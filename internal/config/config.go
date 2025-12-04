package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Service    ServiceConfig    `yaml:"service"`
	Tracing    TracingConfig    `yaml:"tracing"`
	Logging    LoggingConfig    `yaml:"logging"`
	Alerts     AlertsConfig     `yaml:"alerts"`
	Exceptions ExceptionsConfig `yaml:"exceptions"`
	Collector  CollectorConfig  `yaml:"collector"`
	Metrics    MetricsConfig    `yaml:"metrics"`
}

// ServiceConfig contains service-level configuration
type ServiceConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
}

// TracingConfig contains tracing configuration
type TracingConfig struct {
	Enabled             bool    `yaml:"enabled"`
	Endpoint            string  `yaml:"endpoint"`
	Insecure            bool    `yaml:"insecure"`
	SamplingRate        float64 `yaml:"sampling_rate"`
	ExportInterval      string  `yaml:"export_interval"`
	MaxExportBatchSize  int     `yaml:"max_export_batch_size"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level              string `yaml:"level"`
	Format             string `yaml:"format"`
	Output             string `yaml:"output"`
	CorrelationEnabled bool   `yaml:"correlation_enabled"`
	IncludeCaller      bool   `yaml:"include_caller"`
}

// AlertsConfig contains alerting configuration
type AlertsConfig struct {
	Enabled            bool          `yaml:"enabled"`
	EvaluationInterval string        `yaml:"evaluation_interval"`
	Rules              []AlertRule   `yaml:"rules"`
	Channels           []AlertChannel `yaml:"channels"`
}

// AlertRule defines an alerting rule
type AlertRule struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Threshold   float64 `yaml:"threshold"`
	Metric      string  `yaml:"metric"`
	Window      string  `yaml:"window"`
	Severity    string  `yaml:"severity"`
}

// AlertChannel defines how alerts are sent
type AlertChannel struct {
	Type       string   `yaml:"type"`
	Enabled    bool     `yaml:"enabled"`
	URL        string   `yaml:"url,omitempty"`
	Recipients []string `yaml:"recipients,omitempty"`
}

// ExceptionsConfig contains exception tracking configuration
type ExceptionsConfig struct {
	Enabled           bool     `yaml:"enabled"`
	CaptureStackTrace bool     `yaml:"capture_stack_trace"`
	MaxStackDepth     int      `yaml:"max_stack_depth"`
	GroupByMessage    bool     `yaml:"group_by_message"`
	IgnorePatterns    []string `yaml:"ignore_patterns"`
}

// CollectorConfig contains collector configuration
type CollectorConfig struct {
	GRPC      EndpointConfig          `yaml:"grpc"`
	HTTP      EndpointConfig          `yaml:"http"`
	Exporters map[string]ExporterConfig `yaml:"exporters"`
}

// EndpointConfig defines endpoint configuration
type EndpointConfig struct {
	Endpoint string `yaml:"endpoint"`
}

// ExporterConfig defines exporter configuration
type ExporterConfig struct {
	Enabled      bool     `yaml:"enabled"`
	Endpoint     string   `yaml:"endpoint,omitempty"`
	Endpoints    []string `yaml:"endpoints,omitempty"`
	IndexPrefix  string   `yaml:"index_prefix,omitempty"`
}

// MetricsConfig contains metrics configuration
type MetricsConfig struct {
	Enabled        bool               `yaml:"enabled"`
	ExportInterval string             `yaml:"export_interval"`
	Histograms     []HistogramConfig  `yaml:"histograms"`
}

// HistogramConfig defines histogram configuration
type HistogramConfig struct {
	Name    string    `yaml:"name"`
	Buckets []float64 `yaml:"buckets"`
}

// Load loads configuration from a YAML file
func Load(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// GetExportInterval returns the tracing export interval as time.Duration
func (c *TracingConfig) GetExportInterval() time.Duration {
	d, _ := time.ParseDuration(c.ExportInterval)
	return d
}

// GetEvaluationInterval returns the alerts evaluation interval as time.Duration
func (c *AlertsConfig) GetEvaluationInterval() time.Duration {
	d, _ := time.ParseDuration(c.EvaluationInterval)
	return d
}

// GetMetricsExportInterval returns the metrics export interval as time.Duration
func (c *MetricsConfig) GetMetricsExportInterval() time.Duration {
	d, _ := time.ParseDuration(c.ExportInterval)
	return d
}

