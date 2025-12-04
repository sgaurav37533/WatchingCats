# Web UI Dashboard Guide

## 🌐 Overview

The **Web UI Dashboard** is your central control panel for the entire OpenTelemetry observability platform. It provides a unified interface to monitor, control, and visualize all your microservices and telemetry data in real-time.

## 🚀 Quick Start

### Access the Dashboard

```bash
# Option 1: Run locally
make run-webui

# Then open in browser
open http://localhost:3001
```

**Or:**

```bash
# Option 2: With all services
make run-all-local

# In another terminal
./bin/webui
```

### What You'll See

When you open http://localhost:3001, you'll see:

1. **Service Status** - Real-time health of all microservices
2. **Metrics Overview** - Request rates, errors, latency
3. **Observability Tools** - Quick links to Jaeger, Grafana, Prometheus, Kibana
4. **Recent Logs** - Live log stream across all services
5. **Load Generator Control** - Start/stop traffic generation

## 📊 Dashboard Features

### 1. Microservices Status

**Real-time health monitoring** for all services:

- ✅ **Green Card** = Service healthy
- ❌ **Red Card** = Service down or unhealthy
- 🔄 **Auto-refresh** every 5 seconds

**Services monitored:**
- Frontend (Port 8080)
- Cart Service (Port 8081)
- Product Catalog (Port 8082)
- Checkout Service (Port 8083)

### 2. Real-Time Metrics

**Live metrics updated every 5 seconds:**

- **Request Rate** - Requests per second
- **Error Rate** - Percentage of failed requests
- **Avg Latency** - Average response time (ms)
- **P95 Latency** - 95th percentile latency (ms)

**Color coding:**
- Blue = Normal metrics
- Red = Error metrics

### 3. Observability Tools

**One-click access** to all observability backends:

#### Jaeger (Distributed Tracing)
- View trace spans
- Analyze service dependencies
- Debug performance issues
- **Click**: Opens in new tab

#### Grafana (Dashboards)
- Pre-configured dashboards
- Custom visualizations
- Alerting rules
- **Click**: Opens in new tab

#### Prometheus (Metrics)
- Query metrics
- Create custom queries
- View time-series data
- **Click**: Opens in new tab

#### Kibana (Logs)
- Search logs
- Create visualizations
- Analyze patterns
- **Click**: Opens in new tab

### 4. Recent Logs

**Live log stream** showing:
- Timestamp
- Log level (info, warn, error)
- Service name
- Message
- Trace ID (for correlation)

**Color coded by level:**
- 🟢 **Green** = Info
- 🟡 **Yellow** = Warning
- 🔴 **Red** = Error

### 5. Load Generator Control

**Control traffic generation:**

**Start Load Generator:**
- Click "Start Load Generator"
- Simulates 30 users/minute
- Realistic user journeys
- Creates traces, logs, metrics

**Stop Load Generator:**
- Click "Stop Load Generator"
- Stops all simulated traffic

## 🎨 Dashboard Layout

```
┌─────────────────────────────────────────────┐
│  Header                                     │
│  OpenTelemetry Observability Platform       │
│  [Refresh Button]                           │
├─────────────────────────────────────────────┤
│  Microservices Status                       │
│  [Frontend] [Cart] [Product] [Checkout]     │
├─────────────────────────────────────────────┤
│  Real-Time Metrics                          │
│  [Req Rate] [Err Rate] [Latency] [P95]     │
├─────────────────────────────────────────────┤
│  Observability Tools                        │
│  [Jaeger] [Grafana] [Prometheus] [Kibana]  │
├─────────────────────────────────────────────┤
│  Recent Logs                                │
│  [Live log stream...]                       │
├─────────────────────────────────────────────┤
│  Load Generator                             │
│  Status: [Running/Stopped]                  │
│  [Start] [Stop]                             │
└─────────────────────────────────────────────┘
```

## 🔄 Auto-Refresh

The dashboard **automatically updates** every 5 seconds:
- Service status checks
- Metrics refresh
- Log stream updates
- Timestamp updates

**Manual refresh:**
- Click the "Refresh" button in header
- All data refreshes immediately

## 📱 Responsive Design

The dashboard is **fully responsive**:
- Desktop: Full layout with all features
- Tablet: Adapted grid layout
- Mobile: Stacked vertical layout

## 🎯 Common Workflows

### 1. Check System Health

