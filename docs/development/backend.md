# WatchingCat Backend Implementation Guide

**Goal**: Build a unified backend service similar to SigNoz  
**Status**: Implementation Guide  
**Phase**: 2 (Next Steps)

---

## 🎯 Overview

This guide outlines how to build the unified WatchingCat backend service that will:
- Query telemetry data from storage backends
- Provide REST API endpoints
- Handle alert evaluation and notifications
- Serve the frontend application
- Support WebSocket for real-time updates

**Inspired by**: [SigNoz Architecture](https://signoz.io/docs/architecture/)

---

## 📁 Project Structure

```
watchingcat/
├── cmd/
│   ├── backend/
│   │   └── main.go                    // Main service entry point
│   └── migrator/
│       └── main.go                    // Database migrations
├── internal/
│   ├── api/
│   │   ├── server.go                  // HTTP server setup
│   │   ├── middleware.go              // Auth, logging, CORS
│   │   ├── router.go                  // Route definitions
│   │   ├── handlers/
│   │   │   ├── traces.go              // Trace endpoints
│   │   │   ├── metrics.go             // Metrics endpoints
│   │   │   ├── logs.go                // Logs endpoints
│   │   │   ├── services.go            // Services endpoints
│   │   │   ├── alerts.go              // Alerts endpoints
│   │   │   └── health.go              // Health checks
│   │   └── websocket/
│   │       └── hub.go                 // WebSocket hub
│   ├── query/
│   │   ├── interface.go               // Query interface
│   │   ├── jaeger/
│   │   │   ├── client.go              // Jaeger client
│   │   │   └── mapper.go              // Data mapping
│   │   ├── prometheus/
│   │   │   ├── client.go              // Prom client
│   │   │   └── query_builder.go      // Query builder
│   │   └── elasticsearch/
│   │       ├── client.go              // ES client
│   │       └── query_dsl.go          // Query DSL
│   ├── alerts/
│   │   ├── evaluator.go               // Rule evaluation
│   │   ├── manager.go                 // Alert manager
│   │   ├── rules.go                   // Rule definitions
│   │   ├── notifier/
│   │   │   ├── slack.go               // Slack notifier
│   │   │   ├── email.go               // Email notifier
│   │   │   └── webhook.go             // Webhook notifier
│   │   └── store.go                   // Alert state storage
│   ├── auth/
│   │   ├── jwt.go                     // JWT handling
│   │   ├── middleware.go              // Auth middleware
│   │   └── rbac.go                    // RBAC
│   ├── cache/
│   │   └── redis.go                   // Cache layer
│   └── models/
│       ├── trace.go                   // Trace models
│       ├── metric.go                  // Metric models
│       ├── log.go                     // Log models
│       └── alert.go                   // Alert models
├── pkg/
│   └── util/
│       ├── time.go                    // Time utilities
│       └── http.go                    // HTTP utilities
├── configs/
│   └── backend-config.yaml            // Backend configuration
└── api/
    └── openapi.yaml                   // API specification
```

---

## 🏗️ Implementation Steps

### Step 1: Project Setup

**Initialize Backend Service**:

```bash
# Create directories
mkdir -p internal/api/handlers
mkdir -p internal/query/{jaeger,prometheus,elasticsearch}
mkdir -p internal/alerts/notifier
mkdir -p internal/auth

# Create main.go
touch cmd/backend/main.go
```

**Dependencies** (`go.mod`):

```go
require (
    // HTTP Framework
    github.com/gin-gonic/gin v1.9.1
    
    // OpenTelemetry
    go.opentelemetry.io/otel v1.21.0
    go.opentelemetry.io/otel/trace v1.21.0
    go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.46.1
    
    // Clients
    github.com/jaegertracing/jaeger v1.52.0        // Jaeger client
    github.com/prometheus/client_golang v1.18.0    // Prometheus client
    github.com/elastic/go-elasticsearch/v8 v8.11.1 // Elasticsearch client
    
    // Database
    github.com/go-redis/redis/v8 v8.11.5           // Redis cache
    
    // Auth
    github.com/golang-jwt/jwt/v5 v5.2.0            // JWT
    
    // Utilities
    github.com/spf13/viper v1.18.2                 // Config
    go.uber.org/zap v1.26.0                        // Logging
    github.com/gorilla/websocket v1.5.1            // WebSocket
)
```

### Step 2: Basic HTTP Server

**`cmd/backend/main.go`**:

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gaurav/watchingcat/internal/api"
    "github.com/gaurav/watchingcat/internal/config"
    "go.uber.org/zap"
)

