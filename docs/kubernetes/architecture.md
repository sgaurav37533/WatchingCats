# Kubernetes Architecture

**How WatchingCat works in Kubernetes**

---

## Overview

WatchingCat uses a two-tier OpenTelemetry Collector architecture inspired by SigNoz K8s-Infra:

1. **otelAgent** (DaemonSet) - Runs on every node
2. **otelDeployment** (Deployment) - Cluster-level collection

---

## Architecture Diagram

```
┌────────────────────────────────────────────────────────┐
│              Kubernetes Cluster                        │
│                                                        │
│  ┌──────────────────────────────────────────────────┐ │
│  │ otelAgent (DaemonSet - every node)              │ │
│  │                                                  │ │
│  │ Collects:                                        │ │
│  │  • Host metrics (CPU, memory, disk, network)    │ │
│  │  • Kubelet metrics (pod/container resources)    │ │
│  │  • Container logs (Docker/Containerd/CRI-O)     │ │
│  │  • Application traces (OTLP/Jaeger/Zipkin)      │ │
│  └──────────────────────────────────────────────────┘ │
│                        │                               │
│  ┌──────────────────────────────────────────────────┐ │
│  │ otelDeployment (Deployment - cluster-wide)      │ │
│  │                                                  │ │
│  │ Collects:                                        │ │
│  │  • Cluster metrics (API server, kube-state)     │ │
│  │  • Kubernetes events                             │ │
│  │  • Prometheus scraping (annotated pods)         │ │
│  │  • OTLP gateway (receives from agents)          │ │
│  └──────────────────────────────────────────────────┘ │
│                        │                               │
└────────────────────────┼───────────────────────────────┘
                         │
                         ↓
┌────────────────────────────────────────────────────────┐
│          WatchingCat Backend (Deployment)              │
│  • Processes all telemetry                            │
│  • Unified REST API (port 8090)                       │
│  • Integrates with storage backends                   │
└────────────────────────────────────────────────────────┘
                         │
              ┌──────────┼──────────┐
              ↓          ↓          ↓
        ┌─────────┐ ┌─────────┐ ┌───────────────┐
        │ Jaeger  │ │Prometheus│ │Elasticsearch  │
        │(Traces) │ │(Metrics) │ │   (Logs)      │
        └─────────┘ └─────────┘ └───────────────┘
```

---

## Components

### otelAgent (DaemonSet)

**Why DaemonSet?**
- Needs to run on EVERY node
- Accesses host filesystem for logs
- Direct access to kubelet API
- Collects node-specific metrics

**What it collects:**
- ✅ Host metrics (CPU, memory, disk I/O, network)
- ✅ Kubelet metrics (pod/container resource usage)
- ✅ Container logs from `/var/log/pods`
- ✅ Application traces (OTLP, Jaeger, Zipkin receivers)

**Ports:**
- 4317 (OTLP gRPC)
- 4318 (OTLP HTTP)
- 14250 (Jaeger gRPC)
- 14268 (Jaeger HTTP)
- 9411 (Zipkin)
- 8888 (Prometheus metrics)

### otelDeployment (Deployment)

**Why Deployment?**
- Only needs one instance (or replicas for HA)
- Collects cluster-level data (no per-node duplication)
- API server access for cluster metrics
- Centralized Prometheus scraping

**What it collects:**
- ✅ Cluster metrics (deployments, pods, nodes)
- ✅ Kubernetes events
- ✅ Prometheus scraping from annotated pods
- ✅ Acts as OTLP gateway for aggregation

**Ports:**
- 4317 (OTLP gRPC gateway)
- 4318 (OTLP HTTP gateway)
- 8888 (Prometheus metrics)

---

## Data Flow

### 1. Application → otelAgent

```
Application (instrumented) 
  → OTLP/Jaeger/Zipkin
    → otelAgent (running on same node)
```

### 2. otelAgent → Backend

```
otelAgent 
  → Processes & enriches with K8s metadata
    → Sends via OTLP to WatchingCat Backend
```

### 3. Cluster → otelDeployment

```
Kubernetes API Server
  → k8s_cluster receiver
    → otelDeployment
```

### 4. Backend → Storage

```
WatchingCat Backend
  ├→ Jaeger (traces)
  ├→ Prometheus (metrics)
  └→ Elasticsearch (logs)
```

---

## Why Two Collectors?

### Problem: Duplication

If only DaemonSet:
- Cluster metrics collected N times (once per node)
- Events collected N times
- Wastes resources and causes duplicates

### Solution: Two-Tier Architecture

- **DaemonSet** (otelAgent): Node-specific data
- **Deployment** (otelDeployment): Cluster-wide data

**Result**: No duplication, efficient resource usage

---

## Metadata Enrichment

All telemetry is automatically enriched with:

```yaml
k8s.namespace.name: "default"
k8s.pod.name: "my-app-12345"
k8s.pod.uid: "abc-def-ghi"
k8s.node.name: "node-1"
k8s.deployment.name: "my-app"
k8s.container.name: "app"
cluster.name: "production"
cloud.region: "us-west-2"
```

---

## Resource Requirements

### Per Node (otelAgent)
- CPU: 200m (request), 500m (limit)
- Memory: 256Mi (request), 512Mi (limit)

### Cluster-wide (otelDeployment)
- CPU: 500m (request), 1000m (limit)
- Memory: 512Mi (request), 1Gi (limit)

### Backend (per replica)
- CPU: 500m (request), 1000m (limit)
- Memory: 512Mi (request), 1Gi (limit)

**Total for 3-node cluster:**
- CPU: ~4-8 cores
- Memory: ~4-8 GB

---

## Security

### RBAC

**otelAgent permissions:**
- Read: nodes, pods, services, endpoints
- Read: deployments, daemonsets, statefulsets
- Access: /metrics endpoint

**otelDeployment permissions:**
- Read: all namespace resources
- Read: cluster-level resources
- Read: events

**No write permissions** - security-first!

### Pod Security

```yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000

securityContext:
  allowPrivilegeEscalation: false
  capabilities:
    drop: [ALL]
```

---

## Scalability

### Horizontal Scaling

**otelAgent:**
- Automatically scales with nodes (DaemonSet)

**otelDeployment:**
- Can scale replicas for high load
- Use HPA based on CPU/memory

**Backend:**
- Scale replicas based on load
- Use LoadBalancer service

### Vertical Scaling

Increase resources in `values.yaml`:

```yaml
otelAgent:
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
```

---

## Monitoring the Monitors

### Collector Metrics

Both collectors expose Prometheus metrics on port 8888:

```bash
# Check otelAgent metrics
kubectl exec -n observability watchingcat-otel-agent-xxxxx -- \
  curl localhost:8888/metrics
```

### Health Checks

```bash
# Check backend health
kubectl exec -n observability watchingcat-backend-xxxxx -- \
  curl localhost:8090/health
```

---

## Next Steps

- 📖 [Helm Chart Guide](helm-chart.md) - Configure the chart
- 🚀 [Quick Start](quickstart.md) - Deploy to K8s
- ⚙️ [Production Guide](production.md) - Production setup

---

**Last Updated**: December 5, 2025

