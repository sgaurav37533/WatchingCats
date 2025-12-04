package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gaurav/otel-observability/internal/logging"
	"go.uber.org/zap"
)

const serviceName = "webui"
const servicePort = "3001"

type WebUI struct {
	logger *logging.Logger
}

type ServiceStatus struct {
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	Healthy   bool      `json:"healthy"`
	Timestamp time.Time `json:"timestamp"`
}

type DashboardData struct {
	Services       []ServiceStatus
	JaegerURL      string
	GrafanaURL     string
	PrometheusURL  string
	KibanaURL      string
}

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("web/templates/*.html"))
}

func main() {
	logger, err := logging.NewLogger(logging.LoggerConfig{
		Level:       "info",
		Format:      "json",
		ServiceName: serviceName,
	})
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	webui := &WebUI{
		logger: logger,
	}

	mux := http.NewServeMux()
	
	// Web UI pages
	mux.HandleFunc("/", webui.handleDashboard)
	mux.HandleFunc("/api/services", webui.handleServicesAPI)
	mux.HandleFunc("/api/traces", webui.handleTracesAPI)
	mux.HandleFunc("/api/metrics", webui.handleMetricsAPI)
	mux.HandleFunc("/api/logs", webui.handleLogsAPI)
	mux.HandleFunc("/api/loadgen/start", webui.handleStartLoadGen)
	mux.HandleFunc("/api/loadgen/stop", webui.handleStopLoadGen)
	mux.HandleFunc("/health", webui.handleHealth)
	
	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	server := &http.Server{
		Addr:    ":" + servicePort,
		Handler: mux,
	}

	go func() {
		logger.Info("Web UI starting", 
			zap.String("port", servicePort),
			zap.String("url", "http://localhost:"+servicePort),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down Web UI...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func (w *WebUI) handleDashboard(rw http.ResponseWriter, r *http.Request) {
	data := DashboardData{
		Services: []ServiceStatus{
			{Name: "Frontend", URL: "http://localhost:8080", Status: "checking..."},
			{Name: "Cart Service", URL: "http://localhost:8081", Status: "checking..."},
			{Name: "Product Catalog", URL: "http://localhost:8082", Status: "checking..."},
			{Name: "Checkout Service", URL: "http://localhost:8083", Status: "checking..."},
		},
		JaegerURL:     "http://localhost:16686",
		GrafanaURL:    "http://localhost:3000",
		PrometheusURL: "http://localhost:9090",
		KibanaURL:     "http://localhost:5601",
	}

	// Try the new index.html first, fall back to dashboard.html
	templateName := "index.html"
	if err := templates.ExecuteTemplate(rw, templateName, data); err != nil {
		// Try dashboard.html as fallback
		if err2 := templates.ExecuteTemplate(rw, "dashboard.html", data); err2 != nil {
			w.logger.Error("Template error", zap.Error(err), zap.Error(err2))
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func (w *WebUI) handleServicesAPI(rw http.ResponseWriter, r *http.Request) {
	services := []ServiceStatus{
		{Name: "Frontend", URL: "http://localhost:8080", Status: "checking..."},
		{Name: "Cart Service", URL: "http://localhost:8081", Status: "checking..."},
		{Name: "Product Catalog", URL: "http://localhost:8082", Status: "checking..."},
		{Name: "Checkout Service", URL: "http://localhost:8083", Status: "checking..."},
	}

	// Check each service health
	for i := range services {
		healthy := w.checkServiceHealth(services[i].URL + "/health")
		services[i].Healthy = healthy
		if healthy {
			services[i].Status = "healthy"
		} else {
			services[i].Status = "unhealthy"
		}
		services[i].Timestamp = time.Now()
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(services)
}

func (w *WebUI) checkServiceHealth(url string) bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (w *WebUI) handleTracesAPI(rw http.ResponseWriter, r *http.Request) {
	// Proxy request to Jaeger API
	jaegerURL := "http://localhost:16686/api/traces?limit=20&service=frontend"
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(jaegerURL)
	if err != nil {
		w.logger.Error("Failed to fetch traces", zap.Error(err))
		json.NewEncoder(rw).Encode(map[string]string{"error": "Jaeger not available"})
		return
	}
	defer resp.Body.Close()

	rw.Header().Set("Content-Type", "application/json")
	io.Copy(rw, resp.Body)
}

func (w *WebUI) handleMetricsAPI(rw http.ResponseWriter, r *http.Request) {
	// Return mock metrics data
	metrics := map[string]interface{}{
		"request_rate":    150.5,
		"error_rate":      0.05,
		"avg_latency_ms":  245.3,
		"p95_latency_ms":  450.2,
		"total_requests":  125000,
		"total_errors":    6250,
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(metrics)
}

func (w *WebUI) handleLogsAPI(rw http.ResponseWriter, r *http.Request) {
	// Return recent log entries (mock data)
	logs := []map[string]interface{}{
		{
			"timestamp": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
			"level":     "info",
			"service":   "frontend",
			"message":   "Request processed successfully",
			"trace_id":  "abc123def456",
		},
		{
			"timestamp": time.Now().Add(-3 * time.Minute).Format(time.RFC3339),
			"level":     "error",
			"service":   "checkoutservice",
			"message":   "Payment processing failed",
			"trace_id":  "xyz789uvw012",
		},
		{
			"timestamp": time.Now().Add(-1 * time.Minute).Format(time.RFC3339),
			"level":     "warn",
			"service":   "cartservice",
			"message":   "High memory usage detected",
			"trace_id":  "mno345pqr678",
		},
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(logs)
}

func (w *WebUI) handleStartLoadGen(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.logger.Info("Load generator start requested")
	
	// In a real implementation, this would start the load generator process
	response := map[string]interface{}{
		"status":  "started",
		"message": "Load generator started successfully",
		"rate":    "30 requests/min",
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func (w *WebUI) handleStopLoadGen(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.logger.Info("Load generator stop requested")
	
	response := map[string]interface{}{
		"status":  "stopped",
		"message": "Load generator stopped successfully",
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(response)
}

func (w *WebUI) handleHealth(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(map[string]string{"status": "healthy"})
}

