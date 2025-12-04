# Quick Reference Card

## 🚀 Start Commands

```bash
# Local (No Docker) - Fastest
make run-all-local

# Docker (Full Stack) - Complete
make docker-up

# Individual Services
make run-frontend    # Port 8080
make run-cart        # Port 8081
make run-product     # Port 8082
make run-checkout    # Port 8083
make run-loadgen     # Traffic generator
```

## 🌐 Service URLs

| Service | URL | Purpose |
|---------|-----|---------|
| Frontend | http://localhost:8080 | Web API |
| Cart | http://localhost:8081 | Cart management |
| Product Catalog | http://localhost:8082 | Products |
| Checkout | http://localhost:8083 | Orders |
| Jaeger UI | http://localhost:16686 | Traces |
| Grafana | http://localhost:3000 | Dashboards |
| Prometheus | http://localhost:9090 | Metrics |
| Kibana | http://localhost:5601 | Logs |

## 📊 Quick Tests

```bash
# Health checks
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health

# Get homepage
curl http://localhost:8080/

# List products
curl http://localhost:8082/products

# Get specific product
curl http://localhost:8082/product/OLJCESPC7Z

# Add to cart
curl -X POST http://localhost:8081/cart/add \
  -H "Content-Type: application/json" \
  -d '{"user_id":"user-1","product_id":"OLJCESPC7Z","quantity":2}'

# Get cart
curl "http://localhost:8081/cart/?user_id=user-1"
```

## 🛠️ Build Commands

```bash
make build           # Build all services
make clean           # Clean build artifacts
make test            # Run tests
make fmt             # Format code
make deps            # Download dependencies
```

## 🐳 Docker Commands

```bash
make docker-build    # Build Docker images
make docker-up       # Start all services
make docker-down     # Stop all services
make docker-logs     # View logs
make status          # Check service health
```

## 📁 Key Files

```
cmd/
├── frontend/        # Frontend service
├── cartservice/     # Cart service
├── productcatalog/  # Product catalog
├── checkoutservice/ # Checkout service
└── loadgenerator/   # Load generator

configs/
├── config.yaml      # Main configuration
├── otel-collector-config.yaml
└── prometheus.yml

bin/                 # Built binaries
logs/                # Service logs
```

## 🔍 Debugging

```bash
# View service logs (local)
tail -f logs/frontend.log
tail -f logs/cart.log

# View Docker logs
docker-compose logs -f frontend
docker-compose logs -f cartservice

# Check if port is in use
lsof -i :8080

# Kill process on port
kill -9 $(lsof -t -i:8080)
```

## 📊 Observability

### View Traces
1. Open http://localhost:16686
2. Select service: `frontend`
3. Click "Find Traces"
4. Explore trace details

### View Metrics
1. Open http://localhost:3000
2. Login: admin/admin
3. Browse dashboards:
   - **OpenTelemetry Collector Data Flow** - Monitor collector pipeline
   - Create custom application dashboards

### View Logs
1. Open http://localhost:5601
2. Create index: `logs-*`
3. Search and filter

## ⚙️ Configuration

### Change Sampling Rate
Edit `configs/config.yaml`:
```yaml
tracing:
  sampling_rate: 0.5  # 50% sampling
```

### Change Log Level
```yaml
logging:
  level: "debug"  # debug, info, warn, error
```

### Adjust Load
Edit `cmd/loadgenerator/main.go`:
```go
requestsPerMin: 60,  // Increase traffic
```

## 🐛 Common Issues

### Build Errors (macOS 26)
**Fixed!** Makefile uses `CGO_ENABLED=0`

### Port Already in Use
```bash
lsof -i :8080
kill -9 <PID>
```

### Docker Won't Start
```bash
make docker-down
docker system prune
make docker-up
```

### No Traces in Jaeger
- Wait 10-15 seconds
- Check collector: `curl http://localhost:4317`
- Verify sampling_rate in config.yaml

## 📚 Documentation

- `README.md` - Complete guide
- `DEMO_ARCHITECTURE.md` - Architecture details
- `QUICKSTART.md` - Step-by-step tutorial
- `EXAMPLES.md` - Code examples
- `COLLECTOR_DASHBOARD_GUIDE.md` - Collector monitoring
- `WHATS_NEW.md` - Recent changes

## 💡 Pro Tips

1. **Development:** Use `make run-all-local` (no Docker)
2. **Testing:** Use `make docker-up` (full stack)
3. **Debugging:** Check individual service logs
4. **Learning:** Start with one service, add more
5. **Production:** Customize config.yaml for your needs

## 🎯 Quick Demo Flow

```bash
# 1. Build
make build

# 2. Start services
make run-all-local

# 3. Generate traffic (new terminal)
make run-loadgen

# 4. View traces
open http://localhost:16686

# 5. View metrics
open http://localhost:3000

# 6. Stop (Ctrl+C in each terminal)
```

## 📞 Get Help

- Check documentation files
- Review service logs
- Test with curl commands
- Verify configuration

---

**Keep this handy for quick reference!** 📌

