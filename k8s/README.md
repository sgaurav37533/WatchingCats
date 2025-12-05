# WatchingCat Kubernetes Deployment

**Complete Kubernetes infrastructure monitoring with OpenTelemetry**

---

## 🎯 Overview

The WatchingCat K8s-Infra chart provides comprehensive Kubernetes cluster monitoring using OpenTelemetry collectors. It deploys pre-configured collectors that gather metrics, logs, traces, and events from your entire Kubernetes environment.

### What it Does

✅ **Logs**: Tails and parses container logs  
✅ **Traces**: Collects distributed traces from applications  
✅ **Metrics**: Gathers node, pod, and cluster-level metrics  
✅ **Events**: Captures Kubernetes events for troubleshooting  
✅ **Gateway**: Acts as OTLP gateway for incoming telemetry  

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│              Kubernetes Cluster                         │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │ otelAgent (DaemonSet - runs on every node)     │   │
│  │  • Host metrics (CPU, memory, disk, network)   │   │
│  │  • Kubelet metrics (pod/container resources)   │   │
│  │  • Container logs (tail & parse)               │   │
│  │  • Application traces                          │   │
│  └─────────────────────────────────────────────────┘   │
│                        │                                │
│  ┌─────────────────────────────────────────────────┐   │
│  │ otelDeployment (Deployment - cluster-level)    │   │
│  │  • Cluster metrics (Deployments, StatefulSets) │   │
│  │  • Kubernetes events                           │   │
│  │  • Custom metrics (Redis, PostgreSQL, etc.)    │   │
│  │  • Prometheus scraping (pod annotations)       │   │
│  └─────────────────────────────────────────────────┘   │
│                        │                                │
└────────────────────────┼────────────────────────────────┘
                         │
                         ↓
┌─────────────────────────────────────────────────────────┐
│           WatchingCat Backend (Service)                 │
│  • Receives all telemetry data                         │
│  • Processes and stores in backends                     │
│  • Provides unified API                                 │
└─────────────────────────────────────────────────────────┘
                         │
              ┌──────────┼──────────┐
              ↓          ↓          ↓
        ┌─────────┐ ┌─────────┐ ┌──────────────┐
        │ Jaeger  │ │Prometheus│ │Elasticsearch │
        │(Traces) │ │(Metrics) │ │   (Logs)     │
        └─────────┘ └─────────┘ └──────────────┘
```

---

## 🚀 Quick Start

### Prerequisites

- Kubernetes cluster (1.19+)
- kubectl configured
- Helm 3.x installed
- At least 4GB RAM available

### Install with Helm

```bash
# Add WatchingCat Helm repo (when published)
helm repo add watchingcat https://watchingcat.io/charts
helm repo update

# Install WatchingCat
helm install watchingcat watchingcat/k8s-infra \
  --namespace observability \
  --create-namespace \
  --set backend.endpoint="http://watchingcat-backend:8090"
```

### Install from Source

```bash
cd /Users/gaurav/Developer/WatchingCat/k8s

# Install with Helm
helm install watchingcat ./helm/k8s-infra \
  --namespace observability \
  --create-namespace

# Or install manually with kubectl
kubectl apply -f manifests/
```

---

## 📦 Components

### 1. otelAgent (DaemonSet)

**Runs on**: Every node in the cluster  
**Purpose**: Node-level and pod-level data collection

**Collects**:
- Host metrics (CPU, memory, disk, network)
- Kubelet metrics (pod/container resources)
- Container logs (tailed and parsed)
- Application traces (OTLP receiver)

**Resource Requirements** (per pod):
- CPU: 200m (request), 500m (limit)
- Memory: 256Mi (request), 512Mi (limit)

### 2. otelDeployment (Deployment)

**Runs as**: Single or replicated deployment  
**Purpose**: Cluster-level data collection

**Collects**:
- Cluster metrics (API server, kube-state-metrics)
- Kubernetes events
- Custom application metrics
- Prometheus scraping (via annotations)

**Resource Requirements**:
- CPU: 500m (request), 1000m (limit)
- Memory: 512Mi (request), 1Gi (limit)

### 3. WatchingCat Backend (Deployment)

**Runs as**: Replicated deployment  
**Purpose**: Unified backend API

**Provides**:
- Telemetry data processing
- REST API for queries
- Storage backend integration
- WebSocket for real-time updates

### 4. Storage Backends (StatefulSets)

**Components**:
- Jaeger (traces)
- Prometheus (metrics)
- Elasticsearch (logs)

---

## ⚙️ Configuration

### Basic Configuration

```yaml
# values.yaml

