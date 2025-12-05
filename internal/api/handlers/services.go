package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gaurav/watchingcat/internal/dao"
	"go.uber.org/zap"
)

// ServicesHandler handles service-related endpoints
type ServicesHandler struct {
	jaegerDAO *dao.JaegerDAO
	logger    *zap.Logger
}

// NewServicesHandler creates a new services handler
func NewServicesHandler(jaegerDAO *dao.JaegerDAO, logger *zap.Logger) *ServicesHandler {
	return &ServicesHandler{
		jaegerDAO: jaegerDAO,
		logger:    logger,
	}
}

// ListServices lists all services
func (h *ServicesHandler) ListServices(c *gin.Context) {
	h.logger.Info("Fetching services list")

	services, err := h.jaegerDAO.GetServices(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to fetch services", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch services",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"services": services,
		"total":    len(services),
	})
}

// GetService returns information about a specific service
func (h *ServicesHandler) GetService(c *gin.Context) {
	serviceName := c.Param("name")

	h.logger.Info("Fetching service info",
		zap.String("service", serviceName),
	)

	operations, err := h.jaegerDAO.GetOperations(c.Request.Context(), serviceName)
	if err != nil {
		h.logger.Error("Failed to fetch service operations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch service info",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":       serviceName,
		"operations": operations,
	})
}

// GetOperations returns operations for a service
func (h *ServicesHandler) GetOperations(c *gin.Context) {
	serviceName := c.Param("name")

	h.logger.Info("Fetching operations for service",
		zap.String("service", serviceName),
	)

	operations, err := h.jaegerDAO.GetOperations(c.Request.Context(), serviceName)
	if err != nil {
		h.logger.Error("Failed to fetch operations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch operations",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"service":    serviceName,
		"operations": operations,
	})
}

