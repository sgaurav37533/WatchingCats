# OpenTelemetry Demo Architecture

This implementation follows the [official OpenTelemetry Demo architecture](https://opentelemetry.io/docs/demo/architecture/) with a microservices-based e-commerce application.

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Load Generator                              │
│                    (Simulates User Traffic)                          │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             │ HTTP Requests
                             │
┌────────────────────────────▼────────────────────────────────────────┐
│                         Frontend Service                             │
│                         (Port 8080)                                  │
│                   - Homepage                                         │
│                   - Product Pages                                    │
│                   - Cart View                                        │
│                   - Checkout Flow                                    │
└────┬──────────────────┬──────────────────┬─────────────────────────┘
     │                  │                  │
     │ HTTP             │ HTTP             │ HTTP
     │                  │                  │
┌────▼──────────┐  ┌───▼──────────┐  ┌───▼──────────────┐
│  Cart Service │  │   Product    │  │    Checkout      │
│  (Port 8081)  │  │   Catalog    │  │    Service       │
│               │  │  (Port 8082) │  │  (Port 8083)     │
│ - Add Item    │  │               │  │                  │
│ - Remove Item │  │ - List        │  │ - Place Order    │
│ - Get Cart    │  │ - Search      │  │ - Payment        │
└───────┬───────┘  │ - Get Product │  │ - Shipping       │
        │          └───────┬───────┘  └──────┬───────────┘
        │                  │                  │
        │                  │                  │
        └──────────────────┴──────────────────┘
                           │
                    All services emit
                  Traces, Logs, Metrics
                           │
                           ▼
┌──────────────────────────────────────────────────────────────────────┐
│                   OpenTelemetry Collector                            │
│                        (Port 4317/4318)                              │
│                                                                      │
│  ┌────────────┐    ┌─────────────┐    ┌──────────────┐            │
│  │  Receivers │───▶│  Processors │───▶│   Exporters  │            │
│  │ OTLP gRPC  │    │   - Batch   │    │  - Jaeger    │            │
│  │ OTLP HTTP  │    │   - Filter  │    │  - Prometheus│            │
│  └────────────┘    │   - Transform│    │  - Elastic   │            │
│                    └─────────────┘    └──────────────┘            │
└──────────────────────────────────────────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌───────────────┐  ┌──────────────┐  ┌──────────────┐
│    Jaeger     │  │  Prometheus  │  │ Elasticsearch│
│  (Port 16686) │  │ (Port 9090)  │  │ (Port 9200)  │
│               │  │              │  │              │
│  Distributed  │  │   Metrics    │  │     Logs     │
│    Traces     │  │   Storage    │  │   Storage    │
└───────┬───────┘  └──────┬───────┘  └──────┬───────┘
        │                  │                  │
        │                  │                  │
        └──────────────────┴──────────────────┘
                           │
                           ▼
        ┌──────────────────────────────────────┐
        │         Visualization Layer          │
        ├──────────────────┬───────────────────┤
        │    Grafana       │      Kibana       │
        │  (Port 3000)     │   (Port 5601)     │
        │                  │                   │
        │  - Dashboards    │  - Log Analysis   │
        │  - Metrics       │  - Search         │
        │  - Alerts        │  - Visualization  │
        └──────────────────┴───────────────────┘
```

## 🎯 Services

### 1. Frontend Service (Port 8080)
**Language:** Go  
**Purpose:** Web frontend and API gateway

**Endpoints:**
- `GET /` - Homepage
- `GET /product/{id}` - Product details
- `GET /cart` - View cart
- `POST /checkout` - Checkout process
- `GET /health` - Health check

**Instrumentation:**
- HTTP request tracing
- Span creation for each endpoint
- Context propagation to downstream services
- Structured logging with trace correlation

### 2. Cart Service (Port 8081)
**Language:** Go  
**Purpose:** Shopping cart management

**Endpoints:**
- `GET /cart/?user_id={id}` - Get user's cart
- `POST /cart/add` - Add item to cart
- `POST /cart/remove` - Remove item from cart
- `GET /health` - Health check

**Features:**
- In-memory cart storage
- Thread-safe operations
- Full OpenTelemetry instrumentation

### 3. Product Catalog Service (Port 8082)
**Language:** Go  
**Purpose:** Product information and search

**Endpoints:**
- `GET /products` - List all products
- `GET /product/{id}` - Get product details
- `GET /search?q={query}` - Search products
- `GET /health` - Health check

**Features:**
- Pre-loaded product catalog
- Search functionality
- Simulated database latency

### 4. Checkout Service (Port 8083)
**Language:** Go  
**Purpose:** Order processing and payment

**Endpoints:**
- `POST /checkout` - Process order
- `GET /health` - Health check

**Features:**
- Payment processing simulation
- Exception tracking (10% failure rate)
- Order ID generation
- Full transaction tracing

### 5. Load Generator
**Language:** Go  
**Purpose:** Simulate realistic user traffic

**Features:**
- Configurable request rate (30 req/min default)
- Realistic user journeys:
  1. Browse homepage
  2. View products
  3. Add to cart
  4. Checkout (70% completion rate)
- Random delays between actions
- Multiple concurrent users

### 6. OpenTelemetry Collector
**Purpose:** Centralized telemetry collection and export

**Configuration:**
- **Receivers:** OTLP gRPC (4317), OTLP HTTP (4318)
- **Processors:** Batch, Memory Limiter, Attributes
- **Exporters:** Jaeger, Prometheus, Logging

## 📊 Observability Stack

### Tracing (Jaeger)
- **UI:** http://localhost:16686
- **Features:**
  - Service dependency graph
  - Trace search and filtering
  - Span details and timing
  - Error tracking

### Metrics (Prometheus + Grafana)
- **Prometheus:** http://localhost:9090
- **Grafana:** http://localhost:3000 (admin/admin)
- **Features:**
  - Pre-configured dashboards
  - Service metrics
  - Request rates and latencies
  - Error rates

### Logs (Elasticsearch + Kibana)
- **Kibana:** http://localhost:5601
- **Features:**
  - Log aggregation
  - Full-text search
  - Trace correlation
  - Log analysis

## 🚀 Running the Demo

### Option 1: Local Development (No Docker)

```bash
# Build all services
make build

# Run all services locally
make run-all-local

# Or run services individually in separate terminals:
make run-frontend
make run-cart
make run-product
make run-checkout
make run-loadgen
```

### Option 2: Docker Compose (Full Stack)

```bash
# Start all services including backends
make docker-up

# View logs
make docker-logs

# Check service status
make status

# Stop all services
make docker-down
```

### Option 3: Hybrid (Services local, Backends in Docker)

```bash
# Start only backends
docker-compose up -d jaeger prometheus grafana elasticsearch kibana otel-collector

# Run services locally
make run-all-local
```

## 📈 Telemetry Data Flow

### 1. Trace Flow
```
Service → OTLP Exporter → Collector → Jaeger → Grafana/Jaeger UI
```

### 2. Metrics Flow
```
Service → OTLP Exporter → Collector → Prometheus → Grafana
```

### 3. Logs Flow
```
Service → Structured Logs → Collector → Elasticsearch → Kibana
```

## 🔍 Observing the System

### View Distributed Traces

1. Open Jaeger UI: http://localhost:16686
2. Select service: `frontend`
3. Click "Find Traces"
4. Click on any trace to see:
   - Full request path across services
   - Timing breakdown
   - Span attributes
   - Errors and exceptions

### Monitor Metrics

1. Open Grafana: http://localhost:3000
2. Login: admin/admin
3. Explore dashboards:
   - Service overview
   - Request rates
   - Error rates
   - Latency percentiles

### Analyze Logs

1. Open Kibana: http://localhost:5601
2. Create index pattern: `logs-*`
3. Search logs with trace correlation
4. Filter by service, level, or trace ID

## 🎨 Key Features Demonstrated

### ✅ Distributed Tracing
- **Context Propagation:** Trace IDs flow through all services
- **Parent-Child Spans:** Clear service call hierarchy
- **Timing Analysis:** Identify slow operations
- **Error Tracking:** See exactly where failures occur

### ✅ Structured Logging
- **JSON Format:** Machine-readable logs
- **Trace Correlation:** Every log has trace_id and span_id
- **Contextual Data:** Rich metadata in every log entry
- **Log Levels:** Debug, Info, Warn, Error

### ✅ Metrics Collection
- **Request Counts:** Total requests per service
- **Error Rates:** Percentage of failed requests
- **Latency:** P50, P95, P99 percentiles
- **Custom Metrics:** Business-specific measurements

### ✅ Exception Tracking
- **Automatic Capture:** All errors recorded
- **Stack Traces:** Full call stack preserved
- **Grouping:** Similar errors grouped together
- **Trace Correlation:** Link exceptions to traces

### ✅ Service Mesh Observability
- **Service Dependencies:** Visualize service relationships
- **Health Monitoring:** Track service health
- **Load Distribution:** See traffic patterns
- **Failure Detection:** Identify problematic services

## 🔧 Configuration

### Service Configuration
All services use `configs/config.yaml`:
- Tracing endpoint
- Sampling rate
- Log level
- Alert thresholds

### Collector Configuration
`configs/otel-collector-config.yaml`:
- Receivers configuration
- Processor pipelines
- Exporter destinations

### Prometheus Configuration
`configs/prometheus.yml`:
- Scrape intervals
- Service discovery
- Alert rules

## 📝 Development

### Adding a New Service

1. Create service directory:
```bash
mkdir -p cmd/newservice
```

2. Implement service with OpenTelemetry:
```go
// See existing services for examples
```

3. Add to Makefile:
```makefile
SERVICES=frontend cartservice productcatalog checkoutservice newservice
```

4. Add to docker-compose.yaml

5. Build and run:
```bash
make build
./bin/newservice
```

## 🎯 Use Cases

This demo showcases:
1. **Microservices Observability:** Track requests across services
2. **Performance Debugging:** Identify bottlenecks
3. **Error Investigation:** Root cause analysis
4. **Capacity Planning:** Understand load patterns
5. **SLO Monitoring:** Track service level objectives

## 📚 Learn More

- [OpenTelemetry Official Demo](https://opentelemetry.io/docs/demo/)
- [OpenTelemetry Go SDK](https://opentelemetry.io/docs/instrumentation/go/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)

---

**This architecture provides a complete, production-ready example of OpenTelemetry in a microservices environment!** 🎉

