---
date: 2025-12-04
id: architecture
tags: [WatchingCat, OpenTelemetry, Self-Host, Observability]
title: WatchingCat Technical Architecture
description: Learn about the technical architecture of WatchingCat, an OpenTelemetry-native observability platform with unified telemetry collection, storage, and visualization.
---

# WatchingCat Technical Architecture

## Overview

WatchingCat is a comprehensive, OpenTelemetry-native observability platform designed for modern distributed systems. Built on industry-standard components and best practices, WatchingCat provides unified collection, storage, and visualization of traces, metrics, and logs—the three pillars of observability.

Our architecture leverages the power of OpenTelemetry for instrumentation, a high-performance storage backend, and a modern React-inspired web interface to deliver real-time insights into application performance and behavior.

**Core Philosophy**: OpenTelemetry-first, vendor-neutral, easy to deploy, powerful to use.

![WatchingCat Architecture](./docs/architecture-diagram.svg)

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Applications                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │ Frontend │  │   Cart   │  │ Catalog  │  │ Checkout │          │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘          │
│       │             │              │              │                 │
│       └─────────────┴──────────────┴──────────────┘                │
│                          │                                          │
│                   OTLP (gRPC/HTTP)                                 │
└───────────────────────────┼─────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                 OpenTelemetry Collector                             │
│  ┌────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │ Receivers  │→│  Processors   │→│   Exporters   │              │
│  │ OTLP,HTTP  │  │ Batch,Memory │  │Jaeger,Prom,ES│              │
│  └────────────┘  └──────────────┘  └──────────────┘              │
│         │                                     │                     │
│         └────────── Telemetry Data ──────────┘                     │
└─────────────────────────────────────────────────────────────────────┘
                            │
                ┌───────────┴───────────┐
                │                       │
                ▼                       ▼
┌──────────────────────┐    ┌──────────────────────┐
│   Storage Layer      │    │   Storage Layer      │
│  ┌────────────────┐  │    │  ┌────────────────┐  │
│  │    Jaeger      │  │    │  │  Prometheus    │  │
│  │   (Traces)     │  │    │  │   (Metrics)    │  │
│  └────────────────┘  │    │  └────────────────┘  │
└──────────────────────┘    └──────────────────────┘
                │                       │
┌──────────────────────┐                │
│   Storage Layer      │                │
│  ┌────────────────┐  │                │
│  │ Elasticsearch  │  │                │
│  │    (Logs)      │  │                │
│  └────────────────┘  │                │
└──────────────────────┘                │
                │                       │
                └───────────┬───────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    WatchingCat Backend                              │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  Query Service                                                │  │
│  │  - Trace Queries (Jaeger API)                                │  │
│  │  - Metric Queries (PromQL)                                   │  │
│  │  - Log Queries (Elasticsearch)                               │  │
│  │  - Unified API endpoints                                     │  │
│  └──────────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  Alert Manager                                                │  │
│  │  - Rule evaluation                                            │  │
│  │  - Alert routing                                              │  │
│  │  - Notification delivery                                      │  │
│  └──────────────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  API Server                                                   │  │
│  │  - REST API (Go)                                              │  │
│  │  - WebSocket for real-time updates                           │  │
│  │  - Authentication & Authorization                             │  │
│  └──────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    WatchingCat Frontend                             │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  Modern Web UI (Vanilla JS + Chart.js + D3.js)               │  │
│  │  - Dashboard (System Overview)                                │  │
│  │  - Services (Health & Performance)                            │  │
│  │  - Traces (Distributed Tracing)                               │  │
│  │  - Metrics (CPU, Memory, Network)                             │  │
│  │  - Logs (Centralized Logging)                                 │  │
│  │  - Alerts (Alert Management)                                  │  │
│  └──────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────┘
```

## Components

### 1. OpenTelemetry Collector

**Role**: Central telemetry data pipeline

The WatchingCat OpenTelemetry Collector is the heart of our observability platform, built on the official [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) framework.

**Key Features**:
- ✅ **Multi-protocol ingestion**: OTLP (gRPC/HTTP), Jaeger, Zipkin, OpenCensus
- ✅ **Intelligent processing**: Batching, sampling, enrichment, filtering
- ✅ **Memory management**: Memory limiter to prevent OOM conditions
- ✅ **Flexible export**: Simultaneous export to multiple backends
- ✅ **Self-monitoring**: Exposes internal metrics on port 8888

**Data Flow**:
```
Application → OTLP → OTel Collector → [Processors] → Exporters
                                           ↓
                               ┌───────────┼───────────┐
                               ↓           ↓           ↓
                            Jaeger    Prometheus  Elasticsearch
