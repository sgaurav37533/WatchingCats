package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	Jaeger        JaegerConfig        `mapstructure:"jaeger"`
	Prometheus    PrometheusConfig    `mapstructure:"prometheus"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
	Grafana       GrafanaConfig       `mapstructure:"grafana"`
	Kibana        KibanaConfig        `mapstructure:"kibana"`
	Redis         RedisConfig         `mapstructure:"redis"`
	Auth          AuthConfig          `mapstructure:"auth"`
	Alerts        AlertsConfig        `mapstructure:"alerts"`
	CORS          CORSConfig          `mapstructure:"cors"`
	Logging       LoggingConfig       `mapstructure:"logging"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

type JaegerConfig struct {
	URL     string `mapstructure:"url"`
	Timeout string `mapstructure:"timeout"`
}

type PrometheusConfig struct {
	URL     string `mapstructure:"url"`
	Timeout string `mapstructure:"timeout"`
}

type ElasticsearchConfig struct {
	URL     string `mapstructure:"url"`
	Timeout string `mapstructure:"timeout"`
	Index   string `mapstructure:"index"`
}

type GrafanaConfig struct {
	URL string `mapstructure:"url"`
}

type KibanaConfig struct {
	URL string `mapstructure:"url"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

type AuthConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	JWTSecret     string `mapstructure:"jwt_secret"`
	TokenDuration string `mapstructure:"token_duration"`
}

type AlertsConfig struct {
	Enabled            bool                     `mapstructure:"enabled"`
	EvaluationInterval string                   `mapstructure:"evaluation_interval"`
	Channels           []NotificationChannel    `mapstructure:"notification_channels"`
}

type NotificationChannel struct {
	Type       string `mapstructure:"type"`
	WebhookURL string `mapstructure:"webhook_url,omitempty"`
	SMTPHost   string `mapstructure:"smtp_host,omitempty"`
	SMTPPort   int    `mapstructure:"smtp_port,omitempty"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, console
}

// Load reads configuration from file or environment variables
func Load() (*Config, error) {
	viper.SetConfigName("backend-config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	setDefaults()

	// Read from environment variables
	viper.AutomaticEnv()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; use defaults and env vars
			fmt.Println("No config file found, using defaults and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Override with environment variables if set
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}
	if jaegerURL := os.Getenv("JAEGER_URL"); jaegerURL != "" {
		config.Jaeger.URL = jaegerURL
	}
	if promURL := os.Getenv("PROMETHEUS_URL"); promURL != "" {
		config.Prometheus.URL = promURL
	}
	if esURL := os.Getenv("ELASTICSEARCH_URL"); esURL != "" {
		config.Elasticsearch.URL = esURL
	}

	return &config, nil
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8090)
	viper.SetDefault("server.mode", "debug")

	// Jaeger defaults
	viper.SetDefault("jaeger.url", "http://localhost:16686")
	viper.SetDefault("jaeger.timeout", "10s")

	// Prometheus defaults
	viper.SetDefault("prometheus.url", "http://localhost:9090")
	viper.SetDefault("prometheus.timeout", "10s")

	// Elasticsearch defaults
	viper.SetDefault("elasticsearch.url", "http://localhost:9200")
	viper.SetDefault("elasticsearch.timeout", "10s")
	viper.SetDefault("elasticsearch.index", "logs-*")

	// Grafana defaults
	viper.SetDefault("grafana.url", "http://localhost:3000")

	// Kibana defaults
	viper.SetDefault("kibana.url", "http://localhost:5601")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.password", "")

	// Auth defaults
	viper.SetDefault("auth.enabled", false)
	viper.SetDefault("auth.jwt_secret", "change-this-secret-in-production")
	viper.SetDefault("auth.token_duration", "24h")

	// Alerts defaults
	viper.SetDefault("alerts.enabled", false)
	viper.SetDefault("alerts.evaluation_interval", "30s")

	// CORS defaults
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3001", "http://localhost:3000"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Content-Type", "Authorization"})

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
}
