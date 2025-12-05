package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gaurav/watchingcat/internal/api/handlers"
	"github.com/gaurav/watchingcat/internal/api/middleware"
	"github.com/gaurav/watchingcat/internal/config"
	"github.com/gaurav/watchingcat/internal/dao"
	"go.uber.org/zap"
)

// NewRouter creates and configures the API router
func NewRouter(
	cfg *config.Config,
	jaegerDAO *dao.JaegerDAO,
	promDAO *dao.PrometheusDAO,
	esDAO *dao.ElasticsearchDAO,
	logger *zap.Logger,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS(cfg.CORS))

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(jaegerDAO, promDAO, esDAO, logger)
	tracesHandler := handlers.NewTracesHandler(jaegerDAO, logger)
	metricsHandler := handlers.NewMetricsHandler(promDAO, logger)
	logsHandler := handlers.NewLogsHandler(esDAO, logger)
	servicesHandler := handlers.NewServicesHandler(jaegerDAO, logger)

	// Serve static files (Frontend)
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/templates/*")

	// Root - serve main UI
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"JaegerURL":     cfg.Jaeger.URL,
			"GrafanaURL":    cfg.Grafana.URL,
			"PrometheusURL": cfg.Prometheus.URL,
			"KibanaURL":     cfg.Kibana.URL,
			"BackendURL":    fmt.Sprintf("http://localhost:%d", cfg.Server.Port),
		})
	})

	// Health check endpoints
	health := router.Group("/health")
	{
		health.GET("", healthHandler.HealthCheck)
		health.GET("/ready", healthHandler.ReadinessCheck)
		health.GET("/live", healthHandler.LivenessCheck)
	}

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Traces endpoints
		traces := v1.Group("/traces")
		{
			traces.GET("", tracesHandler.ListTraces)
			traces.GET("/:id", tracesHandler.GetTrace)
			traces.POST("/search", tracesHandler.SearchTraces)
		}

		// Services endpoints
		services := v1.Group("/services")
		{
			services.GET("", servicesHandler.ListServices)
			services.GET("/:name", servicesHandler.GetService)
			services.GET("/:name/operations", servicesHandler.GetOperations)
		}

		// Metrics endpoints
		metrics := v1.Group("/metrics")
		{
			metrics.GET("", metricsHandler.GetMetrics)
			metrics.POST("/query", metricsHandler.Query)
			metrics.POST("/query_range", metricsHandler.QueryRange)
			metrics.GET("/labels", metricsHandler.GetLabels)
			metrics.GET("/labels/:name/values", metricsHandler.GetLabelValues)
		}

		// Logs endpoints
		logs := v1.Group("/logs")
		{
			logs.GET("", logsHandler.GetLogs)
			logs.POST("/search", logsHandler.SearchLogs)
			logs.GET("/trace/:traceId", logsHandler.GetLogsByTrace)
		}

		// Dashboards endpoints (future)
		dashboards := v1.Group("/dashboards")
		{
			dashboards.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Dashboards API coming in Phase 2"})
			})
		}

		// Alerts endpoints (future)
		alerts := v1.Group("/alerts")
		{
			alerts.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Alerts API coming in Phase 2"})
			})
		}
	}

	// WebSocket endpoint for real-time updates (future)
	router.GET("/ws", func(c *gin.Context) {
		c.JSON(501, gin.H{"message": "WebSocket support coming in Phase 2"})
	})

	return router
}

