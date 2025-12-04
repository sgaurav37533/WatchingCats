# OpenTelemetry Collector Dashboard - Quick Start

## 🎯 Goal

View real-time metrics about how telemetry data flows through your OpenTelemetry Collector.

## ⚡ Quick Start (3 Steps)

### 1. Start the Stack

```bash
make docker-up
```

Wait 30 seconds for all services to initialize.

### 2. Open the Dashboard

1. Go to **http://localhost:3000** (Grafana)
2. Login: `admin` / `admin`
3. Click menu (☰) → **Dashboards**
4. Open folder: **OpenTelemetry**
5. Click: **OpenTelemetry Collector Data Flow**

### 3. Generate Traffic

```bash
make run-loadgen
```

Watch the dashboard come alive with real-time data!

## 📊 What You'll See

### Process Metrics
- **Memory usage** of the collector
- **CPU utilization**

### Traces Flow
- Traces **received** from your services
- Traces **sent** to Jaeger
- **Export ratio** (should be green/100%)

### Metrics Flow  
- Metrics **received** from your services
- Metrics **sent** to Prometheus
- **Export ratio** (health indicator)

### Data Processing
- **Batch processor** statistics
- **Prometheus scraping** metrics

## 🎨 Dashboard Layout

```
┌─────────────────────────────────────────────┐
│   Process Metrics                           │
│   ├─ Memory RSS          ├─ CPU Usage       │
├─────────────────────────────────────────────┤
│   Traces Pipeline                           │
│   ├─ Ingress (Receiver)  ├─ Egress (Exporter)│
│   ├─ Export Ratio        ├─ Batch Processor │
├─────────────────────────────────────────────┤
│   Metrics Pipeline                          │
│   ├─ Ingress (Receiver)  ├─ Egress (Exporter)│
│   ├─ Export Ratio        ├─ Batch Processor │
├─────────────────────────────────────────────┤
│   Prometheus Scraping                       │
│   ├─ Samples Scraped     ├─ Scrape Duration │
└─────────────────────────────────────────────┘
```

## 🔍 What to Look For

### Healthy System
✅ **Export Ratios**: Green (>95%)  
✅ **Refused Metrics**: 0  
✅ **Failed Exports**: 0  
✅ **Memory**: Steady or slight growth  
✅ **CPU**: <50% normally

### Problems Detected
🔴 **Red Export Ratio**: Data loss occurring  
🔴 **High Refused Count**: Memory limits hit  
🔴 **Failed Exports**: Backend connectivity issues  
🔴 **Rising Memory**: Potential memory leak  
🔴 **High CPU**: Processing bottleneck

## 🛠️ Troubleshooting

### No Data Showing?

```bash
# 1. Verify collector is running
docker-compose ps otel-collector

# 2. Check metrics endpoint
curl http://localhost:8888/metrics | grep otelcol

# 3. Verify Prometheus target
open http://localhost:9090/targets
# Look for: otel-collector target is UP
```

### Dashboard Not Visible?

```bash
# Restart Grafana
docker-compose restart grafana

# Check logs
docker-compose logs grafana
```

### Still Having Issues?

```bash
# Run verification script
make verify-dashboard
```

## 📚 Learn More

- **Detailed Guide**: [COLLECTOR_DASHBOARD_GUIDE.md](COLLECTOR_DASHBOARD_GUIDE.md)
- **Implementation Details**: [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)
- **Official Docs**: https://opentelemetry.io/docs/demo/collector-data-flow-dashboard/

## 💡 Pro Tips

1. **Refresh Rate**: Dashboard auto-refreshes every 10 seconds
2. **Time Range**: Default is last 30 minutes (adjust in top-right)
3. **Zoom In**: Click and drag on any graph to zoom
4. **Panel Details**: Click any panel title → View for full-screen
5. **Export Data**: Panel menu → Inspect → Data

## 🎓 Understanding Export Ratios

The export ratio shows data flow efficiency:

```
Export Ratio = Sent / Received

100% (Green)  = Perfect - All data exported
90% (Yellow)  = Warning - Some data dropped
<80% (Red)    = Critical - Significant data loss
```

**Note**: Metrics ratio can be >100% if processors generate additional metrics!

## 🚀 Next Steps

Once you understand the collector's health:

1. **Create Alerts**: Set up Grafana alerts for export ratios
2. **Optimize Config**: Adjust batch sizes and memory limits
3. **Custom Panels**: Add metrics specific to your use case
4. **Monitor Trends**: Track metrics over time to identify patterns

## 🔗 Related Dashboards

You can create additional dashboards for:
- Application-specific metrics
- Service health and SLOs
- Business metrics
- Custom KPIs

Place them in: `configs/dashboards/your-dashboard.json`

## ✅ Verification Checklist

Before considering setup complete:

- [ ] Can access Grafana at http://localhost:3000
- [ ] Dashboard appears under OpenTelemetry folder
- [ ] Process metrics show CPU and memory
- [ ] With traffic, see spans/metrics flowing
- [ ] Export ratios display (green is good)
- [ ] No errors in Grafana logs

## 🎉 Success!

You now have complete visibility into your OpenTelemetry Collector's data pipeline!

This dashboard helps ensure:
- ✅ Telemetry reaches the collector
- ✅ Data is processed correctly
- ✅ Exports succeed to backends
- ✅ Collector remains healthy

**Happy Monitoring!** 📊

---

**Quick Commands**

```bash
make docker-up           # Start everything
make run-loadgen        # Generate traffic
make verify-dashboard   # Verify setup
make docker-down        # Stop everything
```