1. Open dashboard
2. Look at **Microservices Status** section
3. All green = healthy
4. Any red = investigate that service

### 2. Monitor Performance

1. Watch **Real-Time Metrics**
2. Check request rate and latency
3. Monitor error rate
4. If metrics spike, investigate in Jaeger

### 3. Debug Issues

1. See error in logs
2. Copy the trace ID
3. Click **Jaeger** tool card
4. Search for trace ID
5. Analyze the full trace

### 4. Generate Load for Testing

1. Scroll to **Load Generator** section
2. Click "Start Load Generator"
3. Watch metrics increase
4. View traces in Jaeger
5. Stop when done testing

### 5. View Traces

1. Click **Jaeger** card
2. Select service (e.g., frontend)
3. Click "Find Traces"
4. Explore trace details

### 6. Create Dashboards

1. Click **Grafana** card
2. Login: admin/admin
3. Create new dashboard
4. Add panels with metrics
5. Save dashboard

### 7. Search Logs

1. Click **Kibana** card
2. Create index pattern: `logs-*`
3. Search by trace ID, service, or message
4. Filter and analyze

## 🔧 Configuration

### API Endpoints

The dashboard uses these internal APIs:

- `GET /api/services` - Service health status
- `GET /api/metrics` - Real-time metrics
- `GET /api/logs` - Recent log entries
- `POST /api/loadgen/start` - Start load generator
- `POST /api/loadgen/stop` - Stop load generator
- `GET /health` - Dashboard health check

### Customization

**Change refresh interval:**

Edit `web/static/js/dashboard.js`:

```javascript
// Change from 5000ms (5 seconds) to desired interval
updateInterval = setInterval(() => {
    // ...
}, 5000);
```

**Change service URLs:**

Edit `cmd/webui/main.go`:

```go
frontendURL:    "http://localhost:8080",
cartURL:        "http://localhost:8081",
// ...
```

## 🐛 Troubleshooting

### Dashboard Won't Load

**Check if service is running:**
```bash
curl http://localhost:3001/health
```

**If not running:**
```bash
make run-webui
```

### Services Show as Unhealthy

**Verify services are running:**
```bash
make status
```

**Start missing services:**
```bash
make run-frontend
make run-cart
make run-product
make run-checkout
```

### Metrics Not Updating

**Check API:**
```bash
curl http://localhost:3001/api/metrics
```

**Restart dashboard:**
```bash
# Stop with Ctrl+C
make run-webui
```

### Load Generator Not Working

**Check if loadgen is running:**
```bash
ps aux | grep loadgenerator
```

**Start it:**
```bash
make run-loadgen
```

## 📊 Understanding the Metrics

### Request Rate
- **Good**: 50-200 req/sec (with load generator)
- **Low**: < 10 req/sec (no traffic)
- **High**: > 500 req/sec (heavy load)

### Error Rate
- **Good**: < 5% (0.05)
- **Warning**: 5-10%
- **Critical**: > 10%

### Latency
- **Good**: < 200ms average
- **OK**: 200-500ms
- **Slow**: > 500ms

### P95 Latency
- **Good**: < 500ms
- **OK**: 500ms-1s
- **Slow**: > 1s

## 🎓 Tips & Best Practices

1. **Keep dashboard open** during development
2. **Watch for red cards** (unhealthy services)
3. **Monitor error rate** closely
4. **Use trace IDs** to correlate logs with traces
5. **Start load generator** to see realistic data
6. **Check metrics** before and after changes
7. **Use tool cards** for deep dives
8. **Refresh manually** when troubleshooting

## 🔗 Related Documentation

- **README.md** - Main documentation
- **DEMO_ARCHITECTURE.md** - System architecture
- **QUICK_REFERENCE.md** - Command reference
- **EXAMPLES.md** - Code examples

## 🎉 Quick Demo

```bash
# 1. Start all services
make run-all-local

# 2. Open dashboard
open http://localhost:3001

# 3. Start load generator (from dashboard)
Click "Start Load Generator"

# 4. Watch metrics update
Observe request rate, latency, logs

# 5. View traces
Click "Jaeger" → Select "frontend" → "Find Traces"

# 6. See service dependencies
Explore trace details and timing

# 7. Check logs
Scroll to "Recent Logs" section

# 8. Stop load generator
Click "Stop Load Generator"
```

---

**The Web UI Dashboard is your command center for observability!** 🎮

Open http://localhost:3001 and start monitoring your microservices! 🚀

