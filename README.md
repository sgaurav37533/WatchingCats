# WatchingCat 🐱 Observability Platform

<div align="center">

**An OpenTelemetry-native observability platform inspired by [SigNoz](https://signoz.io)**

*Self-hosted • Modern UI • Production-Ready • Easy to Deploy*

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-Native-blue)](https://opentelemetry.io)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker)](https://docs.docker.com/compose/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](LICENSE)

</div>

---

## 🎯 What is WatchingCat?

**WatchingCat** is a comprehensive, self-hosted observability platform that provides **unified collection, storage, and visualization** of traces, metrics, and logs—the three pillars of observability.

Built on industry-standard components with an OpenTelemetry-first approach, WatchingCat offers:

### Core Capabilities
- 🔍 **Distributed Tracing** - Track requests across microservices with inline trace viewer
- 📊 **Metrics Collection** - Real-time visualization of system and application metrics
- 📝 **Log Management** - Centralized logging with trace correlation
- 🚨 **Alert Management** - Rule-based alerting with multiple notification channels
- 📈 **Service Topology** - Interactive dependency graphs and health monitoring
- 🎨 **Modern Web UI** - Beautiful, responsive interface without framework complexity

### Demo Applications
Experience realistic telemetry with our **5-microservice e-commerce demo**:
- Frontend, Cart, Product Catalog, Checkout, Load Generator
- OpenTelemetry SDK instrumentation
- Realistic traffic patterns and error injection

## ✨ Why WatchingCat?

### 🎯 Product Vision
> **"Democratize observability by making OpenTelemetry accessible to all teams"**

Inspired by SigNoz's architecture, WatchingCat provides:

| Feature | Description | Status |
|---------|-------------|--------|
| **OpenTelemetry-Native** | Built on OTel from day one | ✅ Complete |
| **Unified Backend** | Single API for all telemetry | 🔨 Phase 2 |
| **Modern Web UI** | Beautiful interface without React | ✅ Complete |
| **Self-Hosted** | Full control of your data | ✅ Complete |
| **Easy Deployment** | Docker Compose in 5 minutes | ✅ Complete |
| **Inline Trace Viewer** | No context switching to Jaeger | ✅ Complete |
| **Real-time Metrics** | Live dashboards with Chart.js | ✅ Complete |
| **Service Topology** | Interactive D3.js graphs | ✅ Complete |
| **Alert Management** | Rule evaluation & notifications | 🔨 Phase 2 |
| **ClickHouse Support** | High-performance storage | 📅 Phase 3 |

### 🆚 Comparison with SigNoz

**Similarities:**
- ✅ OpenTelemetry-native architecture
- ✅ Unified observability (traces, metrics, logs)
- ✅ Self-hosted option
- ✅ Modern web interface
- ✅ Production-ready

**Differences:**
- **Storage**: Polyglot (Jaeger/Prometheus/ES) vs ClickHouse-only
- **Frontend**: Vanilla JS vs React
- **Focus**: Educational with demo apps vs Enterprise-first
- **Deployment**: Easier initial setup vs Advanced scaling
- **Target**: Small-medium teams vs Large-scale production

**See [WATCHINGCAT_ARCHITECTURE.md](WATCHINGCAT_ARCHITECTURE.md) for detailed comparison**

### ✨ Key Features

#### 🔍 Observability
- ✅ **Distributed Tracing** with inline viewer (OpenTelemetry-compliant)
- ✅ **Metrics Collection** - CPU, memory, network, request rates
- ✅ **Log Management** - Structured JSON logs with trace correlation
- ✅ **Alert System** - Threshold-based alerts (Phase 2)
- ✅ **Exception Tracking** - Full stack traces with context
- ✅ **Service Health** - Real-time health monitoring

#### 🏗️ Architecture
- ✅ **OpenTelemetry Collector** - Central telemetry pipeline
- ✅ **Polyglot Storage** - Jaeger (traces), Prometheus (metrics), Elasticsearch (logs)
- ✅ **Go Backend** - High-performance API server (Phase 2)
- ✅ **Modern Frontend** - Vanilla JS with Chart.js & D3.js
- ✅ **Docker Compose** - One-command deployment
- ✅ **Production-Ready** - Best practices and patterns

#### 🎨 User Interface
- ✅ **Dashboard** - System overview with key metrics
- ✅ **Services** - Health and performance monitoring
- ✅ **Traces** - Inline trace viewer with span hierarchy
- ✅ **Metrics** - Real-time charts and visualizations
- ✅ **Demo Shop** - E-commerce simulation for testing
- ✅ **Theme Support** - Light and dark modes
- ✅ **Mobile-Ready** - Responsive design

## 🚀 Quick Start

### Prerequisites
- **Go 1.21+** (for local development)
- **Docker & Docker Compose** (for full stack)
- **8GB RAM** recommended for Docker

### 🌐 **NEW: Modern Web UI (Recommended)** ⭐️

A comprehensive observability platform inspired by the [OpenTelemetry Astronomy Shop Demo](https://github.com/open-telemetry/opentelemetry-demo)!

```bash
# 1. Build all services
make build

# 2. Start the Web Dashboard
make run-webui

# 3. Open in your browser
open http://localhost:3001
```

**The Modern Web UI provides:**
- 🎨 **Beautiful Dashboard** - System overview with real-time metrics
- 🖥️ **Services Monitor** - Health and performance of all microservices
- 🔀 **Distributed Traces** - Interactive trace exploration with Jaeger integration
- 📊 **Metrics Visualization** - Charts for CPU, memory, and network
- 🛒 **Demo Shop** - Observatory-themed e-commerce for testing (6 products!)
- 📈 **Real-time Charts** - Request volume, latency percentiles with Chart.js
- 🗺️ **Service Topology** - Interactive dependency graph with D3.js
- 🌓 **Theme Support** - Light and dark mode
- 📱 **Responsive Design** - Works on desktop, tablet, and mobile
- 🎛️ **Load Generator** - Automated traffic simulation
- 🔗 **Direct Integration** - Quick access to Jaeger, Grafana, Prometheus, Kibana

**See [MODERN_UI_GUIDE.md](MODERN_UI_GUIDE.md) for the complete guide!**

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
3. Explore pre-configured dashboards:
   - **OpenTelemetry Collector Data Flow** - Monitor collector health and data pipeline
   - Custom application dashboards (create your own!)
4. See real-time:
   - Request rates
   - Error percentages
   - Latency percentiles (P50, P95, P99)
   - Service health
   - Collector throughput and export ratios

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

## 📚 Documentation

### 🎯 Product Documentation
- **[WATCHINGCAT_ARCHITECTURE.md](WATCHINGCAT_ARCHITECTURE.md)** - ⭐️ Complete architecture overview
- **[PRODUCT_ROADMAP.md](PRODUCT_ROADMAP.md)** - 🗺️ Product vision and roadmap
- **[BACKEND_IMPLEMENTATION_GUIDE.md](BACKEND_IMPLEMENTATION_GUIDE.md)** - 🔧 Building the unified backend
- **[MODERN_UI_GUIDE.md](MODERN_UI_GUIDE.md)** - 🎨 Web UI complete guide
- **[OTEL_PRINCIPLES_UPDATE.md](OTEL_PRINCIPLES_UPDATE.md)** - 📖 OpenTelemetry principles

### 🏗️ Architecture & Design
- [DEMO_ARCHITECTURE.md](DEMO_ARCHITECTURE.md) - Demo application architecture
- [ARCHITECTURE.md](ARCHITECTURE.md) - System design details
- [COLLECTOR_DASHBOARD_GUIDE.md](COLLECTOR_DASHBOARD_GUIDE.md) - Collector monitoring

### 🚀 Getting Started
- [QUICKSTART.md](QUICKSTART.md) - Step-by-step setup guide
- [GETTING_STARTED.md](GETTING_STARTED.md) - Beginner's guide
- [EXAMPLES.md](EXAMPLES.md) - Code examples
- [QUICK_REFERENCE.md](QUICK_REFERENCE.md) - Command reference

### 🔧 Implementation
- [MODERN_UI_IMPLEMENTATION.md](MODERN_UI_IMPLEMENTATION.md) - UI implementation
- [TRACE_VIEWER_FIXES.md](TRACE_VIEWER_FIXES.md) - Trace viewer details
- [WEB_UI_GUIDE.md](WEB_UI_GUIDE.md) - Web UI technical guide

### 🌐 External Resources
- [OpenTelemetry Docs](https://opentelemetry.io/docs/)
- [OpenTelemetry Demo](https://opentelemetry.io/docs/demo/)
- [SigNoz Architecture](https://signoz.io/docs/architecture/)
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

## 🗺️ Product Roadmap

### Phase 1: Foundation ✅ COMPLETE
- ✅ OpenTelemetry Collector setup
- ✅ Multi-backend storage (Jaeger, Prometheus, ES)
- ✅ 5 demo microservices with instrumentation
- ✅ Modern Web UI with inline trace viewer
- ✅ Real-time metrics and service topology
- ✅ Complete documentation

### Phase 2: Production Ready 🔨 CURRENT
- [ ] Unified Go backend service (SigNoz-style)
- [ ] Real data integration (no mocks)
- [ ] Alert Management UI
- [ ] Enhanced trace viewer with real Jaeger data
- [ ] Logs integration in UI
- [ ] WebSocket for real-time updates
- [ ] JWT authentication

**Timeline**: 2-3 weeks | [Full Roadmap →](PRODUCT_ROADMAP.md)

### Phase 3: Advanced Features 📅 Q1 2026
- [ ] ClickHouse migration (optional)
- [ ] Service Level Objectives (SLOs)
- [ ] Anomaly detection (ML-based)
- [ ] Advanced service dependency mapping
- [ ] Incident management

### Phase 4: Enterprise & Cloud 📅 Q2 2026
- [ ] Multi-tenancy support
- [ ] WatchingCat Cloud (SaaS)
- [ ] Enterprise features (SSO, RBAC)
- [ ] Advanced integrations

**See [PRODUCT_ROADMAP.md](PRODUCT_ROADMAP.md) for detailed roadmap and timelines**

## 🤝 Contributing

WatchingCat is open source and welcomes contributions!

### How to Contribute
1. **Try it out** - Use WatchingCat and provide feedback
2. **Report bugs** - Create issues for any problems
3. **Suggest features** - Share your ideas in Discussions
4. **Contribute code** - Submit pull requests
5. **Improve docs** - Help others learn

### Development Workflow
```bash
# Fork and clone
git clone https://github.com/yourusername/WatchingCat
cd WatchingCat

# Create feature branch
git checkout -b feature/amazing-feature

# Make changes and test
make build && make test

# Submit PR
git push origin feature/amazing-feature
```

### Areas We Need Help
- [ ] Real Jaeger API integration in frontend
- [ ] Alert rule builder UI
- [ ] Log viewer implementation
- [ ] Performance optimization
- [ ] Documentation improvements
- [ ] Test coverage
- [ ] Kubernetes manifests

**See [BACKEND_IMPLEMENTATION_GUIDE.md](BACKEND_IMPLEMENTATION_GUIDE.md) for Phase 2 implementation details**

## 🌟 Show Your Support

If you find WatchingCat useful, please:
- ⭐ Star this repository
- 🐦 Share on social media
- 📝 Write a blog post
- 🎥 Create a tutorial
- 💬 Join our community

## 📝 License

**Apache 2.0 License** - Use freely in your projects!

This project is inspired by and follows the principles of:
- [OpenTelemetry](https://opentelemetry.io) (Apache 2.0)
- [SigNoz](https://signoz.io) (MIT/Apache 2.0)
- [Jaeger](https://www.jaegertracing.io) (Apache 2.0)

## 🎉 What's Next?

### For Users
1. 🔍 **Explore the traces** - Click through the inline trace viewer
2. 📊 **Create dashboards** - Build custom Grafana dashboards
3. 🎨 **Customize the UI** - Modify themes and layouts
4. 🚨 **Set up alerts** - Configure meaningful alerts (Phase 2)
5. 🚀 **Deploy to production** - Use as your observability platform

### For Developers
1. 📖 **Read the architecture** - Understand the system design
2. 🔧 **Build the backend** - Help with Phase 2 implementation
3. 🧪 **Add tests** - Improve test coverage
4. 📝 **Write docs** - Help others understand
5. 🌟 **Contribute** - Submit pull requests

### For Organizations
1. 🏢 **Self-host** - Deploy in your infrastructure
2. 📈 **Scale** - Adapt for your workload
3. 🔐 **Secure** - Add enterprise authentication
4. 🎯 **Customize** - Tailor to your needs
5. 💼 **Share feedback** - Help us improve

---

<div align="center">

## 🐱 **WatchingCat: Observability Made Easy**

**Self-Hosted • OpenTelemetry-Native • Production-Ready**

[![GitHub](https://img.shields.io/badge/GitHub-WatchingCat-181717?logo=github)](.)
[![Documentation](https://img.shields.io/badge/Docs-Complete-success)](WATCHINGCAT_ARCHITECTURE.md)
[![Roadmap](https://img.shields.io/badge/Roadmap-View-blue)](PRODUCT_ROADMAP.md)

**You now have a complete observability platform inspired by SigNoz!** 🚀

For questions: Check [documentation](WATCHINGCAT_ARCHITECTURE.md) • Create [issues](../../issues) • Join [discussions](../../discussions)

**Happy Observing!** 👀📊🔍

---

*Built with ❤️ by the observability community*  
*Powered by OpenTelemetry • Inspired by SigNoz*

</div>

