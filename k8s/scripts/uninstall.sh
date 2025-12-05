#!/bin/bash

# WatchingCat K8s-Infra Uninstallation Script

set -e

echo "🐱 WatchingCat Kubernetes Uninstallation"
echo "========================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="${NAMESPACE:-observability}"
RELEASE_NAME="${RELEASE_NAME:-watchingcat}"

echo "Configuration:"
echo "  Namespace: $NAMESPACE"
echo "  Release: $RELEASE_NAME"
echo ""

# Confirmation
read -p "Are you sure you want to uninstall WatchingCat? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "Uninstallation cancelled"
    exit 0
fi

echo ""
echo "🗑️  Uninstalling WatchingCat..."

# Uninstall Helm release
if helm list -n $NAMESPACE | grep -q $RELEASE_NAME; then
    helm uninstall $RELEASE_NAME -n $NAMESPACE
    echo -e "${GREEN}✓ Helm release uninstalled${NC}"
else
    echo -e "${YELLOW}  Release $RELEASE_NAME not found${NC}"
fi

echo ""

# Optionally delete namespace
read -p "Delete namespace $NAMESPACE? (yes/no): " delete_ns
if [ "$delete_ns" == "yes" ]; then
    kubectl delete namespace $NAMESPACE
    echo -e "${GREEN}✓ Namespace deleted${NC}"
else
    echo -e "${YELLOW}  Namespace preserved${NC}"
fi

echo ""
echo -e "${GREEN}✅ Uninstallation complete!${NC}"

