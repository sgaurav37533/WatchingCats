# Quick Start Guide

Get WatchingCat running in 5 minutes!

---

## Prerequisites

- **Docker & Docker Compose** (for Docker setup)
- **Kubernetes 1.19+** and **Helm 3.x** (for K8s setup)
- **8GB RAM** recommended

---

## Option 1: Docker Compose (Recommended for Testing)

### 1. Clone and Start

```bash
# Clone repository
git clone https://github.com/yourusername/WatchingCat.git
cd WatchingCat

# Start all services
docker-compose up -d
```

### 2. Wait for Services

```bash
# Check all services are running
docker-compose ps

# Wait until all services are healthy (~2 minutes)
```

### 3. Access the UI

```bash
# Open WatchingCat UI
open http://localhost:3001
```

**Other interfaces:**
- 📊 **Grafana**: http://localhost:3000 (admin/admin)
- 🔍 **Jaeger**: http://localhost:16686
- 📈 **Prometheus**: http://localhost:9090

### 4. Generate Sample Data

```bash
# Start demo services (in new terminal)
make run-frontend
```

---

## Option 2: Kubernetes (Recommended for Production)

### 1. Install with Helm

```bash
cd WatchingCat/k8s

# Run automated installation
./scripts/install.sh

# Wait for pods to be ready
kubectl wait --for=condition=ready pod --all -n observability --timeout=300s
```

### 2. Access the UI

```bash
# Port forward to access locally
kubectl port-forward -n observability svc/watchingcat-frontend 3001:3001

# Open browser
open http://localhost:3001
```

### 3. Deploy Sample Application

```bash
# Deploy nginx with monitoring
kubectl create deployment nginx --image=nginx --replicas=2

# Annotate for Prometheus scraping
kubectl patch deployment nginx -p '{
  "spec": {
    "template": {
      "metadata": {
        "annotations": {
          "prometheus.io/scrape": "true",
          "prometheus.io/port": "80"
        }
      }
    }
  }
}'
```

📖 **[Full Kubernetes Guide →](../kubernetes/quickstart.md)**

---

## What's Next?

### Explore the UI
1. **Dashboard** - Overview of your system
2. **Services** - All discovered services
3. **Traces** - Distributed traces
4. **Metrics** - System and app metrics
5. **Topology** - Service dependency graph

### Monitor Your Application

Add OpenTelemetry instrumentation to your app:

```go
// Example: Go application
import "go.opentelemetry.io/otel"

// Set OTLP endpoint
export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
```

### Learn More
- 📖 [User Guide](user-guide.md) - Detailed usage
- 🏗️ [Architecture](../architecture/overview.md) - How it works
- ⚙️ [Configuration](configuration.md) - Customize WatchingCat

---

## Troubleshooting

### Docker: Services won't start?

```bash
# Check logs
docker-compose logs

# Restart
docker-compose restart
```

### K8s: Pods not running?

```bash
# Check pod status
kubectl get pods -n observability

# View logs
kubectl logs -n observability <pod-name>
```

### Can't access UI?

- Check if port 3001 is already in use
- Verify services are running
- Check firewall settings

---

## Next Steps

- ✅ [User Guide](user-guide.md) - Learn to use WatchingCat
- ✅ [Configuration](configuration.md) - Customize settings
- ✅ [Kubernetes Guide](../kubernetes/quickstart.md) - Deploy to K8s

**Happy monitoring!** 🐱📊

