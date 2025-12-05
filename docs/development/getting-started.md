# Development Getting Started

**Set up your development environment and start contributing to WatchingCat**

---

## Prerequisites

- **Go 1.22+**
- **Node.js 18+** (for frontend tooling)
- **Docker & Docker Compose**
- **Git**
- **Make**

---

## Setup Development Environment

### 1. Fork and Clone

```bash
# Fork on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/WatchingCat.git
cd WatchingCat

# Add upstream remote
git remote add upstream https://github.com/original/WatchingCat.git
```

### 2. Install Dependencies

```bash
# Go dependencies
go mod download

# Verify installation
go version
make --version
docker --version
```

### 3. Start Backend Services

```bash
# Start Jaeger, Prometheus, Elasticsearch
docker-compose up -d jaeger prometheus elasticsearch grafana otel-collector
```

---

## Development Workflow

### Backend Development

#### Run Backend API

```bash
# Terminal 1: Run backend
make run-backend

# Backend will be available at http://localhost:8090
```

#### Run Tests

```bash
# Unit tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/api/handlers/...
```

#### Code Style

```bash
# Format code
go fmt ./...

# Lint
golangci-lint run
```

### Frontend Development

#### Run Frontend

```bash
# Terminal 2: Run frontend
make run-webui

# Frontend will be available at http://localhost:3001
```

#### Watch for Changes

The Go webui server automatically reloads templates. For static files, just refresh the browser.

---

## Project Structure

```
WatchingCat/
├── cmd/
│   ├── backend/          # Backend API server
│   │   └── main.go
│   └── webui/            # Frontend server
│       └── main.go
├── internal/
│   ├── api/              # API layer
│   │   ├── handlers/     # Request handlers
│   │   └── middleware/   # Middleware
│   ├── config/           # Configuration
│   └── dao/              # Data access
│       ├── jaeger.go
│       ├── prometheus.go
│       └── elasticsearch.go
├── web/
│   ├── templates/        # HTML templates
│   └── static/           # CSS, JS, images
│       ├── css/
│       ├── js/
│       └── images/
├── configs/              # Config files
│   ├── backend-config.yaml
│   ├── otel-collector-config.yaml
│   └── prometheus.yml
├── k8s/                  # Kubernetes manifests
│   ├── helm/
│   └── scripts/
└── docs/                 # Documentation
```

---

## Making Changes

### 1. Create a Branch

```bash
git checkout -b feature/my-feature
```

### 2. Make Changes

Follow our [coding standards](coding-standards.md).

### 3. Test Your Changes

```bash
# Run tests
go test ./...

# Manual testing
make run-backend
# Test API endpoints
```

### 4. Commit

```bash
git add .
git commit -m "Add my feature

- Implemented X
- Fixed Y
- Added tests for Z"
```

### 5. Push and Create PR

```bash
git push origin feature/my-feature
```

Then create a Pull Request on GitHub.

---

## Common Tasks

### Add a New API Endpoint

1. Add handler in `internal/api/handlers/`
2. Register route in `internal/api/router.go`
3. Add tests
4. Update API docs

### Add a New DAO Method

1. Add method in `internal/dao/`
2. Add tests
3. Use in handler

### Update Frontend

1. Modify templates in `web/templates/`
2. Update JS in `web/static/js/`
3. Update CSS in `web/static/css/`
4. Test in browser

### Add Kubernetes Feature

1. Update Helm templates in `k8s/helm/k8s-infra/templates/`
2. Update `values.yaml`
3. Test with `helm template`
4. Deploy and verify

---

## Debugging

### Backend

```bash
# Run with debugger
dlv debug cmd/backend/main.go

# Or use VS Code debugger
# (see .vscode/launch.json)
```

### Check Logs

```bash
# Backend logs
# (logs to stdout)

# Docker logs
docker-compose logs backend

# K8s logs
kubectl logs -n observability watchingcat-backend-xxxxx
```

---

## Testing

### Unit Tests

```bash
go test ./internal/...
```

### Integration Tests

```bash
# Start dependencies
docker-compose up -d

# Run integration tests
go test -tags=integration ./...
```

### End-to-End Tests

```bash
# Start full stack
docker-compose up -d

# Run e2e tests
./scripts/e2e-tests.sh
```

---

## Documentation

### Update Documentation

When adding features:
1. Update relevant docs in `docs/`
2. Update API reference if needed
3. Add examples
4. Update README if major feature

### Build Documentation

```bash
# (If we add doc generation later)
make docs
```

---

## Contribution Guidelines

### Code Review

- All PRs require review
- Address review comments
- Keep PRs focused and small
- Write clear commit messages

### Testing

- Add tests for new features
- Maintain test coverage >70%
- Test edge cases

### Documentation

- Document all public APIs
- Add comments for complex logic
- Update docs with code changes

---

## Getting Help

- 📖 [Backend Guide](backend.md) - Detailed backend docs
- 🎨 [Frontend Guide](frontend.md) - Frontend development
- 💬 [GitHub Discussions](https://github.com/yourusername/WatchingCat/discussions)
- 🐛 [Report Issues](https://github.com/yourusername/WatchingCat/issues)

---

## Next Steps

- 📖 [Backend Development](backend.md) - Backend details
- 🎨 [Frontend Development](frontend.md) - Frontend details
- 🧪 [Testing Guide](testing.md) - Testing practices

**Happy coding!** 🚀