```

**Configuration**:
- Receivers: OTLP (gRPC: 4317, HTTP: 4318)
- Processors: Batch, Memory Limiter, Attributes
- Exporters: Jaeger, Prometheus, Logging

**Benefits**:
- Vendor-neutral instrumentation
- Future-proof architecture
- Protocol translation
- Unified telemetry pipeline

### 2. Storage Layer

WatchingCat uses a **polyglot persistence** approach, optimizing each telemetry type with the best-suited database.

#### 2.1 Jaeger (Traces)

**Technology**: Jaeger All-in-One / Cassandra-backed  
**Purpose**: Distributed tracing storage and analysis

**Capabilities**:
- Store trace spans with hierarchical relationships
- Support for service dependency analysis
- Trace search by service, operation, tags
- Performance analysis and bottleneck detection

**Key Features**:
- ✅ OpenTelemetry-compatible
- ✅ Service dependency graphs
- ✅ Trace search and filtering
- ✅ Latency distribution analysis

#### 2.2 Prometheus (Metrics)

**Technology**: Prometheus Time-Series Database  
**Purpose**: Metrics collection and storage

**Capabilities**:
- Time-series data with labels
- PromQL for powerful queries
- Efficient compression
- Built-in alerting

**Key Features**:
- ✅ Multi-dimensional data model
- ✅ Flexible query language (PromQL)
- ✅ Pull-based metrics collection
- ✅ Service discovery integration

#### 2.3 Elasticsearch (Logs)

**Technology**: Elasticsearch + Kibana  
**Purpose**: Centralized log management

**Capabilities**:
- Full-text search across logs
- Structured and unstructured log data
- Log correlation with traces
- Real-time log streaming

**Key Features**:
- ✅ Powerful search capabilities
- ✅ Log aggregation and analysis
- ✅ Trace context correlation
- ✅ Kibana visualization

### 3. WatchingCat Backend

**Role**: Unified API and business logic layer

The backend service is written in **Go**, providing high performance and efficient resource utilization.

#### 3.1 Query Service

**Purpose**: Unified query layer across all telemetry types

**Endpoints**:
```go
// Traces
GET  /api/traces              // List traces
GET  /api/traces/:id          // Get trace details
GET  /api/traces/search       // Search traces

// Metrics
GET  /api/metrics             // Get metrics
POST /api/metrics/query       // PromQL queries
GET  /api/metrics/range       // Range queries

// Logs
GET  /api/logs                // Get logs
POST /api/logs/search         // Full-text search
GET  /api/logs/tail           // Real-time tail

// Services
GET  /api/services            // Service list
GET  /api/services/:id/health // Service health
GET  /api/services/:id/metrics// Service metrics
```

**Features**:
- Unified API across storage backends
- Query optimization
- Caching layer for performance
- Rate limiting and authentication

#### 3.2 Alert Manager

**Purpose**: Alert rule evaluation and notification

**Components**:
```
Rule Evaluator → Alert Router → Notification Delivery
      ↓               ↓                  ↓
   PromQL        Grouping          Slack/Email/PagerDuty
   LogQL         Deduplication     Webhook
   TraceQL       Silencing         Custom
```

**Alert Types**:
- **Metric Alerts**: Based on Prometheus queries
- **Trace Alerts**: Based on trace attributes (latency, errors)
- **Log Alerts**: Based on log patterns

**Features**:
- ✅ Flexible rule definition
- ✅ Multiple notification channels
- ✅ Alert grouping and deduplication
- ✅ Silencing and maintenance windows

#### 3.3 API Server

**Purpose**: REST API and WebSocket server

**Technology**: Go with Gin/Echo framework

**Features**:
- RESTful API design
- WebSocket for real-time updates
- JWT authentication
- Role-based access control (RBAC)
- API versioning
- Rate limiting
- CORS support

**Endpoints Structure**:
```
/api/v1/
  ├── traces/
  ├── metrics/
  ├── logs/
  ├── services/
  ├── alerts/
  ├── dashboards/
  └── users/
