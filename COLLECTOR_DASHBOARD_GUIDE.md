# OpenTelemetry Collector Data Flow Dashboard Guide

This guide explains the OpenTelemetry Collector Data Flow Dashboard implementation in the WatchingCat project, based on the [official OpenTelemetry documentation](https://opentelemetry.io/docs/demo/collector-data-flow-dashboard/).

## Overview

The Collector Data Flow Dashboard provides visibility into how telemetry data flows through the OpenTelemetry Collector, including:
- Ingress and egress metrics for traces and metrics
- Process-level metrics (CPU, memory)
- Export ratios to understand data flow efficiency
- Batch processor statistics
- Prometheus scraping metrics

## Architecture

The dashboard monitors the OpenTelemetry Collector's internal telemetry, which is:
1. **Exposed** by the collector on port 8888
2. **Scraped** by Prometheus every 15 seconds
3. **Visualized** in Grafana through the dashboard

## Dashboard Sections

### 1. Process Metrics
Monitors the health of the collector process itself:
- **Memory RSS**: Resident Set Size showing actual memory usage
- **CPU Usage**: Percentage of CPU utilized by the collector

### 2. Traces Pipeline
Monitors the flow of trace data through the collector:

#### Ingress (Receiver)
- `otelcol_receiver_accepted_spans`: Traces successfully received
- `otelcol_receiver_refused_spans`: Traces rejected (e.g., memory limits exceeded)
- Labels: `receiver`, `transport`

#### Egress (Exporter)
- `otelcol_exporter_sent_spans`: Traces successfully exported
- `otelcol_exporter_send_failed_spans`: Traces that failed to export
- Labels: `exporter`

#### Export Ratio
Calculated as: `sent_spans / accepted_spans`
- **Green (>0.95)**: Healthy, most traces are being exported
- **Yellow (0.8-0.95)**: Some traces may be dropped
- **Red (<0.8)**: Significant trace loss

### 3. Metrics Pipeline
Monitors the flow of metric data through the collector:

#### Ingress (Receiver)
- `otelcol_receiver_accepted_metric_points`: Metrics successfully received
- `otelcol_receiver_refused_metric_points`: Metrics rejected
- Labels: `receiver`, `transport`

#### Egress (Exporter)
- `otelcol_exporter_sent_metric_points`: Metrics successfully exported
- `otelcol_exporter_send_failed_metric_points`: Metrics that failed to export
- Labels: `exporter`

#### Export Ratio
Calculated as: `sent_metric_points / accepted_metric_points`

**Note**: The export ratio for metrics may be >1 if processors generate additional metrics (e.g., span metrics processor).

### 4. Prometheus Scraping
Shows Prometheus scraping activity:
- **Scrape Samples**: Number of metrics scraped from the collector
- **Scrape Duration**: Time taken to scrape metrics

## Key Metrics to Monitor

### Health Indicators
1. **Process Memory**: Steadily increasing memory may indicate a memory leak
2. **Process CPU**: High sustained CPU usage may indicate processing bottlenecks
3. **Export Ratios**: Low ratios indicate data loss

### Data Flow Issues
1. **Refused Spans/Metrics**: Non-zero values indicate:
   - Memory limit exceeded (configured in `memory_limiter` processor)
   - Queue overflow
   - Rate limiting

2. **Failed Exports**: Non-zero values indicate:
   - Backend unavailability
   - Network issues
   - Authentication failures

### Batch Processor
The batch processor aggregates telemetry before export:
- Monitor batch sizes to understand throughput
- Large spikes may indicate bursty traffic

## Configuration Files

### Collector Configuration
File: `configs/otel-collector-config.yaml`

Key settings:
```yaml
service:
  telemetry:
    metrics:
      address: 0.0.0.0:8888  # Internal metrics endpoint
```

### Prometheus Configuration
File: `configs/prometheus.yml`

The collector is scraped using:
```yaml
scrape_configs:
  - job_name: 'otel-collector'
    static_configs:
      - targets: ['otel-collector:8888']
    metric_relabel_configs:
      - source_labels: [__name__]
        regex: 'otelcol_.*'
        action: keep
```

### Grafana Dashboard Provisioning
Files:
- `configs/grafana-dashboards.yaml`: Dashboard provider configuration
- `configs/dashboards/otel-collector-dataflow.json`: Dashboard definition

## Accessing the Dashboard

1. **Start the stack**:
   ```bash
   docker-compose up -d
   ```

2. **Access Grafana**:
   - URL: http://localhost:3000
   - Username: `admin`
   - Password: `admin`

3. **Find the dashboard**:
   - Click the menu icon (☰) on the left
   - Navigate to **Dashboards**
   - Look for the **OpenTelemetry** folder
   - Select **OpenTelemetry Collector Data Flow**

## Troubleshooting

### Dashboard Shows No Data

1. **Check if collector is running**:
   ```bash
   docker-compose ps otel-collector
   ```

2. **Verify metrics endpoint**:
   ```bash
   curl http://localhost:8888/metrics
   ```
   You should see metrics starting with `otelcol_`

3. **Check Prometheus is scraping**:
   - Go to http://localhost:9090
   - Navigate to Status → Targets
   - Verify `otel-collector` target is UP

4. **Verify Prometheus has data**:
   - Go to http://localhost:9090
   - Query: `otelcol_receiver_accepted_spans`
   - Should show results if traces are flowing

### High Memory Usage

If `otelcol_process_memory_rss` keeps increasing:
1. Check the memory_limiter configuration
2. Review batch processor settings
3. Check for exporter failures causing backlog

### Low Export Ratio

If export ratios are consistently low:
1. Check `otelcol_receiver_refused_*` metrics
2. Increase memory limits if needed
3. Verify exporters are healthy
4. Check network connectivity to backends

## Customization

You can customize the dashboard by:

1. **Adding panels**: Add memory_limiter processor metrics
2. **Adjusting thresholds**: Modify gauge thresholds based on your SLOs
3. **Adding alerts**: Create Grafana alerts based on these metrics
4. **Filtering**: Add template variables to filter by service or exporter

### Example: Adding Memory Limiter Metrics

Add a panel with:
```promql
otelcol_processor_refused_spans{processor="memory_limiter"}
```

## Reference

- [Official OpenTelemetry Collector Data Flow Dashboard Documentation](https://opentelemetry.io/docs/demo/collector-data-flow-dashboard/)
- [OpenTelemetry Collector Metrics](https://opentelemetry.io/docs/collector/internal-telemetry/)
- [Collector Configuration Reference](https://opentelemetry.io/docs/collector/configuration/)

## Integration with WatchingCat

This dashboard complements the WatchingCat observability stack:

- **Traces**: View collector trace flow → Full traces in Jaeger
- **Metrics**: Monitor collector metrics → Application metrics in Grafana
- **Logs**: Collector logging exporter → Logs in Elasticsearch/Kibana

The dashboard helps you understand the "plumbing" of your observability pipeline, ensuring telemetry data flows correctly from your applications to your backends.