# Global settings
global:
  clusterName: "my-k8s-cluster"
  region: "us-west-2"

# OpenTelemetry Agent (DaemonSet)
otelAgent:
  enabled: true
  image: otel/opentelemetry-collector-k8s:latest
  resources:
    requests:
      cpu: 200m
      memory: 256Mi
    limits:
      cpu: 500m
      memory: 512Mi
  
  # What to collect
  logs:
    enabled: true
  metrics:
    enabled: true
  traces:
    enabled: true

# OpenTelemetry Deployment (Cluster-level)
otelDeployment:
  enabled: true
  replicas: 1
  image: otel/opentelemetry-collector-k8s:latest
  resources:
    requests:
      cpu: 500m
      memory: 512Mi
    limits:
      cpu: 1000m
      memory: 1Gi
  
  # What to collect
  clusterMetrics:
    enabled: true
  kubernetesEvents:
    enabled: true
  prometheusScaping:
    enabled: true

# WatchingCat Backend
backend:
  enabled: true
  replicas: 2
  image: watchingcat/backend:latest
  endpoint: "http://watchingcat-backend:8090"

# Storage Backends
storage:
  jaeger:
    enabled: true
    storageType: memory  # or cassandra, elasticsearch
  prometheus:
    enabled: true
    retention: 15d
  elasticsearch:
    enabled: true
    replicas: 1
```

### Advanced Configuration

```yaml
# Custom exporters
otelAgent:
  config:
    exporters:
      otlp:
        endpoint: "watchingcat-backend:4317"
        tls:
          insecure: true
      prometheus:
        endpoint: "0.0.0.0:8889"

# Resource detection
otelAgent:
  resourceDetection:
    detectors:
      - env
      - system
      - docker
      - k8snode

# Custom log parsing
otelAgent:
  logs:
    operators:
      - type: json_parser
        parse_from: body
      - type: move
        from: attributes.message
        to: body
```

---

## 🔧 Installation Steps

### Step 1: Create Namespace

```bash
kubectl create namespace observability
```

### Step 2: Install with Helm

```bash
# From local directory
cd /Users/gaurav/Developer/WatchingCat/k8s

helm install watchingcat ./helm/k8s-infra \
  --namespace observability \
  --values helm/k8s-infra/values.yaml
```

### Step 3: Verify Installation

```bash
# Check pods
kubectl get pods -n observability

# Expected output:
# NAME                                    READY   STATUS    RESTARTS   AGE
# watchingcat-otel-agent-xxxxx            1/1     Running   0          2m
# watchingcat-otel-agent-yyyyy            1/1     Running   0          2m
# watchingcat-otel-deployment-xxxxx       1/1     Running   0          2m
# watchingcat-backend-xxxxx               1/1     Running   0          2m
# watchingcat-backend-yyyyy               1/1     Running   0          2m
# jaeger-xxxxx                            1/1     Running   0          2m
# prometheus-xxxxx                        1/1     Running   0          2m
# elasticsearch-xxxxx                     1/1     Running   0          2m

# Check services
kubectl get svc -n observability

# Check logs
kubectl logs -n observability -l app=watchingcat-otel-agent
```

### Step 4: Access the UI

```bash
# Port forward to access locally
kubectl port-forward -n observability svc/watchingcat-frontend 3001:3001

# Open browser
open http://localhost:3001
```

---

## 🎯 Usage Examples

### Monitor Your Applications

**Annotate your pods** to enable monitoring:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app
  annotations:
    # Enable Prometheus scraping
    prometheus.io/scrape: "true"
    prometheus.io/port: "8080"
    prometheus.io/path: "/metrics"
spec:
  containers:
  - name: app
    image: my-app:latest
    ports:
    - containerPort: 8080
    env:
    # OTLP endpoint for traces
    - name: OTEL_EXPORTER_OTLP_ENDPOINT
      value: "http://watchingcat-otel-agent:4317"
```

