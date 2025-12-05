#!/bin/bash

# WatchingCat K8s-Infra Installation Script

set -e

echo "🐱 WatchingCat Kubernetes Installation"
echo "======================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="${NAMESPACE:-observability}"
RELEASE_NAME="${RELEASE_NAME:-watchingcat}"
CHART_PATH="${CHART_PATH:-./helm/k8s-infra}"

echo "Configuration:"
echo "  Namespace: $NAMESPACE"
echo "  Release: $RELEASE_NAME"
echo "  Chart: $CHART_PATH"
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."

if ! command -v kubectl &> /dev/null; then
    echo -e "${RED}✗ kubectl not found${NC}"
    exit 1
fi
echo -e "${GREEN}✓ kubectl found${NC}"

if ! command -v helm &> /dev/null; then
    echo -e "${RED}✗ Helm not found${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Helm found${NC}"

# Check cluster connection
if ! kubectl cluster-info &> /dev/null; then
    echo -e "${RED}✗ Cannot connect to Kubernetes cluster${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Connected to Kubernetes cluster${NC}"

CLUSTER_NAME=$(kubectl config current-context)
echo -e "${GREEN}  Cluster: $CLUSTER_NAME${NC}"
echo ""

# Create namespace
echo "📦 Creating namespace..."
if kubectl get namespace $NAMESPACE &> /dev/null; then
    echo -e "${YELLOW}  Namespace $NAMESPACE already exists${NC}"
else
    kubectl create namespace $NAMESPACE
    echo -e "${GREEN}✓ Namespace $NAMESPACE created${NC}"
fi
echo ""

# Install with Helm
echo "🚀 Installing WatchingCat..."
helm install $RELEASE_NAME $CHART_PATH \
  --namespace $NAMESPACE \
  --create-namespace \
  --set global.clusterName=$CLUSTER_NAME \
  --wait \
  --timeout 10m

echo ""
echo -e "${GREEN}✓ WatchingCat installed successfully!${NC}"
echo ""

# Wait for pods
echo "⏳ Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod \
  --all \
  --namespace=$NAMESPACE \
  --timeout=300s

echo ""
echo -e "${GREEN}✓ All pods are ready!${NC}"
echo ""

# Show status
echo "📊 Deployment Status:"
echo ""
kubectl get pods -n $NAMESPACE
echo ""

# Show services
echo "🌐 Services:"
echo ""
kubectl get svc -n $NAMESPACE
echo ""

# Get access info
echo "🎯 Access Information:"
echo ""

# Backend API
BACKEND_PORT=$(kubectl get svc -n $NAMESPACE watchingcat-backend -o jsonpath='{.spec.ports[0].port}' 2>/dev/null || echo "8090")
echo "Backend API:"
echo "  kubectl port-forward -n $NAMESPACE svc/watchingcat-backend $BACKEND_PORT:$BACKEND_PORT"
echo "  Then access: http://localhost:$BACKEND_PORT"
echo ""

# Frontend UI
FRONTEND_PORT=$(kubectl get svc -n $NAMESPACE watchingcat-frontend -o jsonpath='{.spec.ports[0].port}' 2>/dev/null || echo "3001")
echo "Frontend UI:"
echo "  kubectl port-forward -n $NAMESPACE svc/watchingcat-frontend $FRONTEND_PORT:$FRONTEND_PORT"
echo "  Then access: http://localhost:$FRONTEND_PORT"
echo ""

# Jaeger
echo "Jaeger UI:"
echo "  kubectl port-forward -n $NAMESPACE svc/jaeger 16686:16686"
echo "  Then access: http://localhost:16686"
echo ""

# Grafana
echo "Grafana:"
echo "  kubectl port-forward -n $NAMESPACE svc/grafana 3000:3000"
echo "  Then access: http://localhost:3000"
echo ""

echo -e "${GREEN}✅ Installation complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Port forward to access the UI"
echo "  2. Deploy your applications with OTLP instrumentation"
echo "  3. Annotate pods for Prometheus scraping"
echo ""
echo "Documentation: https://github.com/yourusername/WatchingCat/tree/main/k8s"

