# Getting Started - Your First Steps

Welcome! This guide will get you running your first OpenTelemetry observability platform in **5 minutes**.

## 🎯 What You'll Achieve

By the end of this guide, you'll have:
- ✅ A running OpenTelemetry collector
- ✅ An instrumented application generating telemetry
- ✅ Distributed traces visible in Jaeger
- ✅ Structured logs with trace correlation
- ✅ Active alerts monitoring error rates
- ✅ Exception tracking with stack traces

## 📋 Prerequisites Check

Before starting, verify you have:

```bash
# Check Go version (need 1.21+)
go version

# Check Docker
docker --version

# Check Docker Compose
docker-compose --version
```

If anything is missing, install it first!

## 🚀 Quick Start (5 Minutes)

### Step 1: Initialize Project (1 min)

```bash
cd /Users/gaurav/Developer/WatchingCat

# Download dependencies
go mod download

# Build binaries
make build
```

Expected output:
```
Building application...
Application built: bin/app
Building collector...
Collector built: bin/collector
```

### Step 2: Start Backend Services (2 min)

```bash
# Start Jaeger, Prometheus, Elasticsearch, etc.
docker-compose up -d

# Verify services are running
docker-compose ps
```

You should see all services with "Up" status.

### Step 3: Start the Collector (1 min)

Open a new terminal:

```bash
cd /Users/gaurav/Developer/WatchingCat
make run-collector
```

Expected output:
```
INFO  Starting OpenTelemetry Collector
INFO  Collector gRPC server starting  endpoint=0.0.0.0:4317
INFO  Collector started successfully
```

### Step 4: Start the Application (1 min)

Open another terminal:

```bash
cd /Users/gaurav/Developer/WatchingCat
make run-app
```

Expected output:
```
INFO  Starting application  service=otel-observability-platform
INFO  Tracing initialized  endpoint=localhost:4317
INFO  Starting sample workload...
INFO  Processing request  request_number=1
```

## 🎨 Explore Your Data

### View Traces in Jaeger

1. Open browser: http://localhost:16686
2. In "Service" dropdown, select: `otel-observability-platform`
3. Click "Find Traces"
4. Click on any trace to explore

**What to look for:**
- Parent span: `process_request`
- Child spans: `database_query`, `external_api_call`, `cache_lookup`
- Trace ID linking all spans
- Timing information showing where time is spent

### View Logs (Console)

In the application terminal, you'll see structured JSON logs:

```json
{
  "timestamp": "2024-12-04T10:30:15.123Z",
  "level": "info",
  "message": "Processing request",
  "trace_id": "a1b2c3d4e5f6...",
  "span_id": "1234567890ab",
  "request_number": 42
}
```

### Watch Alerts (Console)

When error rate exceeds 5%, you'll see:

```
WARN  Alert triggered
  alert_name=high_error_rate
  severity=critical
  value=0.15
  threshold=0.05
```

### View Exception Tracking (Logs)

When exceptions occur:

```
ERROR  Exception recorded
  exception_id=exc_1733310615123456789
  type=*errors.errorString
  message=simulated database connection error
  trace_id=a1b2c3d4...
```

## 📊 Understanding the Output

### Request Flow

```
1. Application creates span: "process_request"
   └─> Sets attributes: request_id, user_id, method
   
2. Child operations create nested spans:
   ├─> database_query (10-60ms)
   ├─> external_api_call (20-100ms)
   └─> cache_lookup (5-20ms)

3. Logs are emitted with trace context

4. Metrics are updated (request count, error rate)

5. If error occurs:
   ├─> Exception is tracked
   ├─> Error recorded on span
   └─> Alert may trigger
```

### Data Flow

```
Application → OTLP Exporter → Collector → Backends
                                    ↓
                         ┌──────────┼──────────┐
                         ↓          ↓          ↓
                      Jaeger   Prometheus  Elasticsearch
                    (Traces)   (Metrics)     (Logs)
```

## 🎮 Try These Actions

### 1. Increase Error Rate

Edit `cmd/app/main.go`:

```go
// Change line ~156
if rand.Float64() < 0.30 {  // Was 0.15 (15%), now 30%
```

Restart app and watch alerts fire!

### 2. View Trace Details

In Jaeger:
1. Find a trace with an error (red indicator)
2. Click to expand
3. See error message and stack trace
4. Note how trace ID links logs and traces

