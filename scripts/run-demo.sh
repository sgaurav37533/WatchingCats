#!/bin/bash

# Demo script to run the complete observability stack

set -e

echo "🎬 Starting OpenTelemetry Observability Demo..."

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Stopping services..."
    pkill -f "cmd/app/main.go" || true
    pkill -f "cmd/collector/main.go" || true
    docker-compose down
    echo "✅ Cleanup complete"
}

trap cleanup EXIT INT TERM

# Start backend services
echo "🐳 Starting backend services (Jaeger, Prometheus, Elasticsearch)..."
docker-compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 10

# Start the collector
echo "📊 Starting OpenTelemetry Collector..."
go run cmd/collector/main.go &
COLLECTOR_PID=$!

sleep 3

# Start the application
echo "🚀 Starting sample application..."
go run cmd/app/main.go &
APP_PID=$!

echo ""
echo "✅ All services started!"
echo ""
echo "📊 Access the UIs:"
echo "  - Jaeger (Traces): http://localhost:16686"
echo "  - Grafana (Metrics): http://localhost:3000"
echo "  - Kibana (Logs): http://localhost:5601"
echo "  - Prometheus: http://localhost:9090"
echo ""
echo "Press Ctrl+C to stop all services..."
echo ""

# Wait for user interrupt
wait