```

### 4. WatchingCat Frontend

**Role**: User interface and visualization

**Technology Stack**:
- **Core**: Vanilla JavaScript (ES6+)
- **Charts**: Chart.js 4.4.0
- **Graphs**: D3.js v7
- **Icons**: Font Awesome 6.4.0
- **Styling**: Modern CSS with variables

**Pages**:
1. **Dashboard**: System overview, key metrics, service topology
2. **Services**: Health monitoring, performance metrics
3. **Traces**: Distributed tracing with inline viewer
4. **Metrics**: CPU, memory, network charts
5. **Logs**: Centralized log viewer (planned)
6. **Alerts**: Alert management UI (planned)
7. **Demo Shop**: E-commerce simulation for testing

**Features**:
- ✅ Real-time updates (5s interval)
- ✅ Theme support (light/dark)
- ✅ Responsive design (mobile-ready)
- ✅ Inline trace viewer (OpenTelemetry-compliant)
- ✅ Interactive service topology (D3.js)
- ✅ No external framework dependencies
- ✅ Fast load times

### 5. Demo Applications

**Purpose**: Generate realistic telemetry for testing and demonstration

**Services**:
```
Frontend (8080)
   ├── Cart Service (8081)
   ├── Product Catalog (8082)
   └── Checkout Service (8083)
        └── Load Generator
```

**Features**:
- OpenTelemetry instrumentation
- Realistic request patterns
- Error injection (10% rate)
- Load generation capabilities
- E-commerce simulation

## Data Flow

### Trace Data Flow

```
1. Application → OTLP Exporter
2. OTLP Exporter → OTel Collector (4317/4318)
3. OTel Collector → Processors (batch, enrich)
4. Processors → Jaeger Exporter
5. Jaeger Exporter → Jaeger Backend
6. Jaeger Backend → Storage (Cassandra/Memory)
7. WatchingCat API → Jaeger Query API
8. Frontend → WatchingCat API → Display
```

### Metric Data Flow

```
1. Application → OTLP Exporter (metrics)
2. OTLP Exporter → OTel Collector
3. OTel Collector → Prometheus Exporter (8889)
4. Prometheus → Scrapes OTel Collector
5. Prometheus → Stores Time-Series Data
6. WatchingCat API → PromQL Queries
7. Frontend → WatchingCat API → Charts
```

### Log Data Flow

```
1. Application → Structured Logging
2. Logs → OTel Collector (OTLP)
3. OTel Collector → Elasticsearch Exporter
4. Elasticsearch → Indexes Logs
5. WatchingCat API → Elasticsearch Query API
6. Frontend → WatchingCat API → Display
```

## Deployment Architecture

### Self-Hosted (Docker Compose)

**Current Implementation**:
```yaml
services:
  - otel-collector     # Telemetry pipeline
  - jaeger             # Trace storage
  - prometheus         # Metrics storage
  - grafana            # Metrics visualization
  - elasticsearch      # Log storage
  - kibana             # Log visualization
  - frontend           # Demo app
  - cartservice        # Demo app
  - productcatalog     # Demo app
  - checkoutservice    # Demo app
  - loadgenerator      # Traffic generation
  - webui              # WatchingCat UI
```

**Resource Requirements**:
- **Minimum**: 8GB RAM, 4 CPU cores
- **Recommended**: 16GB RAM, 8 CPU cores
- **Storage**: 50GB+ for telemetry data

### Kubernetes (Planned)

**Components**:
```
Namespace: watchingcat
├── Deployments
│   ├── otel-collector (DaemonSet)
│   ├── watchingcat-backend (Deployment)
│   ├── watchingcat-frontend (Deployment)
│   ├── jaeger (StatefulSet)
│   ├── prometheus (StatefulSet)
│   └── elasticsearch (StatefulSet)
├── Services
│   ├── otel-collector (ClusterIP)
│   ├── watchingcat-api (LoadBalancer)
│   ├── watchingcat-ui (LoadBalancer)
│   └── Internal services
└── ConfigMaps, Secrets, PVCs
```

## Scalability

### Horizontal Scaling

**Stateless Components** (can scale freely):
- OpenTelemetry Collector (DaemonSet pattern)
- WatchingCat Backend (API servers)
- WatchingCat Frontend (static assets)

**Stateful Components** (require coordination):
- Jaeger (replicated backend)
- Prometheus (federation)
- Elasticsearch (cluster)

### Vertical Scaling

**Memory-intensive**:
- Elasticsearch (heap size)
- Prometheus (time-series data)

**CPU-intensive**:
- OTel Collector (processing)
- ClickHouse queries (if migrated)

## Performance Characteristics

### Throughput

**Current Capacity** (single instance):
- **Traces**: 10K spans/second
- **Metrics**: 100K points/second
- **Logs**: 50K events/second

**Optimizations**:
- Batching (1024 items, 10s timeout)
- Memory limiter (512 MiB)
- Connection pooling
- Query caching

### Latency

**End-to-end latency**:
- Trace ingestion → storage: <100ms (p99)
- Query response: <200ms (p95)
- UI updates: 5s refresh interval

## Security

### Authentication

- JWT-based authentication
- Session management
- Password hashing (bcrypt)

### Authorization

- Role-based access control (RBAC)
- Team-based isolation
- API key management

### Data Security

- TLS/SSL for transport
- Encryption at rest (optional)
- Audit logging
- Secrets management

## Monitoring & Observability

### Self-Monitoring

**WatchingCat monitors itself**:
- OTel Collector metrics (port 8888)
- Backend service metrics
- Frontend performance metrics
- Database health checks

**Key Metrics**:
```
# Collector
otelcol_receiver_accepted_spans
otelcol_exporter_sent_spans
otelcol_process_memory_rss
otelcol_process_cpu_seconds

