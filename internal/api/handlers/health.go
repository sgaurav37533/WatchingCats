package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gaurav/watchingcat/internal/dao"
	"go.uber.org/zap"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	jaegerDAO *dao.JaegerDAO
	promDAO   *dao.PrometheusDAO
	esDAO     *dao.ElasticsearchDAO
	logger    *zap.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(
	jaegerDAO *dao.JaegerDAO,
	promDAO *dao.PrometheusDAO,
	esDAO *dao.ElasticsearchDAO,
	logger *zap.Logger,
) *HealthHandler {
	return &HealthHandler{
		jaegerDAO: jaegerDAO,
		promDAO:   promDAO,
		esDAO:     esDAO,
		logger:    logger,
	}
}

// HealthCheck returns the overall health status
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := gin.H{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"services":  gin.H{},
	}

	services := status["services"].(gin.H)

	// Check Jaeger
	if err := h.jaegerDAO.Ping(ctx); err != nil {
		services["jaeger"] = gin.H{"status": "unhealthy", "error": err.Error()}
	} else {
		services["jaeger"] = gin.H{"status": "healthy"}
	}

	// Check Prometheus
	if err := h.promDAO.Ping(ctx); err != nil {
		services["prometheus"] = gin.H{"status": "unhealthy", "error": err.Error()}
	} else {
		services["prometheus"] = gin.H{"status": "healthy"}
	}

	// Check Elasticsearch
	if err := h.esDAO.Ping(ctx); err != nil {
		services["elasticsearch"] = gin.H{"status": "unhealthy", "error": err.Error()}
	} else {
		services["elasticsearch"] = gin.H{"status": "healthy"}
	}

	// Determine overall status
	allHealthy := true
	for _, svc := range services {
		if svcMap, ok := svc.(gin.H); ok {
			if svcMap["status"] != "healthy" {
				allHealthy = false
				break
			}
		}
	}

	if !allHealthy {
		status["status"] = "degraded"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	c.JSON(http.StatusOK, status)
}

// ReadinessCheck returns readiness status
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check critical services
	if err := h.jaegerDAO.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"reason": "jaeger unavailable",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// LivenessCheck returns liveness status
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

