package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gaurav/watchingcat/internal/api"
	"github.com/gaurav/watchingcat/internal/config"
	"github.com/gaurav/watchingcat/internal/dao"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting WatchingCat Backend Service")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded",
		zap.String("jaeger_url", cfg.Jaeger.URL),
		zap.String("prometheus_url", cfg.Prometheus.URL),
		zap.String("elasticsearch_url", cfg.Elasticsearch.URL),
		zap.Int("port", cfg.Server.Port),
	)

	// Initialize DAOs (Data Access Objects)
	logger.Info("Initializing data access objects...")
	
	jaegerDAO := dao.NewJaegerDAO(cfg.Jaeger.URL, logger)
	promDAO := dao.NewPrometheusDAO(cfg.Prometheus.URL, logger)
	esDAO := dao.NewElasticsearchDAO(cfg.Elasticsearch.URL, logger)

	// Test connections
	if err := jaegerDAO.Ping(context.Background()); err != nil {
		logger.Warn("Failed to connect to Jaeger", zap.Error(err))
	} else {
		logger.Info("Connected to Jaeger successfully")
	}

	if err := promDAO.Ping(context.Background()); err != nil {
		logger.Warn("Failed to connect to Prometheus", zap.Error(err))
	} else {
		logger.Info("Connected to Prometheus successfully")
	}

	if err := esDAO.Ping(context.Background()); err != nil {
		logger.Warn("Failed to connect to Elasticsearch", zap.Error(err))
	} else {
		logger.Info("Connected to Elasticsearch successfully")
	}

	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize API router
	logger.Info("Initializing API router...")
	router := api.NewRouter(cfg, jaegerDAO, promDAO, esDAO, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting HTTP server",
			zap.Int("port", cfg.Server.Port),
			zap.String("mode", cfg.Server.Mode),
		)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("WatchingCat Backend is ready!",
		zap.String("url", fmt.Sprintf("http://localhost:%d", cfg.Server.Port)),
		zap.String("health", fmt.Sprintf("http://localhost:%d/health", cfg.Server.Port)),
		zap.String("api", fmt.Sprintf("http://localhost:%d/api/v1", cfg.Server.Port)),
	)

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}

