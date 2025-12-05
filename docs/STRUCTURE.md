
# 📁 WatchingCat Repository Structure

**Clean and organized project structure**

---

## Root Directory

```
WatchingCat/
├── README.md              ⭐ Main entry point
├── docker-compose.yaml    🐳 Docker setup
├── Makefile              🔧 Build commands
├── go.mod / go.sum       📦 Go dependencies
├── Dockerfile.*          🐳 Container images
│
├── cmd/                  🚀 Applications
│   ├── backend/          Backend API server
│   └── webui/            Frontend server
│
├── internal/             🔒 Private code
│   ├── api/              API layer
│   ├── config/           Configuration
│   └── dao/              Data access
│
├── web/                  🎨 Frontend
│   ├── templates/        HTML templates
│   └── static/           CSS, JS, images
│
├── configs/              ⚙️  Configuration files
│   ├── backend-config.yaml
│   ├── otel-collector-config.yaml
│   └── prometheus.yml
│
├── k8s/                  ☸️  Kubernetes
│   ├── helm/             Helm charts
│   └── scripts/          K8s scripts
│
├── docs/                 📚 Documentation
│   ├── README.md         Documentation index
│   ├── architecture/     Architecture docs
│   ├── guides/           User guides
│   ├── kubernetes/       K8s docs
│   ├── development/      Dev guides
│   └── images/           Screenshots
│
└── scripts/              🔨 Utility scripts
```

---

## Key Files

### Root Level
- **README.md** - Project overview and quick start
- **docker-compose.yaml** - Multi-service Docker setup
- **Makefile** - Build and run commands
- **go.mod** - Go module dependencies

### Source Code
- **cmd/** - Application entry points
- **internal/** - Private Go packages
- **web/** - Frontend assets

### Configuration
- **configs/** - All configuration files
- **k8s/** - Kubernetes manifests and Helm charts

### Documentation
- **docs/** - All documentation organized by topic

---

## Quick Navigation

### For Users
- 📖 Start: [README.md](README.md)
- 🚀 Quick Start: [docs/guides/quickstart.md](docs/guides/quickstart.md)
- 📚 User Guide: [docs/guides/user-guide.md](docs/guides/user-guide.md)

### For Operators
- ☸️  K8s Deploy: [docs/kubernetes/quickstart.md](docs/kubernetes/quickstart.md)
- 🐳 Docker: [docker-compose.yaml](docker-compose.yaml)
- ⚙️  Config: [configs/](configs/)

### For Developers
- 💻 Dev Setup: [docs/development/getting-started.md](docs/development/getting-started.md)
- 🏗️  Architecture: [docs/architecture/overview.md](docs/architecture/overview.md)
- 🔨 Build: [Makefile](Makefile)

---

**Last Updated**: December 5, 2025
