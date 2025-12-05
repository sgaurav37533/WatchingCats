# WatchingCat Kubernetes Deployment

Deploy WatchingCat to Kubernetes with Helm.

---

## Quick Install

```bash
# Install
./scripts/install.sh

# Access UI
kubectl port-forward -n observability svc/watchingcat-frontend 3001:3001
open http://localhost:3001
```

---

## Documentation

📖 **Complete documentation**: [/docs/kubernetes/](../docs/kubernetes/)

- **[Quick Start](../docs/kubernetes/quickstart.md)** - 5-minute deployment
- **[Helm Chart](../docs/kubernetes/helm-chart.md)** - Chart configuration
- **[Architecture](../docs/kubernetes/architecture.md)** - How it works

---

## Structure

```
k8s/
├── helm/
│   └── k8s-infra/          # Helm chart
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/      # K8s manifests
└── scripts/
    ├── install.sh          # Installation
    └── uninstall.sh        # Cleanup
```

---

## Uninstall

```bash
./scripts/uninstall.sh
```

---

See [full documentation](../docs/kubernetes/) for details.

