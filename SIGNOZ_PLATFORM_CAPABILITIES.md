# SigNoz Platform Capabilities Implementation Plan

**Goal**: Implement SigNoz observability platform capabilities in WatchingCat

**Reference**: [SigNoz Platform Repository](https://github.com/SigNoz/signoz)

---

## 🎯 What is SigNoz Platform?

The actual SigNoz platform (not the marketing website) consists of:

1. **Frontend** (React + TypeScript)
   - Dashboard builder
   - Query builder
   - Trace explorer
   - Logs explorer
   - Alerts UI
   - Service map

2. **Backend** (Go)
   - Query Service
   - Alert Manager
   - Rule Engine
   - API Server

3. **Storage** (ClickHouse)
   - Unified storage for traces, metrics, logs
   - Optimized schemas
   - Efficient querying

4. **OpenTelemetry Collector**
   - Custom processors
   - ClickHouse exporters

---

## 📊 SigNoz Platform Capabilities vs WatchingCat

### Current State Comparison

| Capability | SigNoz Platform | WatchingCat Phase 1 | Implementation Plan |
|------------|-----------------|---------------------|---------------------|
| **Core Platform** |
| Distributed Tracing | ✅ Advanced | ✅ Complete | ✅ Already have |
| Metrics Collection | ✅ Advanced | ✅ Basic | 🔨 Phase 2 |
| Log Management | ✅ Advanced | 🔨 Phase 2 | 🔨 Phase 2 |
| Service Map | ✅ Advanced | ✅ Basic D3.js | 🔨 Phase 2 (enhance) |
| **Query & Analysis** |
| ClickHouse SQL | ✅ Yes | ❌ No | 📅 Phase 3 |
| Query Builder UI | ✅ Yes | ❌ No | 🔨 Phase 2 |
| Saved Queries | ✅ Yes | ❌ No | 🔨 Phase 2 |
| Dashboard Builder | ✅ Yes | ⚠️ Basic | 🔨 Phase 2 |
| **Alerting** |
| Alert Rules | ✅ Yes | 🔨 Phase 2 | 🔨 Phase 2 |
| Alert Channels | ✅ Multiple | 🔨 Phase 2 | 🔨 Phase 2 |
| Alert Templates | ✅ Yes | 📅 Phase 3 | 📅 Phase 3 |
| **Advanced Features** |
| SLOs | ✅ Yes | 📅 Phase 3 | 📅 Phase 3 |
| Anomaly Detection | ✅ Yes | 📅 Phase 3 | 📅 Phase 3 |
| Error Tracking | ✅ Yes | ✅ Basic | 🔨 Phase 2 |
| **Backend** |
| Unified API | ✅ Yes | 🔨 Phase 2 | 🔨 Phase 2 |
| GraphQL | ✅ Yes | ❌ No | 📅 Phase 3 |
| WebSocket | ✅ Yes | 🔨 Phase 2 | 🔨 Phase 2 |
| **Authentication** |
| JWT Auth | ✅ Yes | 🔨 Phase 2 | 🔨 Phase 2 |
| OAuth | ✅ Yes | 📅 Phase 4 | 📅 Phase 4 |
| RBAC | ✅ Yes | 📅 Phase 2 | 🔨 Phase 2 |
| **Deployment** |
| Docker Compose | ✅ Yes | ✅ Yes | ✅ Already have |
| Kubernetes | ✅ Helm | 📅 Phase 2 | 📅 Phase 2 |
| Cloud | ✅ SigNoz Cloud | 📅 Phase 4 | 📅 Phase 4 |

---

## 🏗️ SigNoz Architecture Components

### 1. Frontend (React Application)

**Key Features**:
```typescript
// Dashboard Components
- DashboardBuilder
- VisualizationTypes (Time Series, Bar, Pie, Table)
- QueryBuilder
- VariableEditor
- PanelEditor

// Trace Components
- TraceExplorer
- TraceDetail
- SpanDetail
- FlameGraph
- GanttChart

// Logs Components
- LogsExplorer
- LogsViewer
- LogsFilter
- LogsAggregation

// Metrics Components
- MetricsExplorer
- MetricsBrowser
- QueryEditor (PromQL-like)

// Alerts Components
- AlertsList
- AlertRuleEditor
- AlertChannels
- AlertHistory

// Service Map
- ServiceTopology
- ServiceDetail
- DependencyGraph
```

**Technology Stack**:
- React 18+
- TypeScript
- Recharts (charting)
- React Query (data fetching)
- Zustand (state management)
- TailwindCSS (styling)

### 2. Backend (Query Service)

**Architecture**:
```go
package main

import (
    "github.com/SigNoz/signoz/pkg/query-service/app"
    "github.com/SigNoz/signoz/pkg/query-service/dao"
    "github.com/SigNoz/signoz/pkg/query-service/rules"
)

// Key Components:
// 1. ClickHouse DAO (Data Access Object)
// 2. Query Builder
// 3. Aggregation Engine
// 4. Alert Rules Engine
// 5. License Manager
// 6. Usage Tracker
```

**API Endpoints**:
```
GET    /api/v1/traces
GET    /api/v1/traces/:id
POST   /api/v1/traces/search
GET    /api/v1/services
GET    /api/v1/operations

GET    /api/v1/metrics
POST   /api/v1/query_range
POST   /api/v1/query
GET    /api/v1/metrics/labels

GET    /api/v1/logs
POST   /api/v1/logs/search
POST   /api/v1/logs/aggregate
GET    /api/v1/logs/tail

GET    /api/v1/rules
POST   /api/v1/rules
PUT    /api/v1/rules/:id
DELETE /api/v1/rules/:id
POST   /api/v1/rules/test

GET    /api/v1/dashboards
POST   /api/v1/dashboards
PUT    /api/v1/dashboards/:id
```

### 3. ClickHouse Schemas

**Traces Table**:
```sql
CREATE TABLE signoz_traces.distributed_signoz_index_v2 ON CLUSTER cluster
(
    timestamp DateTime64(9) CODEC(DoubleDelta, LZ4),
    traceID FixedString(32) CODEC(ZSTD(1)),
    spanID String CODEC(ZSTD(1)),
    parentSpanID String CODEC(ZSTD(1)),
    serviceName LowCardinality(String) CODEC(ZSTD(1)),
    name LowCardinality(String) CODEC(ZSTD(1)),
    kind Int8 CODEC(T64, ZSTD(1)),
    durationNano UInt64 CODEC(T64, ZSTD(1)),
    statusCode Int16 CODEC(T64, ZSTD(1)),
    ...
) ENGINE = Distributed('cluster', 'signoz_traces', 'signoz_index_v2', cityHash64(traceID))
```

**Metrics Table**:
```sql
CREATE TABLE signoz_metrics.distributed_samples_v2 ON CLUSTER cluster
(
    metric_name LowCardinality(String) CODEC(ZSTD(1)),
    timestamp_ms Int64 CODEC(DoubleDelta, ZSTD(1)),
    value Float64 CODEC(Gorilla, ZSTD(1)),
    fingerprint UInt64 CODEC(Delta, ZSTD(1)),
    ...
) ENGINE = Distributed('cluster', 'signoz_metrics', 'samples_v2', fingerprint)
```

**Logs Table**:
```sql
CREATE TABLE signoz_logs.distributed_logs ON CLUSTER cluster
(
    timestamp DateTime64(9) CODEC(DoubleDelta, LZ4),
    id String CODEC(ZSTD(1)),
    trace_id String CODEC(ZSTD(1)),
    span_id String CODEC(ZSTD(1)),
    severity_text LowCardinality(String) CODEC(ZSTD(1)),
    body String CODEC(ZSTD(1)),
    ...
) ENGINE = Distributed('cluster', 'signoz_logs', 'logs', cityHash64(id))
```

---

## 🎯 Implementation Plan for WatchingCat

### Phase 2A: Core Backend (Weeks 1-2)

**Goal**: Build unified backend matching SigNoz Query Service

#### Week 1: Foundation
```bash
# Project Structure
cmd/backend/
├── main.go                 # Main entry
├── config/
│   └── config.go          # Configuration
internal/
├── api/
│   ├── router.go          # API router
│   ├── handlers/
│   │   ├── traces.go      # Trace handlers
│   │   ├── metrics.go     # Metrics handlers
│   │   └── logs.go        # Logs handlers
│   └── middleware/
│       ├── auth.go        # Authentication
│       └── cors.go        # CORS
├── dao/
│   ├── jaeger.go          # Jaeger DAO
│   ├── prometheus.go      # Prometheus DAO
│   └── elasticsearch.go   # ES DAO
└── models/
    ├── trace.go
    ├── metric.go
    └── log.go
```

**Tasks**:
- [ ] Setup Go project structure
- [ ] Implement Jaeger client
- [ ] Implement Prometheus client
- [ ] Implement Elasticsearch client
- [ ] Build unified API layer
- [ ] Add authentication middleware

#### Week 2: Advanced Features
- [ ] Query builder for complex queries
- [ ] Aggregation engine
- [ ] Caching layer (Redis)
- [ ] WebSocket support
- [ ] API documentation (Swagger)

### Phase 2B: Enhanced Frontend (Weeks 3-4)

**Goal**: Build advanced UI components matching SigNoz

#### Week 3: Query & Dashboard Builder
```javascript
// New Components
components/
├── QueryBuilder/
│   ├── QueryEditor.js
│   ├── FilterBuilder.js
│   └── AggregationBuilder.js
├── Dashboard/
│   ├── DashboardBuilder.js
│   ├── PanelEditor.js
│   └── VariableEditor.js
└── Visualization/
    ├── TimeSeriesChart.js
    ├── BarChart.js
    ├── PieChart.js
    └── TableView.js
```

**Tasks**:
- [ ] Build query builder UI
- [ ] Implement dashboard builder
- [ ] Add panel editor
- [ ] Create visualization library
- [ ] Add saved queries

#### Week 4: Logs & Alerts
- [ ] Logs explorer UI
- [ ] Logs filtering and search
- [ ] Alert rule editor
- [ ] Alert channels configuration
- [ ] Alert history view

### Phase 2C: Integration (Week 5)

**Tasks**:
- [ ] Connect frontend to backend API
- [ ] Real data in all visualizations
- [ ] WebSocket real-time updates
- [ ] End-to-end testing
- [ ] Performance optimization

---

## 🚀 Specific SigNoz Features to Implement

### 1. Advanced Trace Explorer

**SigNoz Capabilities**:
- Complex trace filtering
- Trace aggregation
- Trace comparison
- Custom span attributes filtering
- Tag-based search

**Implementation**:
```javascript
// TraceExplorer Component
const TraceExplorer = () => {
  return (
    <div>
      <FilterBuilder
        filters={[
          { key: 'service.name', operator: '=', value: 'frontend' },
          { key: 'duration', operator: '>', value: 1000 },
          { key: 'status.code', operator: '=', value: 'ERROR' }
        ]}
      />
      <AggregationView
        groupBy="service.name"
        aggregation="P95(duration)"
      />
      <TraceList traces={filteredTraces} />
    </div>
  )
}
```

### 2. Query Builder

**SigNoz Capabilities**:
- Visual query construction
- Metric/span/log queries
- Function library (avg, max, rate, etc.)
- Group by multiple dimensions

**Implementation**:
```javascript
// QueryBuilder Component
const QueryBuilder = () => {
  const [query, setQuery] = useState({
    metric: 'http_requests_total',
    filters: [],
    groupBy: ['service', 'method'],
    aggregation: 'rate',
    timeRange: '1h'
  })
  
  return (
    <QueryEditorUI
      query={query}
      onChange={setQuery}
      onExecute={executeQuery}
    />
  )
}
```

### 3. Dashboard Builder

**SigNoz Capabilities**:
- Drag-and-drop panels
- Multiple visualization types
- Variables and templating
- Panel linking
- Time range picker

**Implementation**:
```javascript
// DashboardBuilder Component
const DashboardBuilder = () => {
  const [panels, setPanels] = useState([])
  
  return (
    <GridLayout
      layout={panels}
      onLayoutChange={setPanels}
      draggableHandle=".panel-header"
    >
      {panels.map(panel => (
        <Panel
          key={panel.id}
          type={panel.type}
          query={panel.query}
          options={panel.options}
        />
      ))}
    </GridLayout>
  )
}
```

### 4. Logs Explorer

**SigNoz Capabilities**:
- Full-text search
- Field filtering
- Log patterns
- Histogram view
- Trace correlation
- Live tail

**Implementation**:
```javascript
// LogsExplorer Component
const LogsExplorer = () => {
  return (
    <div>
      <LogsFilter
        fields={['level', 'service', 'message']}
        search={searchQuery}
      />
      <LogsHistogram data={logStats} />
      <LogsTable
        logs={filteredLogs}
        onTraceClick={viewTrace}
        liveTail={isLiveTail}
      />
    </div>
  )
}
```

### 5. Alert Rules Engine

**SigNoz Capabilities**:
- Threshold alerts
- Composite alerts
- Alert templates
- Multiple channels
- Alert history

**Implementation**:
```go
// Alert Rule Engine (Go)
type AlertRule struct {
    ID          string
    Name        string
    Query       string
    Condition   string  // e.g., "value > 100"
    Duration    string  // e.g., "5m"
    Labels      map[string]string
    Annotations map[string]string
}

func (a *AlertManager) EvaluateRules(ctx context.Context) {
    for _, rule := range a.rules {
        result := a.executeQuery(rule.Query)
        if a.checkCondition(result, rule.Condition) {
            a.fireAlert(rule)
        }
    }
}
```

---

## 📊 Technology Stack Alignment

### WatchingCat → SigNoz Alignment

| Component | WatchingCat Current | SigNoz | Recommendation |
|-----------|---------------------|---------|----------------|
| **Frontend** | Vanilla JS | React + TS | Keep Vanilla or migrate to React |
| **Charts** | Chart.js | Recharts | Keep Chart.js (lighter) |
| **Topology** | D3.js | D3.js | ✅ Already aligned |
| **Backend** | Go (planned) | Go | ✅ Already aligned |
| **Storage** | Polyglot | ClickHouse | Migrate in Phase 3 |
| **Query** | PromQL/JaegerQL | ClickHouse SQL | Gradual transition |
| **State** | Local | Zustand | Add state management |

---

## 🎯 Feature Priority Matrix

### High Priority (Phase 2 - Immediate)
1. ✅ **Unified Backend API** - Critical foundation
2. ✅ **Real Data Integration** - Remove all mocks
3. ✅ **Query Builder UI** - Essential UX
4. ✅ **Alert Management** - Production requirement
5. ✅ **Logs Explorer** - Complete observability

### Medium Priority (Phase 3)
6. **Dashboard Builder** - Custom dashboards
7. **ClickHouse Migration** - Performance boost
8. **Advanced Filtering** - Power user features
9. **SLO Tracking** - Reliability focus
10. **Saved Queries** - Productivity

### Lower Priority (Phase 4)
11. **Multi-tenancy** - Enterprise
12. **OAuth Integration** - Enterprise
13. **Advanced Analytics** - ML features
14. **Cloud Deployment** - SaaS

---

## 🔧 Implementation Code Examples

### 1. Unified Backend API (Go)

```go
// cmd/backend/main.go
package main

import (
    "context"
    "fmt"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/watchingcat/internal/dao"
    "github.com/watchingcat/internal/api"
)

func main() {
    // Initialize DAOs
    jaegerDAO := dao.NewJaegerDAO("http://localhost:16686")
    promDAO := dao.NewPrometheusDAO("http://localhost:9090")
    esDAO := dao.NewElasticsearchDAO("http://localhost:9200")
    
    // Initialize API server
    router := gin.Default()
    
    // Setup routes
    api.SetupTraceRoutes(router, jaegerDAO)
    api.SetupMetricsRoutes(router, promDAO)
    api.SetupLogsRoutes(router, esDAO)
    api.SetupAlertsRoutes(router)
    
    // Start server
    router.Run(":8090")
}
```

### 2. Query Builder (JavaScript)

```javascript
// components/QueryBuilder.js
class QueryBuilder {
    constructor() {
        this.query = {
            metric: '',
            filters: [],
            groupBy: [],
            aggregation: null,
            timeRange: '1h'
        }
    }
    
    addFilter(key, operator, value) {
        this.query.filters.push({ key, operator, value })
        return this
    }
    
    setAggregation(fn, field) {
        this.query.aggregation = { fn, field }
        return this
    }
    
    build() {
        // Convert to PromQL or ClickHouse SQL
        return this.toPromQL()
    }
    
    toPromQL() {
        let q = this.query.metric
        if (this.query.filters.length > 0) {
            const filters = this.query.filters
                .map(f => `${f.key}${f.operator}"${f.value}"`)
                .join(',')
            q += `{${filters}}`
        }
        if (this.query.aggregation) {
            q = `${this.query.aggregation.fn}(${q})`
        }
        return q
    }
}
```

### 3. Dashboard Builder (JavaScript)

```javascript
// components/DashboardBuilder.js
const DashboardBuilder = () => {
    const [panels, setPanels] = useState([
        {
            id: 'panel-1',
            type: 'timeseries',
            title: 'Request Rate',
            query: 'rate(http_requests_total[5m])',
            gridPos: { x: 0, y: 0, w: 12, h: 8 }
        }
    ])
    
    const addPanel = (type) => {
        const newPanel = {
            id: `panel-${Date.now()}`,
            type,
            title: 'New Panel',
            query: '',
            gridPos: { x: 0, y: 0, w: 6, h: 8 }
        }
        setPanels([...panels, newPanel])
    }
    
    const updatePanel = (id, updates) => {
        setPanels(panels.map(p => 
            p.id === id ? { ...p, ...updates } : p
        ))
    }
    
    return (
        <div className="dashboard-builder">
            <Toolbar onAddPanel={addPanel} />
            <GridLayout layout={panels} onLayoutChange={updateLayout}>
                {panels.map(panel => (
                    <Panel
                        key={panel.id}
                        {...panel}
                        onUpdate={(updates) => updatePanel(panel.id, updates)}
                    />
                ))}
            </GridLayout>
        </div>
    )
}
```

---

## 📈 Success Metrics

### Phase 2 Completion Criteria

**Backend**:
- [ ] All API endpoints implemented (30+)
- [ ] Response time <200ms (p95)
- [ ] Test coverage >80%
- [ ] API documentation complete

**Frontend**:
- [ ] Query builder functional
- [ ] Dashboard builder working
- [ ] Logs explorer operational
- [ ] Alert UI complete
- [ ] Real data in all views

**Integration**:
- [ ] End-to-end tests passing
- [ ] Performance benchmarks met
- [ ] Zero critical bugs
- [ ] Documentation updated

---

## 🎯 Next Steps

### Week 1: Start Backend Implementation
1. Review [BACKEND_IMPLEMENTATION_GUIDE.md](BACKEND_IMPLEMENTATION_GUIDE.md)
2. Set up Go project structure
3. Implement Jaeger client
4. Create basic API endpoints
5. Test with existing frontend

### Week 2: Advanced Backend
1. Add Prometheus client
2. Add Elasticsearch client
3. Build query builder
4. Implement WebSocket
5. Add authentication

### Week 3-4: Frontend Enhancement
1. Build query builder UI
2. Create dashboard builder
3. Add logs explorer
4. Implement alert UI
5. Real-time updates

### Week 5: Integration & Polish
1. Connect all components
2. End-to-end testing
3. Performance optimization
4. Documentation
5. Deploy staging

---

## 📚 Resources

### SigNoz Resources
- [SigNoz Platform Repo](https://github.com/SigNoz/signoz)
- [SigNoz Docs](https://signoz.io/docs/)
- [SigNoz Architecture](https://signoz.io/docs/architecture/)
- [SigNoz Blog](https://signoz.io/blog/)

### WatchingCat Resources
- [BACKEND_IMPLEMENTATION_GUIDE.md](BACKEND_IMPLEMENTATION_GUIDE.md)
- [WATCHINGCAT_ARCHITECTURE.md](WATCHINGCAT_ARCHITECTURE.md)
- [PRODUCT_ROADMAP.md](PRODUCT_ROADMAP.md)

---

## 🎉 Summary

**We're implementing SigNoz PLATFORM capabilities** (not the marketing website):

✅ **What We Have** (Phase 1):
- Basic observability platform
- Beautiful UI with trace viewer
- Demo applications
- Excellent documentation

🔨 **What We're Building** (Phase 2):
- Unified backend API (like SigNoz Query Service)
- Advanced query builder
- Dashboard builder
- Logs explorer
- Alert management

📅 **Future** (Phase 3-4):
- ClickHouse integration (optional)
- Advanced analytics
- Enterprise features
- Cloud deployment

**Timeline**: 5 weeks for Phase 2
**Result**: Feature parity with SigNoz core platform

---

**Ready to start Phase 2 implementation!** 🚀

---

**Last Updated**: December 4, 2025  
**Status**: Implementation Plan Ready  
**Reference**: SigNoz Platform (not marketing site)

