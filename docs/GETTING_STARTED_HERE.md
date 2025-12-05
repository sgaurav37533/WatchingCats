
# 🎯 Getting Started with WatchingCat

**Start your journey here!**

---

## For Everyone: Start Here

### 1. Read the Main README
📖 **[README.md](README.md)** - Project overview, features, and quick start

This is your entry point. It has:
- What WatchingCat is
- Key features
- Quick start for Docker and Kubernetes
- Links to detailed documentation

---

## Then Choose Your Path

### 🚀 I want to try it quickly
→ **[docs/guides/quickstart.md](docs/guides/quickstart.md)**
- 5-minute Docker setup
- 5-minute K8s setup
- Get it running fast

### 📚 I want complete installation guide
→ **[docs/guides/installation.md](docs/guides/installation.md)**
- Detailed Docker setup
- Configuration options
- Troubleshooting

### ☸️ I want to deploy to Kubernetes
→ **[docs/kubernetes/quickstart.md](docs/kubernetes/quickstart.md)**
- Helm installation
- K8s architecture
- Production setup

### 🏗️ I want to understand how it works
→ **[docs/architecture/overview.md](docs/architecture/overview.md)**
- System architecture
- Component details
- Data flow

### 💻 I want to contribute
→ **[docs/development/getting-started.md](docs/development/getting-started.md)**
- Dev environment setup
- Code structure
- Development workflow

### 🗂️ I want to browse all documentation
→ **[docs/README.md](docs/README.md)**
- Complete documentation index
- All topics organized
- Quick search

---

## Documentation Structure

```
docs/
├── guides/         📚 User guides
├── kubernetes/     ☸️  K8s deployment
├── architecture/   🏗️  How it works
└── development/    💻 Contributing
```

---

## Quick Commands

### Run with Docker
```bash
docker-compose up -d
open http://localhost:3001
```

### Run with Kubernetes
```bash
cd k8s && ./scripts/install.sh
kubectl port-forward -n observability svc/watchingcat-frontend 3001:3001
open http://localhost:3001
```

### Run Backend (Development)
```bash
make run-backend
```

---

## Need Help?

- 📖 Browse: [docs/README.md](docs/README.md)
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/WatchingCat/issues)
- 💬 Discuss: [GitHub Discussions](https://github.com/yourusername/WatchingCat/discussions)

---

**Happy monitoring!** 🐱📊