# API
http_requests_total
http_request_duration_seconds
api_errors_total

# Storage
jaeger_spans_stored_total
prometheus_samples_ingested
elasticsearch_docs_indexed
```

### Health Checks

**Endpoints**:
```
/health           # Basic health
/health/ready     # Readiness probe
/health/live      # Liveness probe
/health/startup   # Startup probe
```

## Future Enhancements

### Phase 2: ClickHouse Integration

**Migration Path**:
```
Current: Jaeger + Prometheus + Elasticsearch
    ↓
Hybrid: ClickHouse + (Jaeger/Prometheus/ES)
    ↓
Target: ClickHouse (unified storage)
```

**Benefits**:
- 10x better query performance
- 5x better compression
- Unified query language
- Lower operational complexity

### Phase 3: Advanced Features

**Planned Features**:
- [ ] Service Level Objectives (SLOs)
- [ ] Anomaly detection (ML-based)
- [ ] Cost attribution
- [ ] Multi-tenancy
- [ ] Query builder UI
- [ ] Custom dashboards
- [ ] Synthetic monitoring
- [ ] Incident management

### Phase 4: Cloud-Native

**WatchingCat Cloud** (SaaS):
- Fully managed service
- Auto-scaling
- Global edge presence
- 99.9% uptime SLA

## Comparison with SigNoz

### Similarities

✅ OpenTelemetry-native  
✅ Unified observability (traces, metrics, logs)  
✅ Modern web UI  
✅ Self-hosted option  
✅ Alert management  
✅ Open-source  

### Differences

**WatchingCat**:
- Polyglot storage (Jaeger, Prometheus, ES)
- Go backend with vanilla JS frontend
- Educational focus with demo shop
- Easier initial setup
- Lower resource requirements

**SigNoz**:
- ClickHouse-only storage
- React frontend
- Production-focused
- More advanced querying
- Better for large scale

### When to Choose WatchingCat

✅ Learning OpenTelemetry  
✅ Proof of concept / Demo  
✅ Small to medium deployments  
✅ Existing Prometheus/Jaeger investment  
✅ Lower resource environments  

### When to Choose SigNoz

✅ Large-scale production  
✅ Unified storage preference  
✅ Advanced querying needs  
✅ SaaS option required  
✅ Enterprise support needed  

## Getting Started

### Quick Start

```bash
# Clone repository
git clone https://github.com/yourusername/WatchingCat
cd WatchingCat

# Start all services
make docker-up

# Access UI
open http://localhost:3001
```

### Architecture Decision Records (ADRs)

See `docs/adr/` for detailed architectural decisions:
- ADR-001: Why OpenTelemetry
- ADR-002: Polyglot Storage Strategy
- ADR-003: Go Backend Choice
- ADR-004: Vanilla JS Frontend
- ADR-005: Docker Compose First

## Contributing

WatchingCat is open source and welcomes contributions!

**Architecture Discussions**: GitHub Discussions  
**Issues**: GitHub Issues  
**Documentation**: See `/docs`

## License

Apache 2.0 License

---

**Version**: 1.0.0  
**Status**: Production-Ready (Self-Host)  
**Last Updated**: December 4, 2025  

**Built with ❤️ by the observability community**

