#!/bin/bash
# Verification script for OpenTelemetry Collector Dashboard implementation

set -e

echo "🔍 Verifying OpenTelemetry Collector Dashboard Implementation..."
echo ""

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check counter
checks_passed=0
checks_failed=0

# Function to check file exists
check_file() {
    if [ -f "$1" ]; then
        echo -e "${GREEN}✓${NC} File exists: $1"
        ((checks_passed++))
        return 0
    else
        echo -e "${RED}✗${NC} File missing: $1"
        ((checks_failed++))
        return 1
    fi
}

# Function to check directory exists
check_dir() {
    if [ -d "$1" ]; then
        echo -e "${GREEN}✓${NC} Directory exists: $1"
        ((checks_passed++))
        return 0
    else
        echo -e "${RED}✗${NC} Directory missing: $1"
        ((checks_failed++))
        return 1
    fi
}

# Function to check string in file
check_string_in_file() {
    if grep -q "$2" "$1" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Found '$2' in $1"
        ((checks_passed++))
        return 0
    else
        echo -e "${RED}✗${NC} Missing '$2' in $1"
        ((checks_failed++))
        return 1
    fi
}

echo "📁 Checking Required Files..."
echo "================================"
check_file "configs/grafana-dashboards.yaml"
check_file "configs/dashboards/otel-collector-dataflow.json"
check_file "COLLECTOR_DASHBOARD_GUIDE.md"
check_file "IMPLEMENTATION_SUMMARY.md"
echo ""

echo "🔧 Checking Configuration..."
echo "================================"
check_string_in_file "configs/otel-collector-config.yaml" "address: 0.0.0.0:8888"
check_string_in_file "configs/prometheus.yml" "job_name: 'otel-collector'"
check_string_in_file "configs/prometheus.yml" "otel-collector:8888"
check_string_in_file "docker-compose.yaml" "grafana-dashboards.yaml"
check_string_in_file "docker-compose.yaml" "./configs/dashboards"
echo ""

echo "📊 Checking Dashboard JSON..."
echo "================================"
check_string_in_file "configs/dashboards/otel-collector-dataflow.json" "OpenTelemetry Collector Data Flow"
check_string_in_file "configs/dashboards/otel-collector-dataflow.json" "otelcol_receiver_accepted_spans"
check_string_in_file "configs/dashboards/otel-collector-dataflow.json" "otelcol_exporter_sent_spans"
check_string_in_file "configs/dashboards/otel-collector-dataflow.json" "otelcol_process_memory_rss"
check_string_in_file "configs/dashboards/otel-collector-dataflow.json" "Process Metrics"
check_string_in_file "configs/dashboards/otel-collector-dataflow.json" "Traces Pipeline"
check_string_in_file "configs/dashboards/otel-collector-dataflow.json" "Metrics Pipeline"
echo ""

echo "📚 Checking Documentation Updates..."
echo "================================"
check_string_in_file "README.md" "COLLECTOR_DASHBOARD_GUIDE.md"
check_string_in_file "QUICK_REFERENCE.md" "OpenTelemetry Collector Data Flow"
echo ""

# Check if docker-compose.yaml is valid
echo "🐳 Validating Docker Compose..."
echo "================================"
if command -v docker-compose &> /dev/null; then
    if docker-compose config > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC} docker-compose.yaml is valid"
        ((checks_passed++))
    else
        echo -e "${RED}✗${NC} docker-compose.yaml has errors"
        ((checks_failed++))
    fi
else
    echo -e "${YELLOW}⚠${NC} docker-compose not installed, skipping validation"
fi
echo ""

# Check if services are running (optional)
echo "🔌 Checking Services (if running)..."
echo "================================"
if command -v docker-compose &> /dev/null; then
    if docker-compose ps otel-collector 2>&1 | grep -q "Up"; then
        echo -e "${GREEN}✓${NC} OpenTelemetry Collector is running"
        
        # Check if metrics endpoint is accessible
        if curl -s http://localhost:8888/metrics > /dev/null 2>&1; then
            echo -e "${GREEN}✓${NC} Collector metrics endpoint accessible"
            
            # Check for specific metrics
            if curl -s http://localhost:8888/metrics | grep -q "otelcol_receiver_accepted_spans"; then
                echo -e "${GREEN}✓${NC} Collector is exposing expected metrics"
            else
                echo -e "${YELLOW}⚠${NC} Collector running but no span metrics yet (may need traffic)"
            fi
        else
            echo -e "${YELLOW}⚠${NC} Collector running but metrics endpoint not accessible"
        fi
        
        # Check Prometheus
        if docker-compose ps prometheus 2>&1 | grep -q "Up"; then
            echo -e "${GREEN}✓${NC} Prometheus is running"
            
            if curl -s http://localhost:9090/-/healthy > /dev/null 2>&1; then
                echo -e "${GREEN}✓${NC} Prometheus is healthy"
            fi
        fi
        
        # Check Grafana
        if docker-compose ps grafana 2>&1 | grep -q "Up"; then
            echo -e "${GREEN}✓${NC} Grafana is running"
            
            if curl -s http://localhost:3000/api/health > /dev/null 2>&1; then
                echo -e "${GREEN}✓${NC} Grafana is healthy"
                echo -e "${GREEN}→${NC} Dashboard should be available at: http://localhost:3000"
            fi
        fi
    else
        echo -e "${YELLOW}⚠${NC} Services not running. Run 'make docker-up' to start them."
    fi
else
    echo -e "${YELLOW}⚠${NC} docker-compose not installed, skipping service checks"
fi
echo ""

# Summary
echo "================================"
echo "📊 Verification Summary"
echo "================================"
echo -e "Checks passed: ${GREEN}${checks_passed}${NC}"
if [ $checks_failed -gt 0 ]; then
    echo -e "Checks failed: ${RED}${checks_failed}${NC}"
else
    echo -e "Checks failed: ${checks_failed}"
fi
echo ""

if [ $checks_failed -eq 0 ]; then
    echo -e "${GREEN}✅ All checks passed! Implementation is complete.${NC}"
    echo ""
    echo "Next steps:"
    echo "1. Start the stack: make docker-up"
    echo "2. Wait 30 seconds for services to initialize"
    echo "3. Open Grafana: http://localhost:3000 (admin/admin)"
    echo "4. Navigate to: Dashboards → OpenTelemetry → OpenTelemetry Collector Data Flow"
    echo "5. Start load generator: make run-loadgen"
    echo ""
    echo "📖 Read the guide: COLLECTOR_DASHBOARD_GUIDE.md"
    exit 0
else
    echo -e "${RED}❌ Some checks failed. Please review the errors above.${NC}"
    exit 1
fi