func main() {
    // Initialize logger
    logger, _ := zap.NewProduction()
    defer logger.Sync()

    // Load configuration
    cfg, err := config.Load("configs/backend-config.yaml")
    if err != nil {
        logger.Fatal("Failed to load config", zap.Error(err))
    }

    // Initialize API server
    router := api.NewRouter(cfg, logger)

    // Create HTTP server
    srv := &http.Server{
        Addr:    fmt.Sprintf(":%d", cfg.Port),
        Handler: router,
    }

    // Start server
    go func() {
        logger.Info("Starting WatchingCat Backend", 
            zap.Int("port", cfg.Port))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Server failed", zap.Error(err))
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Server forced to shutdown", zap.Error(err))
    }

    logger.Info("Server exited")
}
```

### Step 3: Router Setup

**`internal/api/router.go`**:

```go
package api

import (
    "github.com/gin-gonic/gin"
    "github.com/gaurav/watchingcat/internal/api/handlers"
    "github.com/gaurav/watchingcat/internal/query"
    "go.uber.org/zap"
)

func NewRouter(cfg *config.Config, logger *zap.Logger) *gin.Engine {
    router := gin.New()
    
    // Middleware
    router.Use(gin.Recovery())
    router.Use(LoggerMiddleware(logger))
    router.Use(CORSMiddleware())
    
    // Initialize query clients
    jaegerClient := query.NewJaegerClient(cfg.Jaeger)
    promClient := query.NewPrometheusClient(cfg.Prometheus)
    esClient := query.NewElasticsearchClient(cfg.Elasticsearch)
    
    // Initialize handlers
    tracesHandler := handlers.NewTracesHandler(jaegerClient, logger)
    metricsHandler := handlers.NewMetricsHandler(promClient, logger)
    logsHandler := handlers.NewLogsHandler(esClient, logger)
    servicesHandler := handlers.NewServicesHandler(jaegerClient, logger)
    
    // Static files (Frontend)
    router.Static("/static", "./web/static")
    router.LoadHTMLGlob("web/templates/*")
    
    // Root
    router.GET("/", func(c *gin.Context) {
        c.HTML(200, "index.html", gin.H{
            "JaegerURL":     cfg.Jaeger.URL,
            "GrafanaURL":    cfg.Grafana.URL,
            "PrometheusURL": cfg.Prometheus.URL,
            "KibanaURL":     cfg.Kibana.URL,
        })
    })
    
    // Health endpoints
    router.GET("/health", handlers.HealthCheck)
    router.GET("/health/ready", handlers.ReadinessCheck)
    router.GET("/health/live", handlers.LivenessCheck)
    
    // API v1
    v1 := router.Group("/api/v1")
    {
        // Traces
        traces := v1.Group("/traces")
        {
            traces.GET("", tracesHandler.ListTraces)
            traces.GET("/:id", tracesHandler.GetTrace)
            traces.POST("/search", tracesHandler.SearchTraces)
        }
        
        // Services
        services := v1.Group("/services")
        {
            services.GET("", servicesHandler.ListServices)
            services.GET("/:name", servicesHandler.GetService)
            services.GET("/:name/operations", servicesHandler.GetOperations)
            services.GET("/:name/metrics", servicesHandler.GetServiceMetrics)
        }
        
        // Metrics
        metrics := v1.Group("/metrics")
        {
            metrics.GET("", metricsHandler.GetMetrics)
            metrics.POST("/query", metricsHandler.QueryMetrics)
            metrics.GET("/labels", metricsHandler.GetLabels)
        }
        
        // Logs
        logs := v1.Group("/logs")
        {
            logs.GET("", logsHandler.GetLogs)
            logs.POST("/search", logsHandler.SearchLogs)
            logs.GET("/tail", logsHandler.TailLogs)
        }
        
        // Alerts (Future)
        alerts := v1.Group("/alerts")
        {
            alerts.GET("", handlers.ListAlerts)
            alerts.POST("", handlers.CreateAlert)
            alerts.PUT("/:id", handlers.UpdateAlert)
            alerts.DELETE("/:id", handlers.DeleteAlert)
        }
    }
    
    // WebSocket for real-time updates
    router.GET("/ws", handlers.WebSocketHandler)
    
    return router
}
```

### Step 4: Jaeger Client Implementation

**`internal/query/jaeger/client.go`**:

```go
package jaeger

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    baseURL    string
    httpClient *http.Client
}

type Trace struct {
    TraceID   string    `json:"traceID"`
    Spans     []Span    `json:"spans"`
    Processes map[string]Process `json:"processes"`
}

type Span struct {
    TraceID       string                 `json:"traceID"`
    SpanID        string                 `json:"spanID"`
    OperationName string                 `json:"operationName"`
    References    []Reference            `json:"references"`
    StartTime     int64                  `json:"startTime"`
    Duration      int64                  `json:"duration"`
    Tags          []Tag                  `json:"tags"`
    Logs          []Log                  `json:"logs"`
    ProcessID     string                 `json:"processID"`
}