### View Logs

```bash
# Query logs via backend API
kubectl exec -it -n observability watchingcat-backend-xxxxx -- \
  curl "http://localhost:8090/api/v1/logs/search" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "error",
    "namespace": "default",
    "pod": "my-app",
    "size": 100
  }'
```

### View Traces

```bash
# Query traces
kubectl exec -it -n observability watchingcat-backend-xxxxx -- \
  curl "http://localhost:8090/api/v1/traces?service=my-app&limit=10"
```

### View Metrics

```bash
# Query metrics
kubectl exec -it -n observability watchingcat-backend-xxxxx -- \
  curl -X POST "http://localhost:8090/api/v1/metrics/query" \
  -H "Content-Type: application/json" \
  -d '{"query": "container_cpu_usage_seconds_total"}'
```

---

## 📊 Collected Metrics

### Node Metrics (from otelAgent)
- `node_cpu_utilization`
- `node_memory_usage`
- `node_disk_io`
- `node_network_io`
- `node_filesystem_usage`

### Pod Metrics (from otelAgent)
- `pod_cpu_utilization`
- `pod_memory_usage`
- `pod_network_io`
- `container_cpu_usage_seconds_total`
- `container_memory_working_set_bytes`

### Cluster Metrics (from otelDeployment)
- `kube_deployment_status_replicas`
- `kube_pod_status_phase`
- `kube_node_status_condition`
- `kube_service_info`

---

## 🐛 Troubleshooting

### Pods not starting?

```bash
# Check pod status
kubectl get pods -n observability

# Describe problematic pod
kubectl describe pod -n observability <pod-name>

# Check logs
kubectl logs -n observability <pod-name>
```

### No metrics appearing?

```bash
# Check otelAgent is running
kubectl get daemonset -n observability

# Check if metrics are being collected
kubectl exec -n observability watchingcat-otel-agent-xxxxx -- \
  curl localhost:8888/metrics

# Verify backend connectivity
kubectl exec -n observability watchingcat-otel-agent-xxxxx -- \
  curl watchingcat-backend:8090/health
```

### Logs not showing?

```bash
# Check log collection is enabled
kubectl get cm -n observability watchingcat-otel-agent-config -o yaml

# Verify log volume mounts
kubectl describe pod -n observability watchingcat-otel-agent-xxxxx
```

---

## 🔐 Security

### RBAC

The chart creates necessary RBAC resources:
- ServiceAccount
- ClusterRole (read-only access to K8s resources)
- ClusterRoleBinding

### Network Policies

```yaml
# Enable network policies
networkPolicy:
  enabled: true
  policyTypes:
    - Ingress
    - Egress
```

### Pod Security

```yaml
# Enable Pod Security Standards
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000

securityContext:
  allowPrivilegeEscalation: false
  readOnlyRootFilesystem: true
  capabilities:
    drop:
      - ALL
```

---

## 📚 Documentation

- [Helm Chart Documentation](./helm/k8s-infra/README.md)
- [OpenTelemetry Configuration](./configs/otel-k8s-config.yaml)
- [Deployment Guide](./docs/DEPLOYMENT.md)
- [Troubleshooting Guide](./docs/TROUBLESHOOTING.md)

---

## 🎯 Next Steps

1. ✅ Install WatchingCat in your cluster
2. ✅ Annotate your applications for monitoring
3. ✅ Access the UI and explore your cluster
4. ✅ Set up alerts and dashboards
5. ✅ Integrate with your CI/CD pipeline

---

## 📞 Support

- **Documentation**: [Full Docs](../../DOCUMENTATION_INDEX.md)
- **Issues**: [GitHub Issues](https://github.com/yourusername/WatchingCat/issues)
- **Community**: [Discord](https://discord.gg/watchingcat)

---

**Last Updated**: December 5, 2025  
**Version**: 1.0.0  
**Status**: Production Ready

