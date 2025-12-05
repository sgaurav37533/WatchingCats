# WatchingCat 🐱

**Modern Kubernetes Observability Platform powered by OpenTelemetry**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-1.19+-326CE5?logo=kubernetes)](docs/kubernetes/)
[![OpenTelemetry](https://img.shields.io/badge/OpenTelemetry-Native-blue)](https://opentelemetry.io)

WatchingCat is a lightweight, OpenTelemetry-native observability platform inspired by SigNoz. Monitor your applications and infrastructure with traces, metrics, and logs - all in one place.

---

## ✨ Features

- 🔍 **Distributed Tracing** - Visualize request flows across services
- 📊 **Metrics Collection** - Monitor system and application metrics
- 📝 **Log Aggregation** - Centralized log management with search
- 🌐 **Service Topology** - Interactive service dependency graphs
- ☸️ **Kubernetes Native** - First-class K8s support with Helm chart
- 🎯 **OpenTelemetry** - Standards-based, vendor-neutral telemetry

---

## 🚀 Quick Start

### Docker Compose (Fastest)

```bash
# Clone the repository
git clone https://github.com/yourusername/WatchingCat.git
cd WatchingCat

# Start all services
docker-compose up -d

# Access the UI
open http://localhost:3001
```

**Access Points:**
- 🎨 **UI**: http://localhost:3001
- 📊 **Grafana**: http://localhost:3000
- 🔍 **Jaeger**: http://localhost:16686
- 📈 **Prometheus**: http://localhost:9090

### Kubernetes (Production)

```bash
# Install with Helm
cd k8s
./scripts/install.sh

# Access the UI
kubectl port-forward -n observability svc/watchingcat-frontend 3001:3001
open http://localhost:3001
```

📖 **[Complete Kubernetes Guide →](docs/kubernetes/)**

---

## 📚 Documentation

### Getting Started
- **[Quick Start Guide](docs/guides/quickstart.md)** - Get running in 5 minutes
- **[Installation Guide](docs/guides/installation.md)** - Detailed setup instructions
- **[Configuration](docs/guides/configuration.md)** - Configure WatchingCat

### Kubernetes
- **[K8s Quick Start](docs/kubernetes/quickstart.md)** - Deploy to K8s in 5 minutes
- **[Helm Chart Guide](docs/kubernetes/helm-chart.md)** - Chart configuration
- **[K8s Architecture](docs/kubernetes/architecture.md)** - How it works in K8s

### Architecture
- **[System Architecture](docs/architecture/overview.md)** - High-level design
- **[Components](docs/architecture/components.md)** - Component details
- **[Data Flow](docs/architecture/data-flow.md)** - How data flows

### API & Development
- **[API Reference](docs/api/reference.md)** - REST API documentation
- **[Development Guide](docs/development/getting-started.md)** - Contribute to WatchingCat
- **[Backend Guide](docs/development/backend.md)** - Backend development

📑 **[Full Documentation Index →](docs/)**

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Applications                          │
│         (instrumented with OpenTelemetry SDK)           │
└────────────────────┬────────────────────────────────────┘
                     │ OTLP/Jaeger/Zipkin
                     ↓
┌─────────────────────────────────────────────────────────┐
│            OpenTelemetry Collector                      │
│  • Receives traces, metrics, logs                       │
│  • Processes and enriches data                          │
│  • Routes to storage backends                           │
└────────────────────┬────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        ↓            ↓            ↓
   ┌────────┐  ┌──────────┐  ┌──────────────┐
   │ Jaeger │  │Prometheus│  │Elasticsearch │
   │(Traces)│  │(Metrics) │  │   (Logs)     │
   └────┬───┘  └─────┬────┘  └──────┬───────┘
        │            │              │
        └────────────┼──────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│              WatchingCat Backend                        │
│  • Unified REST API                                     │
│  • Query interface                                      │
│  • Data aggregation                                     │
└────────────────────┬────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────────────────┐
│              WatchingCat Frontend                       │
│  • Modern React-like UI                                 │
│  • Interactive dashboards                               │
│  • Trace visualization                                  │
└─────────────────────────────────────────────────────────┘
```

---

## 🎯 Use Cases

### Development
- Debug distributed systems
- Identify performance bottlenecks
- Trace request flows
- Monitor local services

### Production
- Monitor application health
- Track SLOs and SLIs
- Incident investigation
- Capacity planning

### Kubernetes
- Monitor cluster health
- Track pod resource usage
- Debug microservices
- View service dependencies

---

## 🛠️ Tech Stack

**Backend:**
- Go (Gin framework)
- OpenTelemetry Go SDK
- Viper (configuration)

**Storage:**
- Jaeger (traces)
- Prometheus (metrics)
- Elasticsearch (logs)

**Frontend:**
- Modern JavaScript
- Chart.js (metrics visualization)
- D3.js (topology graphs)

**Infrastructure:**
- Docker & Docker Compose
- Kubernetes & Helm
- OpenTelemetry Collector

---

## 📊 Screenshots

### Dashboard
<img width="2526" height="1250" alt="image" src="https://github.com/user-attachments/assets/6050a2b2-ce1b-4d6f-981a-fa2ba5b3a004" />


### Trace Viewer
![Traces](docs/images/traces.png)

### Service Topology
![Topology](docs/images/topology.png)

---

## 🤝 Contributing

We welcome contributions! See our [Development Guide](docs/development/getting-started.md) for details.

```bash
# Fork and clone
git clone https://github.com/yourusername/WatchingCat.git

# Create a branch
git checkout -b feature/my-feature

# Make changes and commit
git commit -am "Add my feature"

# Push and create PR
git push origin feature/my-feature
```

---

## 📝 License

Apache License 2.0 - See [LICENSE](LICENSE) for details.

---

## 🙏 Acknowledgments

- Inspired by [SigNoz](https://signoz.io/)
- Built with [OpenTelemetry](https://opentelemetry.io/)
- Uses [Jaeger](https://www.jaegertracing.io/), [Prometheus](https://prometheus.io/), [Elasticsearch](https://www.elastic.co/)

---

## 📞 Support

- 📖 **Documentation**: [docs/](docs/)
- 🐛 **Issues**: [GitHub Issues](https://github.com/yourusername/WatchingCat/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/yourusername/WatchingCat/discussions)

---

<div align="center">

**Built with ❤️ for observability**

⭐ Star us on GitHub | 🐛 Report Issues | 🤝 Contribute

</div>
