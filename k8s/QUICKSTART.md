# WatchingCat Kubernetes Quick Start

**Get WatchingCat running in your K8s cluster in 5 minutes!**

---

## 🚀 5-Minute Install

### Prerequisites

- ✅ Kubernetes cluster (1.19+)
- ✅ kubectl configured
- ✅ Helm 3.x installed
- ✅ At least 4GB RAM available in cluster

### Step 1: Clone Repository

```bash
git clone https://github.com/yourusername/WatchingCat.git
cd WatchingCat/k8s
```

### Step 2: Install

```bash
# Make install script executable
chmod +x scripts/install.sh

# Run installation
./scripts/install.sh
```

**That's it!** WatchingCat is now running in your cluster.

---

## 🎯 Access the UI

### Option 1: Port Forward (Quick)

```bash
# Access Frontend UI
kubectl port-forward -n observability svc/watchingcat-frontend 3001:3001

# Open browser
open http://localhost:3001
```

### Option 2: LoadBalancer (Production)

```bash
# Get external IP
kubectl get svc -n observability watchingcat-frontend

# Access via external IP
open http://<EXTERNAL-IP>:3001
```

### Option 3: Ingress (Recommended)

Create an Ingress resource:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: watchingcat
  namespace: observability
spec:
  rules:
  - host: watchingcat.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: watchingcat-frontend
            port:
              number: 3001
```

---

## 🧪 Test It Works

### 1. Check Pods

```bash
kubectl get pods -n observability

# All should show STATUS: Running
```

### 2. Test Backend API

```bash
# Port forward
kubectl port-forward -n observability svc/watchingcat-backend 8090:8090

# Test health
curl http://localhost:8090/health | jq
```

### 3. Deploy Sample Application

```yaml
# sample-app.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-app
  namespace: default
spec:
  replicas: 2
  selector:
    matchLabels:
      app: sample-app
  template:
    metadata:
      labels:
        app: sample-app
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      containers:
      - name: app
        image: nginx:latest
        ports:
        - containerPort: 8080
        env:
        - name: OTEL_EXPORTER_OTLP_ENDPOINT
          value: "http://watchingcat-otel-agent.observability.svc.cluster.local:4317"
        - name: OTEL_SERVICE_NAME
          value: "sample-app"
```

```bash
kubectl apply -f sample-app.yaml

# Wait a minute, then check in WatchingCat UI
# You should see: sample-app in services!
```

---

## 📊 What You Get

### Monitoring Capabilities

✅ **All Nodes**: CPU, memory, disk, network metrics  
✅ **All Pods**: Resource usage and health  
✅ **All Logs**: Container logs with K8s metadata  
✅ **All Traces**: Distributed tracing across services  
✅ **Cluster Metrics**: Deployments, StatefulSets, DaemonSets  
✅ **K8s Events**: Real-time cluster events  

### Automatic Collection

- **Host Metrics**: Collected every 20s from each node
- **Kubelet Metrics**: Pod/container resources every 20s
- **Logs**: Tailed in real-time from all containers
- **Traces**: Received on OTLP/Jaeger/Zipkin ports
- **Events**: Captured as they happen

---

## 🔧 Common Tasks

### View Logs

```bash
# OTel Agent logs (DaemonSet)
kubectl logs -n observability -l app=watchingcat-otel-agent

# Backend logs
kubectl logs -n observability -l app=watchingcat-backend

# Specific pod
kubectl logs -n observability <pod-name> -f
```

### Scale Components

```bash
# Scale backend
kubectl scale deployment watchingcat-backend \
  -n observability \
  --replicas=3

# Scale OTel deployment
kubectl scale deployment watchingcat-otel-deployment \
  -n observability \
  --replicas=2
```

### Update Configuration

```bash
# Edit values
helm upgrade watchingcat ./helm/k8s-infra \
  -n observability \
  --set backend.replicas=3

# Or edit ConfigMap directly
kubectl edit cm watchingcat-otel-agent-config -n observability

# Restart pods to apply
kubectl rollout restart daemonset/watchingcat-otel-agent -n observability
```

---

## 🔐 Security Best Practices

### 1. Use Network Policies

```yaml
networkPolicy:
  enabled: true
```

### 2. Enable Pod Security

```yaml
podSecurityContext:
  runAsNonRoot: true
  runAsUser: 1000
  fsGroup: 1000
```

### 3. Restrict RBAC

The chart uses minimal read-only permissions by default.

### 4. Use Secrets for Sensitive Data

```bash
# Create secret
kubectl create secret generic watchingcat-secrets \
  -n observability \
  --from-literal=jwt-secret=your-secret-key

# Reference in values.yaml
backend:
  env:
    - name: JWT_SECRET
      valueFrom:
        secretKeyRef:
          name: watchingcat-secrets
          key: jwt-secret
```

---

## 🗑️ Uninstall

```bash
# Quick uninstall
./scripts/uninstall.sh

# Or manually
helm uninstall watchingcat -n observability
kubectl delete namespace observability
```

---

## 📚 Learn More

### Documentation
- [Full K8s Documentation](./README.md)
- [Helm Chart Values](./helm/k8s-infra/values.yaml)
- [WatchingCat Architecture](../WATCHINGCAT_ARCHITECTURE.md)

### Examples
- [Sample Applications](./examples/)
- [Custom Dashboards](./examples/dashboards/)
- [Alert Rules](./examples/alerts/)

---

## 🎯 Success Checklist

After installation, verify:

- [ ] All pods are running (`kubectl get pods -n observability`)
- [ ] Backend API is healthy (`curl http://localhost:8090/health`)
- [ ] Frontend UI is accessible (port-forward and open browser)
- [ ] OTel agents are collecting metrics (check Prometheus)
- [ ] Logs are being collected (check Elasticsearch)
- [ ] Traces work (deploy sample app and check)

---

## 💡 Tips

### Development Cluster

For Minikube/Kind:

```bash
# Start Minikube with enough resources
minikube start --cpus=4 --memory=8192

# Install WatchingCat
./scripts/install.sh

# Access via Minikube
minikube service watchingcat-frontend -n observability
```

### Production Cluster

1. Use persistent storage for databases
2. Scale backend (replicas: 3+)
3. Configure ingress with TLS
4. Set up monitoring alerts
5. Configure backup/restore

---

<div align="center">

## ✅ **WatchingCat K8s Integration Complete!**

**Your cluster is now fully observable!**

🎯 Install: `./scripts/install.sh`  
🌐 Access: Port forward port 3001  
📊 Monitor: All nodes, pods, and services  

**Happy Monitoring!** 🐱📊

</div>

---

**Last Updated**: December 5, 2025  
**Version**: 1.0.0  
**Tested on**: Kubernetes 1.19+, 1.27+