type Process struct {
    ServiceName string `json:"serviceName"`
    Tags        []Tag  `json:"tags"`
}

type Tag struct {
    Key   string      `json:"key"`
    Type  string      `json:"type"`
    Value interface{} `json:"value"`
}

type Reference struct {
    RefType string `json:"refType"`
    TraceID string `json:"traceID"`
    SpanID  string `json:"spanID"`
}

type Log struct {
    Timestamp int64  `json:"timestamp"`
    Fields    []Tag  `json:"fields"`
}

func NewClient(baseURL string) *Client {
    return &Client{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *Client) GetTrace(ctx context.Context, traceID string) (*Trace, error) {
    url := fmt.Sprintf("%s/api/traces/%s", c.baseURL, traceID)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("jaeger returned status %d", resp.StatusCode)
    }
    
    var result struct {
        Data []Trace `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    if len(result.Data) == 0 {
        return nil, fmt.Errorf("trace not found")
    }
    
    return &result.Data[0], nil
}

func (c *Client) SearchTraces(ctx context.Context, serviceName string, limit int) ([]Trace, error) {
    url := fmt.Sprintf("%s/api/traces?service=%s&limit=%d", c.baseURL, serviceName, limit)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Data []Trace `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result.Data, nil
}

func (c *Client) GetServices(ctx context.Context) ([]string, error) {
    url := fmt.Sprintf("%s/api/services", c.baseURL)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Data []string `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result.Data, nil
}
```

### Step 5: API Handlers Example

**`internal/api/handlers/traces.go`**:

```go
package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gaurav/watchingcat/internal/query/jaeger"
    "go.uber.org/zap"
)

type TracesHandler struct {
    jaegerClient *jaeger.Client
    logger       *zap.Logger
}

func NewTracesHandler(client *jaeger.Client, logger *zap.Logger) *TracesHandler {
    return &TracesHandler{
        jaegerClient: client,
        logger:       logger,
    }
}

// GET /api/v1/traces
func (h *TracesHandler) ListTraces(c *gin.Context) {
    service := c.DefaultQuery("service", "frontend")
    limit := c.DefaultQuery("limit", "20")
    
    traces, err := h.jaegerClient.SearchTraces(c.Request.Context(), service, 20)
    if err != nil {
        h.logger.Error("Failed to fetch traces", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to fetch traces",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "traces": traces,
        "total":  len(traces),
    })
}

// GET /api/v1/traces/:id
func (h *TracesHandler) GetTrace(c *gin.Context) {
    traceID := c.Param("id")
    
    trace, err := h.jaegerClient.GetTrace(c.Request.Context(), traceID)
    if err != nil {
        h.logger.Error("Failed to fetch trace", 
            zap.String("traceID", traceID),
            zap.Error(err))
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Trace not found",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "trace": trace,
    })
}

// POST /api/v1/traces/search
func (h *TracesHandler) SearchTraces(c *gin.Context) {
    var req struct {
        Service   string `json:"service"`
        Operation string `json:"operation"`
        MinDuration int  `json:"minDuration"`
        MaxDuration int  `json:"maxDuration"`
        Limit     int    `json:"limit"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request",
        })
        return
    }
    
    traces, err := h.jaegerClient.SearchTraces(
        c.Request.Context(), 
        req.Service, 
        req.Limit,
    )
    if err != nil {
        h.logger.Error("Failed to search traces", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Search failed",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "traces": traces,
        "total":  len(traces),
    })
}
```

### Step 6: Configuration

**`configs/backend-config.yaml`**:

```yaml
server:
  port: 8090
  mode: release  # debug, release, test

jaeger:
  url: http://localhost:16686
  timeout: 10s

prometheus:
  url: http://localhost:9090
  timeout: 10s

elasticsearch:
  url: http://localhost:9200
  timeout: 10s
  index: logs-*

grafana:
  url: http://localhost:3000

kibana:
  url: http://localhost:5601

redis:
  host: localhost
  port: 6379
  db: 0
  password: ""

auth:
  enabled: false
  jwt_secret: "your-secret-key-change-in-production"
  token_duration: 24h

alerts:
  enabled: true
  evaluation_interval: 30s
  notification_channels:
    - type: slack
      webhook_url: ${SLACK_WEBHOOK_URL}
    - type: email
      smtp_host: smtp.gmail.com
      smtp_port: 587

cors:
  allowed_origins:
    - http://localhost:3001
    - http://localhost:3000
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
  allowed_headers:
    - Content-Type
    - Authorization

logging:
  level: info  # debug, info, warn, error
  format: json # json, console
```

### Step 7: Update WebUI to Use New API

**Update `cmd/webui/main.go`**:

```go
// Change API endpoints to use backend service
const backendURL = "http://localhost:8090"

