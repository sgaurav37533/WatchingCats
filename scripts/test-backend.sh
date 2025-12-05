#!/bin/bash

echo "🧪 Testing WatchingCat Backend..."
echo ""

# Check if backend directory exists
if [ ! -d "cmd/backend" ]; then
    echo "❌ Backend directory not found"
    exit 1
fi

echo "✅ Backend directory found"
echo ""

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    exit 1
fi

echo "✅ Go is installed: $(go version)"
echo ""

# Verify backend files
echo "📁 Verifying backend files..."
files=(
    "cmd/backend/main.go"
    "internal/config/config.go"
    "internal/dao/jaeger.go"
    "internal/dao/prometheus.go"
    "internal/dao/elasticsearch.go"
    "internal/api/router.go"
    "configs/backend-config.yaml"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✅ $file"
    else
        echo "  ❌ $file (missing)"
        exit 1
    fi
done

echo ""
echo "✅ All backend files present"
echo ""

echo "📊 Backend Statistics:"
echo "  Files: $(find internal cmd/backend -name '*.go' 2>/dev/null | wc -l | xargs)"
echo "  Lines: $(find internal cmd/backend -name '*.go' -exec cat {} \; 2>/dev/null | wc -l | xargs)"
echo ""

echo "✅ Backend foundation is ready!"
echo ""
echo "Next steps:"
echo "  1. Start Docker services: make docker-up"
echo "  2. Run backend: make run-backend"
echo "  3. Test API: curl http://localhost:8090/health"
echo ""