### 3. Check Metrics

View application metrics:

```bash
# In app terminal, every 10 requests:
INFO  Metrics snapshot
  total_requests=10
  total_errors=2
  error_rate=0.20
```

### 4. Explore Different Log Levels

Edit `configs/config.yaml`:

```yaml
logging:
  level: "debug"  # Change from "info" to see more detail
```

Restart and see debug logs about each operation.

## 🔍 Common Patterns You'll See

### Pattern 1: Successful Request

```
Trace:
  process_request (200ms)
    ├─ database_query (45ms) ✓
    ├─ external_api_call (80ms) ✓
    └─ cache_lookup (12ms) ✓

Logs:
  INFO: Processing request
  DEBUG: Database query executed
  DEBUG: External API called
  DEBUG: Cache checked
  INFO: Request completed successfully
```

### Pattern 2: Failed Request

```
Trace:
  process_request (150ms) ⚠
    ├─ database_query (40ms) ✗ ERROR
    ├─ external_api_call (skipped)
    └─ cache_lookup (skipped)

Logs:
  INFO: Processing request
  ERROR: Request failed (error: database connection error)

Exception:
  ID: exc_1234567890
  Type: *errors.errorString
  Message: simulated database connection error
  Stack trace: [10 frames]

Alert:
  Name: high_error_rate
  Value: 0.25
  Threshold: 0.05
```

## 🛠 Troubleshooting

### Issue: "Connection refused" on collector

**Solution:**
```bash
# Check if port 4317 is available
lsof -i :4317

# If occupied, kill the process
kill -9 <PID>

# Restart collector
make run-collector
```

### Issue: No traces in Jaeger

**Check:**
1. Is collector running? ✓
2. Is Jaeger running? `docker-compose ps jaeger` ✓
3. Is app using correct endpoint? Check `configs/config.yaml`
4. Wait 10-15 seconds for data to export

### Issue: Docker services won't start

**Solution:**
```bash
# Clean up
docker-compose down -v

# Remove old containers
docker system prune

# Start fresh
docker-compose up -d
```

## 📈 Next Steps

Now that you have it running:

1. **Explore the Code**
   - Read `cmd/app/main.go` - see instrumentation patterns
   - Check `internal/tracing/tracing.go` - tracing implementation
   - Review `internal/alerts/alerts.go` - alerting logic

2. **Read Documentation**
   - `README.md` - Full documentation
   - `ARCHITECTURE.md` - System design
   - `EXAMPLES.md` - Code examples

3. **Customize Configuration**
   - Edit `configs/config.yaml`
   - Adjust sampling rates
   - Change alert thresholds
   - Modify log levels

4. **Integrate Your Application**
   - Copy instrumentation patterns
   - Add tracing to your code
   - Set up logging
   - Configure alerts

5. **Deploy to Production**
   - Review security settings
   - Configure TLS
   - Set up authentication
   - Scale the collector

## 🎓 Key Concepts Learned

- ✅ **Distributed Tracing**: Following requests across services
- ✅ **Trace Context**: Correlation between logs, traces, and metrics
- ✅ **Structured Logging**: Machine-readable log format
- ✅ **Observability**: Understanding system behavior from outputs
- ✅ **Alerting**: Proactive issue detection
- ✅ **Exception Tracking**: Error monitoring with context

## 🎉 Congratulations!

You now have a fully functional OpenTelemetry observability platform! 

The best way to learn more is to:
1. Generate some load
2. Trigger some errors
3. Explore traces in Jaeger
4. Watch logs and alerts
5. Modify the code and experiment

**Happy Observing!** 👀

---

## Quick Reference

**Start Everything:**
```bash
docker-compose up -d
make run-collector  # Terminal 1
make run-app        # Terminal 2
```

**Stop Everything:**
```bash
# Ctrl+C in both terminals
docker-compose down
```

**Access UIs:**
- Jaeger: http://localhost:16686
- Grafana: http://localhost:3000 (admin/admin)
- Kibana: http://localhost:5601
- Prometheus: http://localhost:9090

**Common Commands:**
```bash
make build          # Build binaries
make test           # Run tests
make clean          # Clean build artifacts
make fmt            # Format code
make help           # Show all commands
```

Need help? Check `README.md` for detailed documentation!