func (w *WebUI) handleTracesAPI(rw http.ResponseWriter, r *http.Request) {
    // Proxy to unified backend
    url := backendURL + "/api/v1/traces?service=frontend&limit=20"
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        w.logger.Error("Failed to fetch traces", zap.Error(err))
        json.NewEncoder(rw).Encode(map[string]string{
            "error": "Backend not available",
        })
        return
    }
    defer resp.Body.Close()
    
    rw.Header().Set("Content-Type", "application/json")
    io.Copy(rw, resp.Body)
}
```

### Step 8: Docker Compose Integration

**Update `docker-compose.yaml`**:

```yaml
services:
  # Add WatchingCat Backend
  watchingcat-backend:
    build:
      context: .
      dockerfile: Dockerfile.backend
    ports:
      - "8090:8090"
    environment:
      - CONFIG_PATH=/app/configs/backend-config.yaml
    volumes:
      - ./configs/backend-config.yaml:/app/configs/backend-config.yaml
    depends_on:
      - otel-collector
      - jaeger
      - prometheus
      - elasticsearch
    networks:
      - otel-network

  # Update WebUI to use backend
  webui:
    # ... existing config ...
    environment:
      - BACKEND_URL=http://watchingcat-backend:8090
    depends_on:
      - watchingcat-backend
```

**Create `Dockerfile.backend`**:

```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./cmd/backend/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /backend .
COPY configs/backend-config.yaml configs/backend-config.yaml
COPY web/ ./web/

EXPOSE 8090

CMD ["./backend"]
```

---

## 🚀 Implementation Timeline

### Week 1: Foundation
- [x] Architecture design
- [ ] Project structure setup
- [ ] Basic HTTP server
- [ ] Configuration management
- [ ] Health check endpoints

### Week 2: Query Layer
- [ ] Jaeger client implementation
- [ ] Prometheus client implementation
- [ ] Elasticsearch client implementation
- [ ] Unified query interface
- [ ] API handlers

### Week 3: Integration
- [ ] Frontend API integration
- [ ] Real data in UI
- [ ] WebSocket support
- [ ] Error handling
- [ ] Testing

### Week 4: Polish
- [ ] Performance optimization
- [ ] Documentation
- [ ] Deployment guides
- [ ] Demo environment

---

## 🧪 Testing Strategy

### Unit Tests

```go
// Example test
func TestGetTrace(t *testing.T) {
    client := jaeger.NewClient("http://localhost:16686")
    
    trace, err := client.GetTrace(context.Background(), "test-trace-id")
    assert.NoError(t, err)
    assert.NotNil(t, trace)
    assert.Equal(t, "test-trace-id", trace.TraceID)
}
```

### Integration Tests

- Test end-to-end data flow
- Verify API responses
- Check error handling
- Validate data consistency

### Performance Tests

- Load testing with k6/Gatling
- Benchmark query performance
- Memory profiling
- CPU profiling

---

## 📊 Monitoring the Backend

### Self-Monitoring

The backend should instrument itself:

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/metric"
)

// Metrics
var (
    httpRequestsTotal metric.Int64Counter
    httpRequestDuration metric.Float64Histogram
    apiErrorsTotal metric.Int64Counter
)

// Initialize metrics
func initMetrics() {
    meter := otel.Meter("watchingcat-backend")
    
    httpRequestsTotal, _ = meter.Int64Counter(
        "http_requests_total",
        metric.WithDescription("Total HTTP requests"),
    )
    
    httpRequestDuration, _ = meter.Float64Histogram(
        "http_request_duration_seconds",
        metric.WithDescription("HTTP request duration"),
    )
    
    apiErrorsTotal, _ = meter.Int64Counter(
        "api_errors_total",
        metric.WithDescription("Total API errors"),
    )
}
```

---

## 🎯 Success Criteria

### Functional Requirements

- [ ] All APIs return real data (no mocks)
- [ ] <200ms API response time (p95)
- [ ] Support 1000+ concurrent users
- [ ] Handle 10K spans/second
- [ ] Zero data loss

### Non-Functional Requirements

- [ ] 99.9% uptime
- [ ] Graceful degradation
- [ ] Proper error handling
- [ ] Comprehensive logging
- [ ] Security best practices

---

## 📚 References

### Similar Projects
- [SigNoz](https://github.com/SigNoz/signoz)
- [Grafana](https://github.com/grafana/grafana)
- [Jaeger](https://github.com/jaegertracing/jaeger)

### Documentation
- [OpenTelemetry Specification](https://opentelemetry.io/docs/specs/otel/)
- [Jaeger API](https://www.jaegertracing.io/docs/latest/apis/)
- [Prometheus API](https://prometheus.io/docs/prometheus/latest/querying/api/)
- [Elasticsearch API](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html)

---

**Next Steps**: Start with Step 1 and build incrementally!

**Last Updated**: December 4, 2025

