# Quick Start Guide

Get up and running with the OpenTelemetry Observability Platform in minutes!

## Prerequisites

- **Go 1.21+**: [Download Go](https://go.dev/dl/)
- **Docker & Docker Compose**: [Install Docker](https://docs.docker.com/get-docker/)
- **Git**: For cloning the repository

## Installation Steps

### 1. Clone and Setup

```bash
# Clone the repository (if from git)
cd WatchingCat

# Run setup script
chmod +x scripts/setup.sh
./scripts/setup.sh
```

### 2. Start Backend Services

Start Jaeger, Prometheus, Elasticsearch, Kibana, and Grafana:

```bash
docker-compose up -d
```

Verify all services are running:

```bash
docker-compose ps
```

### 3. Run the Collector

In a new terminal:

```bash
make run-collector
```

You should see:

```
INFO  Collector started successfully
INFO  Collector gRPC server starting  endpoint=0.0.0.0:4317
```

### 4. Run the Sample Application

In another terminal:

```bash
make run-app
```

You should see:

```
INFO  Starting application  service=otel-observability-platform
INFO  Tracing initialized  endpoint=localhost:4317
INFO  Starting sample workload...
INFO  Processing request  request_number=1
```

## Access the UIs

Once everything is running, access the following UIs:

| Service | URL | Credentials | Purpose |
|---------|-----|-------------|---------|
| **Jaeger** | http://localhost:16686 | - | View distributed traces |
| **Grafana** | http://localhost:3000 | admin/admin | Metrics dashboards |
| **Kibana** | http://localhost:5601 | - | Log analysis |
| **Prometheus** | http://localhost:9090 | - | Metrics queries |

## Exploring the Data

### View Traces in Jaeger

1. Open http://localhost:16686
2. Select **Service**: `otel-observability-platform`
3. Click **Find Traces**
4. Click on any trace to see the detailed span hierarchy

You'll see spans like:
- `process_request` (parent span)
  - `database_query` (child span)
  - `external_api_call` (child span)
  - `cache_lookup` (child span)

### View Logs

The application logs to console with structured JSON format. Each log entry includes:
- `timestamp`: When the log was created
- `level`: Log level (info, warn, error)
- `message`: Log message
- `trace_id`: Correlation with traces
- `span_id`: Current span ID

### Monitor Alerts

Watch the console output for alerts when error rates exceed thresholds:

```
WARN  Alert triggered  
  alert_name=high_error_rate 
  severity=critical 
  value=0.15 
  threshold=0.05
```

### View Exceptions

Exceptions are automatically tracked and correlated with traces. Check logs for:

```
ERROR  Exception recorded  
  exception_id=exc_1234567890 
  type=*errors.errorString 
  message=simulated database connection error
  trace_id=abc123...
```

## Testing the System

### Generate Load

The sample app automatically generates requests. To customize:

1. Edit `cmd/app/main.go`
2. Modify the ticker interval in `runSampleWorkload()`
3. Restart the application

### Trigger Alerts

To trigger the high error rate alert:

1. Edit `cmd/app/main.go`
2. Increase the error probability: `if rand.Float64() < 0.15` → `0.30`
3. Restart the application
4. Watch for alert notifications

### Add Custom Instrumentation

Add tracing to your code:

```go
import (
    "github.com/gaurav/otel-observability/internal/tracing"
    "go.opentelemetry.io/otel/attribute"
)

func myFunction(ctx context.Context) {
    tracer := tracing.GetTracer("my-component")
    ctx, span := tracing.StartSpan(ctx, tracer, "my_operation")
    defer span.End()
    
    // Add attributes
    tracing.AddSpanAttributes(span,
        attribute.String("user.id", "12345"),
    )
    
    // Your code here...
    
    // Record errors
    if err != nil {
        tracing.RecordError(span, err)
    }
}
```

## Configuration

### Modify Settings

Edit `configs/config.yaml` to customize:

```yaml
# Change sampling rate (0.0 to 1.0)
tracing:
  sampling_rate: 0.5  # Sample 50% of traces

# Change log level
logging:
  level: "debug"  # Show debug logs

# Adjust alert thresholds
alerts:
  rules:
    - name: "high_error_rate"
      threshold: 0.10  # Alert at 10% error rate
```

### Apply Changes

Restart the application for changes to take effect:

```bash
# Stop with Ctrl+C, then:
make run-app
```

## Common Commands

```bash
# Build all binaries
make build

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean

# Download dependencies
make deps

# View all commands
make help
```

## Quick Demo Script

Run everything automatically:

```bash
chmod +x scripts/run-demo.sh
./scripts/run-demo.sh
```

This will:
1. Start all backend services
2. Launch the collector
3. Start the sample application
4. Display access URLs

Press `Ctrl+C` to stop everything.

## Troubleshooting

### Collector Won't Start

**Error**: "failed to listen: address already in use"

**Solution**: Port 4317 is already in use. Stop other OTLP collectors:
```bash
lsof -i :4317
kill -9 <PID>
```

### No Traces in Jaeger

**Check**:
1. Is the collector running?
2. Is Jaeger running? `docker-compose ps`
3. Is the app configured correctly? Check `configs/config.yaml`
4. Check collector logs for export errors

### High Memory Usage

**Solution**: Adjust batch sizes in `configs/config.yaml`:
```yaml
tracing:
  max_export_batch_size: 256  # Reduce from 512
```

### Docker Services Not Starting

**Check Docker resources**:
```bash
docker-compose down
docker system prune -a  # WARNING: Removes all unused containers
docker-compose up -d
```

## Next Steps

1. **Read the Architecture**: See `ARCHITECTURE.md` for system design
2. **Explore Examples**: Check `cmd/app/main.go` for instrumentation patterns
3. **Add Custom Components**: Extend `internal/` packages
4. **Deploy to Production**: See deployment guides in `README.md`
5. **Integrate Your App**: Use the instrumentation libraries

## Getting Help

- Check `README.md` for detailed documentation
- Review `ARCHITECTURE.md` for system design
- Look at test files (`*_test.go`) for usage examples
- Examine sample app in `cmd/app/main.go`

## Clean Up

When you're done:

```bash
# Stop application and collector (Ctrl+C in their terminals)

# Stop and remove Docker containers
docker-compose down

# Remove volumes (optional)
docker-compose down -v

# Clean build artifacts
make clean
```

---

**Congratulations!** 🎉 You now have a fully functional OpenTelemetry observability platform running locally!

