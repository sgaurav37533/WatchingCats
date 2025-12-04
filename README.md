# OpenTelemetry Observability Platform

A complete, production-ready microservices observability platform built with **Go** and **OpenTelemetry**, following the [official OpenTelemetry Demo architecture](https://opentelemetry.io/docs/demo/architecture/).

## 🎯 What This Is

A fully functional e-commerce microservices application with complete observability:
- **5 Microservices** (Frontend, Cart, Product Catalog, Checkout, Load Generator)
- **Distributed Tracing** across all services
- **Structured Logging** with trace correlation
- **Real-time Metrics** and alerting
- **Exception Tracking** with stack traces
- **Complete Visualization Stack** (Jaeger, Grafana, Kibana)

## ✨ Features

### 🔍 Observability
- ✅ **Distributed Tracing** - Track requests across microservices
- ✅ **Log Management** - Structured JSON logs with trace correlation
- ✅ **Metrics Collection** - Request rates, latencies, error rates
- ✅ **Alerting System** - Real-time threshold monitoring
- ✅ **Exception Tracking** - Full stack traces with context
- ✅ **Service Mesh Visualization** - See service dependencies

### 🏗️ Architecture
- ✅ **Microservices** - 5 independent services
- ✅ **OpenTelemetry SDK** - Industry-standard instrumentation
- ✅ **OTLP Protocol** - Standard telemetry export
- ✅ **Load Generator** - Realistic traffic simulation
- ✅ **Docker Compose** - One-command deployment
- ✅ **Production-Ready** - Best practices and patterns

## 🚀 Quick Start

### Prerequisites
- **Go 1.21+** (for local development)
- **Docker & Docker Compose** (for full stack)
- **8GB RAM** recommended for Docker

### 🌐 **NEW: Web Dashboard (Recommended)**

The easiest way to use the platform:

```bash
# 1. Build all services
make build

# 2. Start the Web Dashboard
make run-webui

# 3. Open in your browser
open http://localhost:3001
```

**The Web Dashboard provides:**
- 📊 Real-time service health monitoring
- 📈 Live metrics (request rate, errors, latency)
- 📝 Recent logs with trace correlation
- 🎛️ Load generator controls
- 🔗 One-click access to Jaeger, Grafana, Prometheus, Kibana

### Option 1: Run Locally (Fastest - No Docker)

```bash
# 1. Build all services
make build

# 2. Run all services
make run-all-local

# 3. Test the frontend
curl http://localhost:8080/
```

**Services will be available at:**
- **Web Dashboard:** http://localhost:3001 ⭐️ (PRIMARY)
- Frontend API: http://localhost:8080
- Cart: http://localhost:8081
- Product Catalog: http://localhost:8082
- Checkout: http://localhost:8083

### Option 2: Full Stack with Docker

```bash
# 1. Start everything (services + backends)
make docker-up

# 2. Wait 30 seconds for services to start

# 3. Access the UIs
```

**Access Points:**
- 🌐 **Web Dashboard:** http://localhost:3001 ⭐️ (START HERE)
- **Frontend API:** http://localhost:8080
- **Jaeger (Traces):** http://localhost:16686
- **Grafana (Dashboards):** http://localhost:3000 (admin/admin)
- **Prometheus (Metrics):** http://localhost:9090
- **Kibana (Logs):** http://localhost:5601

### Option 3: Hybrid (Best for Development)

```bash
# 1. Start only backends in Docker
docker-compose up -d jaeger prometheus grafana otel-collector

# 2. Run services locally
make run-all-local
```

## 📊 Exploring the System

### View Distributed Traces (Jaeger)

1. Open http://localhost:16686
2. Select Service: `frontend`
3. Click "Find Traces"
4. Click any trace to see:
   - Complete request path
   - Timing breakdown
   - Service dependencies
   - Errors and exceptions

**What you'll see:**
```
frontend (200ms)
  ├─ GET / (5ms)
  ├─ productcatalog.ListProducts (45ms)
  │   └─ database_query (30ms)
  ├─ cartservice.GetCart (20ms)
  └─ checkoutservice.PlaceOrder (130ms)
      ├─ payment_processing (80ms)
      └─ shipping_calculation (50ms)
```

### Monitor Metrics (Grafana)

1. Open http://localhost:3000
2. Login: `admin` / `admin`
3. Explore pre-configured dashboards
4. See real-time:
   - Request rates
   - Error percentages
   - Latency percentiles (P50, P95, P99)
   - Service health

### Analyze Logs (Kibana)

1. Open http://localhost:5601
2. Create index pattern: `logs-*`
3. Search logs by:
   - Trace ID (find all logs for a request)
   - Service name
   - Log level
   - Time range

## 🏗️ Architecture

```
Load Generator → Frontend → Cart Service
                         → Product Catalog
                         → Checkout Service
                              ↓
                    All emit telemetry
                              ↓
                   OTLP Collector
                              ↓
                 ┌────────────┼────────────┐
                 ↓            ↓            ↓
              Jaeger     Prometheus  Elasticsearch
            (Traces)     (Metrics)      (Logs)
                 ↓            ↓            ↓
              Grafana  ←──────┴───────→ Kibana
           (Visualization)
```

See [DEMO_ARCHITECTURE.md](DEMO_ARCHITECTURE.md) for detailed architecture.

## 🎮 Try These Scenarios

### 1. Trace a Successful Purchase

```bash
# Make a purchase
curl -X POST http://localhost:8083/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-123",
    "email": "user@example.com",
    "address": {
      "street": "123 Main St",
      "city": "San Francisco",
      "state": "CA",
      "zip": "94102",
      "country": "USA"
    },
    "credit_card": {
      "number": "4111111111111111",
      "cvv": "123",
      "expiry": "12/25"
    }
  }'

# View the trace in Jaeger
# You'll see the complete flow through all services
```

### 2. Observe Error Handling

The checkout service has a 10% failure rate. Watch in Jaeger:
- Failed requests marked in red
- Exception details in span
- Stack traces captured
- Error propagation through services

### 3. Monitor Load

```bash
# Start load generator
make run-loadgen

# Watch in Grafana:
# - Request rate increases
# - Latency distribution
# - Error rate (should be ~10%)
```

### 4. Correlate Logs and Traces

1. Find a trace ID in Jaeger
2. Copy the trace ID
3. Search for it in Kibana
4. See all logs for that request across all services

## 📁 Project Structure

```
WatchingCat/
├── cmd/
│   ├── frontend/          # Frontend service (8080)
│   ├── cartservice/       # Cart service (8081)
│   ├── productcatalog/    # Product catalog (8082)
│   ├── checkoutservice/   # Checkout service (8083)
│   ├── loadgenerator/     # Load generator
│   └── collector/         # Custom collector
├── internal/
│   ├── tracing/           # OpenTelemetry tracing
│   ├── logging/           # Structured logging
│   ├── alerts/            # Alerting system
│   ├── exceptions/        # Exception tracking
│   └── config/            # Configuration
├── configs/
│   ├── config.yaml        # Service configuration
│   ├── otel-collector-config.yaml
│   ├── prometheus.yml
│   └── grafana-datasources.yaml
├── docker-compose.yaml    # Full stack deployment
├── Dockerfile.service     # Service container image
└── Makefile              # Build automation
```

## 🛠️ Development

### Build Services

```bash
# Build all services
make build

# Build specific service
make build-service SERVICE=frontend

# Clean and rebuild
make clean && make build
```

### Run Services Individually

```bash
# Terminal 1
make run-frontend

# Terminal 2
make run-cart

# Terminal 3
make run-product

# Terminal 4
make run-checkout

# Terminal 5
make run-loadgen
```

### View Logs

```bash
# Docker logs
make docker-logs

# Local logs (when using run-all-local)
tail -f logs/frontend.log
tail -f logs/cart.log
```

### Check Service Health

```bash
make status

# Or manually
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

## 🔧 Configuration

### Service Configuration (`configs/config.yaml`)

```yaml
tracing:
  enabled: true
  endpoint: "localhost:4317"
  sampling_rate: 1.0  # 100% sampling

logging:
  level: "info"
  format: "json"

alerts:
  enabled: true
  rules:
    - name: "high_error_rate"
      threshold: 0.05  # 5%
```

### Customize Load Generator

Edit `cmd/loadgenerator/main.go`:

```go
requestsPerMin: 30,  // Adjust traffic rate
```

## 📊 Metrics Available

### Per Service
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request latency
- `http_requests_errors_total` - Error count

### System-Wide
- Request rate (req/sec)
- Error rate (%)
- P50, P95, P99 latencies
- Service availability

## 🎓 Learning Resources

### Included Documentation
- [DEMO_ARCHITECTURE.md](DEMO_ARCHITECTURE.md) - Detailed architecture
- [QUICKSTART.md](QUICKSTART.md) - Step-by-step guide
- [EXAMPLES.md](EXAMPLES.md) - Code examples
- [ARCHITECTURE.md](ARCHITECTURE.md) - System design

### External Resources
- [OpenTelemetry Docs](https://opentelemetry.io/docs/)
- [OpenTelemetry Demo](https://opentelemetry.io/docs/demo/)
- [Jaeger Docs](https://www.jaegertracing.io/docs/)
- [Grafana Docs](https://grafana.com/docs/)

## 🐛 Troubleshooting

### Services Won't Start

```bash
# Check if ports are in use
lsof -i :8080
lsof -i :8081
lsof -i :8082
lsof -i :8083

# Kill processes if needed
kill -9 <PID>
```

### Docker Issues

```bash
# Clean up
make docker-down
docker system prune -a

# Restart
make docker-up
```

### Build Errors on macOS 26

The Makefile is already configured with `CGO_ENABLED=0` to fix this.

### No Traces in Jaeger

1. Check collector is running: `curl http://localhost:4317`
2. Check service logs for export errors
3. Wait 10-15 seconds for data to appear
4. Verify sampling rate in config.yaml

## 🚢 Deployment

### Docker Images

```bash
# Build images for all services
make docker-build

# Push to registry (customize)
docker tag otel-frontend:latest your-registry/frontend:latest
docker push your-registry/frontend:latest
```

### Kubernetes

See `k8s/` directory (coming soon) for Kubernetes manifests.

## 🤝 Contributing

This is a reference implementation. Feel free to:
- Add new services
- Implement additional features
- Improve instrumentation
- Add more backends

## 📝 License

MIT License - Use freely in your projects!

## 🎉 What's Next?

1. **Explore the traces** - See how requests flow through services
2. **Create dashboards** - Build custom Grafana dashboards
3. **Add your service** - Integrate your own microservice
4. **Customize alerts** - Set up meaningful alerts
5. **Deploy to production** - Use as a template for your system

---

**You now have a complete, production-ready OpenTelemetry observability platform!** 🚀

For questions or issues, check the documentation or create an issue.

**Happy Observing!** 👀

