# Grafana Dashboards

This directory contains Grafana dashboard definitions that are automatically provisioned when Grafana starts.

## Available Dashboards

### OpenTelemetry Collector Data Flow
**File**: `otel-collector-dataflow.json`  
**Location in Grafana**: Dashboards → OpenTelemetry → OpenTelemetry Collector Data Flow

**Purpose**: Monitors the health and data flow of the OpenTelemetry Collector

**Panels**:
- Process Metrics (CPU, Memory)
- Traces Pipeline (Ingress/Egress)
- Metrics Pipeline (Ingress/Egress)
- Export Ratios (Health indicators)
- Batch Processor Statistics
- Prometheus Scraping Metrics

**Key Metrics**:
- `otelcol_receiver_accepted_spans` - Spans received
- `otelcol_receiver_refused_spans` - Spans rejected
- `otelcol_exporter_sent_spans` - Spans exported
- `otelcol_exporter_send_failed_spans` - Export failures
- `otelcol_process_memory_rss` - Memory usage
- `otelcol_process_cpu_seconds` - CPU usage

## Auto-Provisioning

Dashboards in this directory are automatically loaded by Grafana on startup via the provisioning configuration defined in `../grafana-dashboards.yaml`.

## Adding New Dashboards

To add a new dashboard:

1. Create or export your dashboard JSON
2. Place it in this directory: `configs/dashboards/your-dashboard.json`
3. Restart Grafana or wait for auto-reload (10 seconds)
4. Dashboard will appear in Grafana under the "OpenTelemetry" folder

## Customization

You can customize existing dashboards by:

1. Opening them in Grafana
2. Making your changes
3. Exporting the JSON
4. Replacing the file in this directory
5. Restarting Grafana to load changes

**Tip**: Set `allowUiUpdates: true` in the provisioning config to allow live editing in Grafana UI.

## Documentation

For detailed information about the Collector Dashboard:
- See `../../COLLECTOR_DASHBOARD_GUIDE.md`
- Official docs: https://opentelemetry.io/docs/demo/collector-data-flow-dashboard/

## Verification

To verify the dashboard is properly configured:

```bash
make verify-dashboard
```

This will check:
- ✅ Dashboard files exist
- ✅ Provisioning is configured
- ✅ Required metrics are defined
- ✅ Docker Compose mounts are correct
- ✅ Services are healthy (if running)

