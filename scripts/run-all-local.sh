#!/bin/bash

# Script to run all services locally for development

set -e

echo "🚀 Starting all OpenTelemetry Demo services locally..."

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Stopping all services..."
    pkill -f "cmd/frontend" || true
    pkill -f "cmd/cartservice" || true
    pkill -f "cmd/productcatalog" || true
    pkill -f "cmd/checkoutservice" || true
    pkill -f "cmd/loadgenerator" || true
    pkill -f "cmd/collector" || true
    echo "✅ All services stopped"
}

trap cleanup EXIT INT TERM

# Start services in background
echo "📊 Starting collector..."
go run cmd/collector/main.go > logs/collector.log 2>&1 &
COLLECTOR_PID=$!
sleep 2

echo "🌐 Starting frontend (port 8080)..."
go run cmd/frontend/main.go > logs/frontend.log 2>&1 &
FRONTEND_PID=$!

echo "🛒 Starting cart service (port 8081)..."
go run cmd/cartservice/main.go > logs/cart.log 2>&1 &
CART_PID=$!

echo "📦 Starting product catalog (port 8082)..."
go run cmd/productcatalog/main.go > logs/product.log 2>&1 &
PRODUCT_PID=$!

echo "💳 Starting checkout service (port 8083)..."
go run cmd/checkoutservice/main.go > logs/checkout.log 2>&1 &
CHECKOUT_PID=$!

sleep 3

echo "🔄 Starting load generator..."
go run cmd/loadgenerator/main.go > logs/loadgen.log 2>&1 &
LOADGEN_PID=$!

echo ""
echo "✅ All services started!"
echo ""
echo "📊 Service URLs:"
echo "  - Frontend: http://localhost:8080"
echo "  - Cart: http://localhost:8081"
echo "  - Product Catalog: http://localhost:8082"
echo "  - Checkout: http://localhost:8083"
echo ""
echo "📝 Logs are in the logs/ directory"
echo "   tail -f logs/frontend.log"
echo ""
echo "Press Ctrl+C to stop all services..."
echo ""

# Wait for user interrupt
wait

