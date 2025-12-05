.PHONY: all build run test clean help

# Variables
SERVICES=frontend cartservice productcatalog checkoutservice loadgenerator webui backend
GOCMD=go
GOBUILD=CGO_ENABLED=0 $(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOMOD=$(GOCMD) mod

all: build

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build all microservices
build:
	@echo "Building all services..."
	@mkdir -p bin
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		$(GOBUILD) -o bin/$$service cmd/$$service/main.go; \
	done
	@echo "Building collector..."; \
	$(GOBUILD) -o bin/collector cmd/collector/main.go
	@echo "✅ All services built successfully"

## run-backend: Run the unified backend API server (NEW - PHASE 2)
run-backend:
	@echo "🚀 Starting WatchingCat Backend API..."
	@echo "API: http://localhost:8090"
	@echo "Health: http://localhost:8090/health"
	@echo "API Docs: http://localhost:8090/api/v1"
	$(GOCMD) run cmd/backend/main.go

## run-webui: Run the Web UI dashboard (PRIMARY INTERFACE)
run-webui:
	@echo "🌐 Starting Web UI Dashboard..."
	@echo "Open http://localhost:3001 in your browser"
	$(GOCMD) run cmd/webui/main.go

## run-frontend: Run the frontend service
run-frontend:
	$(GOCMD) run cmd/frontend/main.go

## run-cart: Run the cart service
run-cart:
	$(GOCMD) run cmd/cartservice/main.go

## run-product: Run the product catalog service
run-product:
	$(GOCMD) run cmd/productcatalog/main.go

## run-checkout: Run the checkout service
run-checkout:
	$(GOCMD) run cmd/checkoutservice/main.go

## run-loadgen: Run the load generator
run-loadgen:
	$(GOCMD) run cmd/loadgenerator/main.go

## run-collector: Run the collector service
run-collector:
	$(GOCMD) run cmd/collector/main.go

## run-all-local: Run all services locally
run-all-local:
	@echo "Starting all services..."
	@./scripts/run-all-local.sh

## docker-up: Start all services with Docker Compose
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d
	@echo "✅ All services started"
	@echo ""
	@echo "🌐 PRIMARY: Web Dashboard → http://localhost:3001"
	@echo ""
	@echo "Other UIs:"
	@echo "  - Frontend API: http://localhost:8080"
	@echo "  - Jaeger: http://localhost:16686"
	@echo "  - Grafana: http://localhost:3000 (admin/admin)"
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Kibana: http://localhost:5601"

## docker-down: Stop all Docker Compose services
docker-down:
	docker-compose down

## docker-logs: Show logs from all services
docker-logs:
	docker-compose logs -f

## status: Check status of all services
status:
	@echo "Checking service status..."
	@curl -s http://localhost:8090/health && echo "✅ Backend API: healthy" || echo "❌ Backend API: down"
	@curl -s http://localhost:3001/health && echo "✅ Web UI: healthy" || echo "❌ Web UI: down"
	@curl -s http://localhost:8080/health && echo "✅ Frontend: healthy" || echo "❌ Frontend: down"
	@curl -s http://localhost:8081/health && echo "✅ Cart: healthy" || echo "❌ Cart: down"
	@curl -s http://localhost:8082/health && echo "✅ Product Catalog: healthy" || echo "❌ Product Catalog: down"
	@curl -s http://localhost:8083/health && echo "✅ Checkout: healthy" || echo "❌ Checkout: down"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.txt
	@echo "✅ Clean complete"

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "✅ Format complete"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "✅ Dependencies downloaded"

## verify-dashboard: Verify OpenTelemetry Collector Dashboard implementation
verify-dashboard:
	@./scripts/verify-dashboard.sh
