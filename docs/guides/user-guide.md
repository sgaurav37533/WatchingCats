# WatchingCat Modern UI Guide

## 🎨 Overview

The WatchingCat Modern UI is a comprehensive observability platform interface inspired by the [OpenTelemetry Astronomy Shop Demo](https://github.com/open-telemetry/opentelemetry-demo). It provides a production-ready, feature-rich interface for monitoring distributed systems.

## ✨ Features

### 1. **Dashboard** - Complete System Overview
- **System Status Banner**: Real-time health indicators for all services
- **Key Metrics Cards**: Request rate, success rate, latency, and error rate
- **Real-time Charts**: Request volume and response time percentiles
- **Service Topology**: Interactive visualization of microservice dependencies
- **Observability Tools**: Quick access to Jaeger, Grafana, Prometheus, and Kibana

### 2. **Services** - Microservices Monitoring
- Detailed health status for each service
- Real-time performance metrics (requests/sec, latency, error rate)
- Service endpoint information
- Visual health indicators

### 3. **Traces** - Distributed Tracing
- Trace search and filtering
- Service-specific trace filtering
- Trace timeline visualization
- Direct links to detailed trace analysis in Jaeger

### 4. **Metrics** - System Performance
- CPU usage by service (bar chart)
- Memory usage by service (bar chart)
- Network throughput over time (line chart)
- Real-time metric updates

### 5. **Demo Shop** - E-commerce Simulation
- Interactive product catalog (Observatory Shop theme)
- Shopping cart functionality
- Checkout simulation to generate traces
- Load generator controls
- Realistic telemetry data generation

## 🚀 Quick Start

### Access the UI

```bash
# Option 1: Run locally
make run-webui

# Option 2: Run in Docker
make docker-up

# Open your browser
open http://localhost:3001
```

### Navigation

The UI has 5 main sections accessible via the top navigation bar:

1. **Dashboard** (🏠) - System overview and key metrics
2. **Services** (🖥️) - Microservices health and performance
3. **Traces** (🔀) - Distributed trace exploration
4. **Metrics** (📊) - Detailed system metrics
5. **Demo Shop** (🛒) - E-commerce simulation for testing

## 📊 Dashboard Features

### Status Banner

Shows real-time system health:
- **All Systems Operational** (🟢) - All services healthy
- **Degraded Performance** (🟡) - Some services unhealthy
- **System Down** (🔴) - All services down

### Key Metrics

Four primary metrics with trend indicators:

**Request Rate**
- Current requests per second
- Trend indicator (↑ or ↓)

**Success Rate**
- Percentage of successful requests
- Calculated as 100% minus error rate

**Average Latency**
- Mean response time across all services
- Measured in milliseconds

**Error Rate**
- Percentage of failed requests
- Critical threshold: > 5%

### Real-time Charts

**Request Volume Chart**
- Shows request rate over time
- Adjustable time range (5m, 15m, 1h, 24h)
- Area chart with gradient fill

**Response Times Chart**
- Three percentile lines:
  - **P50** (🟢) - Median latency
  - **P95** (🟡) - 95th percentile
  - **P99** (🔴) - 99th percentile

### Service Topology

Interactive visualization showing:
- **Services**: Frontend, Cart, Product Catalog, Checkout
- **Backend Systems**: OpenTelemetry Collector, Jaeger, Prometheus, ELK
- **Data Flow**: Animated dots showing telemetry flow
- **Health Status**: Green indicators on healthy services

**Interaction:**
- Hover over nodes for details
- Click "Reset View" to reposition
- Status indicators update in real-time

### Observability Tools

Quick access cards for:

**Jaeger** (🔀)
- Distributed tracing
- Request flow analysis
- Shows trace count

**Grafana** (📊)
- Metrics dashboards
- Custom visualizations
- Shows dashboard count

**Prometheus** (🔥)
- Metrics storage
- PromQL queries
- Shows metric count

**Kibana** (🔍)
- Log search
- Analysis and visualization
- Shows log count

## 🖥️ Services Page

### Service Cards

Each microservice displays:

**Header**
- Service name
- Endpoint URL
- Health badge (Healthy/Unhealthy)

**Statistics**
- **Request Rate**: Requests per second
- **Latency**: Average response time
- **Error Rate**: Percentage of failed requests

**Features**
- Auto-refresh every 5 seconds
- Click to view detailed metrics
- Color-coded health indicators

## 🔀 Traces Page

### Trace Search

**Search Bar**
- Search by trace ID
- Search by service name
- Real-time filtering

**Service Filter**
- Filter by specific service
- Options: Frontend, Cart, Catalog, Checkout
- "All Services" to view everything

**Trace List**

Each trace displays:
- **Trace ID**: 32-character hex identifier
- **Duration**: Total time in milliseconds
- **Timestamp**: When trace was created
- **Services**: Which services were involved

**Interaction**
- Click any trace to open in Jaeger
- Sorted by most recent first
- Auto-refresh available

## 📊 Metrics Page

### CPU Usage Chart
- Bar chart showing CPU percentage per service
- Color-coded by service
- Scale: 0-100%

### Memory Usage Chart
- Bar chart showing memory in MB per service
- Color-coded by service
- Real-time updates

### Network Throughput Chart
- Line chart showing MB/s over time
- Area chart with gradient
- 20-minute time window

## 🛒 Demo Shop

### Product Catalog

**Observatory-themed Products:**
1. Observatory Telescope Pro - $2,999.99 🔭
2. Star Chart Collection - $149.99 🗺️
3. Night Vision Binoculars - $599.99 🔍
4. Astronomy Guide Book - $49.99 📚
5. Portable Planetarium - $899.99 🌍
6. Space Photography Kit - $1,299.99 📷

### Shopping Cart

**Features:**
- Add products to cart
- View cart total
- Remove items
- Checkout process

**Telemetry Generation:**
When you checkout:
- Generates distributed traces
- Creates metrics in Prometheus
- Logs events in ELK stack
- Viewable in Jaeger immediately

### Load Generator

**Purpose:** Simulate automated shopping behavior

**Controls:**
- **Start**: Begin generating requests
- **Stop**: Halt request generation

**Rate:** 30 requests per minute by default

**What it does:**
- Browses products randomly
- Adds items to cart
- Completes checkouts
- Generates realistic telemetry

## 🎨 Theme Support

### Light / Dark Mode

**Toggle Theme:**
- Click the moon icon (🌙) in the top-right
- Switches between light and dark themes
- Preference saved locally

**Color Schemes:**

**Light Theme:**
- White backgrounds
- Dark text
- High contrast

**Dark Theme:**
- Dark blue backgrounds
- Light text
- Reduced eye strain

## 🔄 Auto-Refresh

**Dashboard Page:**
- Refreshes every 5 seconds
- Updates metrics automatically
- Real-time status indicators

**Manual Refresh:**
- Click refresh button (🔄) anytime
- Updates current page immediately
- Shows toast notification

## 📱 Responsive Design

**Desktop (> 1200px):**
- Full 4-column metric cards
- Side-by-side charts
- Complete navigation bar

**Tablet (768px - 1200px):**
- 2-column metric cards
- Stacked charts
- Collapsible navigation

**Mobile (< 768px):**
- Single column layout
- Stacked components
- Touch-optimized controls

## 🔔 Toast Notifications

**Types:**
- **Info** (🔵) - General information
- **Success** (🟢) - Action completed
- **Warning** (🟡) - Caution advised
- **Error** (🔴) - Action failed

**Behavior:**
- Appears top-right
- Slides in with animation
- Auto-dismiss after 3 seconds
- Stack multiple notifications

## 🎯 Use Cases

### 1. Monitor Service Health

**Scenario:** Check if all services are running

**Steps:**
1. Go to Dashboard
2. Check status banner (should be green)
3. Review service topology
4. Visit Services page for details

### 2. Investigate High Latency

**Scenario:** Response times are slow

**Steps:**
1. Check Dashboard latency metric
2. View response times chart (P95, P99)
3. Go to Metrics page
4. Review CPU and memory usage
5. Check Traces for slow requests

### 3. Debug Failed Requests

**Scenario:** Errors appearing in logs

**Steps:**
1. Check error rate on Dashboard
2. Go to Traces page
3. Filter by service
4. Click trace to view in Jaeger
5. Analyze span errors

### 4. Generate Test Traffic

**Scenario:** Need telemetry data for testing

**Steps:**
1. Go to Demo Shop
2. Click "Start" load generator
3. Or manually shop and checkout
4. View traces in Jaeger
5. Check metrics in Grafana

### 5. Demonstrate Observability

**Scenario:** Show observability to stakeholders

**Steps:**
1. Start load generator
2. Open Dashboard
3. Show real-time metrics
4. Display service topology
5. Click into Jaeger/Grafana
6. Demonstrate trace correlation

## 🔧 Customization

### Modify Time Ranges

Edit chart time ranges in `modern-app.js`:

```javascript
function updateCharts(timeRange) {
    const count = timeRange === '5m' ? 5 : 
                  timeRange === '15m' ? 15 : 
                  timeRange === '1h' ? 60 : 24;
    // Update charts...
}
```

### Add Products

Add new products to the shop in `modern-app.js`:

```javascript
const products = [
    { id: 7, name: 'Your Product', price: 99.99, image: '🎁' },
    // ...
];
```

### Customize Colors

Edit color scheme in `modern-ui.css`:

```css
:root {
    --primary: #6366f1;  /* Change primary color */
    --success: #10b981;  /* Change success color */
    /* ... */
}
```

### Modify Refresh Rate

Change auto-refresh interval in `modern-app.js`:

```javascript
function startAutoRefresh() {
    state.refreshInterval = setInterval(() => {
        // ...
    }, 5000);  // Change this value (milliseconds)
}
```

## 📊 API Endpoints

The UI consumes these APIs:

### GET /api/services
Returns service health status

**Response:**
```json
[
    {
        "name": "Frontend",
        "url": "http://localhost:8080",
        "status": "healthy",
        "healthy": true,
        "timestamp": "2025-12-04T10:30:00Z"
    }
]
```

### GET /api/metrics
Returns current system metrics

**Response:**
```json
{
    "request_rate": 150.5,
    "error_rate": 0.05,
    "avg_latency_ms": 245.3,
    "p95_latency_ms": 450.2,
    "total_requests": 125000,
    "total_errors": 6250
}
```

### GET /api/logs
Returns recent log entries

**Response:**
```json
[
    {
        "timestamp": "2025-12-04T10:25:00Z",
        "level": "error",
        "service": "checkoutservice",
        "message": "Payment processing failed",
        "trace_id": "xyz789uvw012"
    }
]
```

### POST /api/loadgen/start
Starts the load generator

**Response:**
```json
{
    "status": "started",
    "message": "Load generator started successfully",
    "rate": "30 requests/min"
}
```

### POST /api/loadgen/stop
Stops the load generator

**Response:**
```json
{
    "status": "stopped",
    "message": "Load generator stopped successfully"
}
```

## 🐛 Troubleshooting

### UI Not Loading

**Problem:** Blank page or 404 errors

**Solution:**
1. Check if webui service is running:
   ```bash
   curl http://localhost:3001/health
   ```
2. Verify template files exist:
   ```bash
   ls web/templates/
   ```
3. Check server logs for errors

### Charts Not Rendering

**Problem:** Chart areas are blank

**Solution:**
1. Check browser console for errors
2. Verify Chart.js is loaded
3. Ensure canvas elements have proper IDs
4. Try hard refresh (Cmd+Shift+R / Ctrl+Shift+R)

### Data Not Updating

**Problem:** Metrics show "--" or old data

**Solution:**
1. Check API endpoints:
   ```bash
   curl http://localhost:3001/api/metrics
   curl http://localhost:3001/api/services
   ```
2. Verify backend services are running
3. Check network tab in browser DevTools

### Topology Not Showing

**Problem:** Service topology is empty

**Solution:**
1. Check if D3.js is loaded
2. Verify container has proper dimensions
3. Check browser console for errors
4. Click "Reset View" button

## 🚀 Performance Tips

### Optimize for Many Services

If you have > 10 services:
1. Increase grid columns in CSS
2. Add pagination to services page
3. Implement virtual scrolling

### Reduce Data Transfer

For better performance:
1. Implement data pagination
2. Add caching for metrics
3. Use WebSockets for real-time updates
4. Compress API responses

## 📚 Related Documentation

- [OpenTelemetry Demo](https://github.com/open-telemetry/opentelemetry-demo)
- [Collector Dashboard Guide](COLLECTOR_DASHBOARD_GUIDE.md)
- [Quick Reference](QUICK_REFERENCE.md)
- [Architecture](ARCHITECTURE.md)

## 🎉 What's Next?

**Enhancements to Consider:**
1. Add user authentication
2. Implement alerting UI
3. Create custom dashboard builder
4. Add service dependency graph
5. Implement log streaming
6. Add incident management
7. Create SLO/SLI tracking

## 💡 Tips & Tricks

1. **Keyboard Shortcuts**: Press `r` to refresh, `d` for dashboard
2. **Quick Toggle**: Double-click theme button for instant switch
3. **Cart Shortcuts**: Press `c` to view cart, `Enter` to checkout
4. **Trace Search**: Use `Cmd/Ctrl + F` in traces page
5. **Screenshot**: Use browser tools to capture for reports

---

**Built with ❤️ inspired by the OpenTelemetry Community**

**Version**: 1.0.0  
**Last Updated**: December 4, 2025  
**Status**: ✅ Production Ready

