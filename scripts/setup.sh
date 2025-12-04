#!/bin/bash

# Setup script for OpenTelemetry Observability Platform

set -e

echo "🚀 Setting up OpenTelemetry Observability Platform..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

echo "✅ Go version: $(go version)"

# Check if Docker is installed (optional but recommended)
if command -v docker &> /dev/null; then
    echo "✅ Docker version: $(docker --version)"
else
    echo "⚠️  Docker is not installed. Backend services will not be available."
fi

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
go mod download
go mod tidy

echo "✅ Dependencies downloaded"

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p bin
mkdir -p logs

# Build the project
echo "🔨 Building application..."
make build

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Start backend services: docker-compose up -d"
echo "  2. Run the collector: make run-collector"
echo "  3. Run the application: make run-app"
echo ""
echo "Access UIs:"
echo "  - Jaeger UI: http://localhost:16686"
echo "  - Grafana: http://localhost:3000 (admin/admin)"
echo "  - Kibana: http://localhost:5601"
echo "  - Prometheus: http://localhost:9090"
echo ""

