package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gaurav/watchingcat/internal/dao"
	"go.uber.org/zap"
)

// MetricsHandler handles metrics-related endpoints
type MetricsHandler struct {
	promDAO *dao.PrometheusDAO
	logger  *zap.Logger
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(promDAO *dao.PrometheusDAO, logger *zap.Logger) *MetricsHandler {
	return &MetricsHandler{
		promDAO: promDAO,
		logger:  logger,
	}
}

// GetMetrics returns current metrics
func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	// Return mock metrics for now
	c.JSON(http.StatusOK, gin.H{
		"message": "Use /query or /query_range for actual metrics",
		"example": gin.H{
			"query":       "up",
			"query_range": "rate(http_requests_total[5m])",
		},
	})
}

// Query executes an instant query
func (h *MetricsHandler) Query(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
		Time  int64  `json:"time"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	h.logger.Info("Executing metric query",
		zap.String("query", req.Query),
	)

	var timestamp time.Time
	if req.Time > 0 {
		timestamp = time.Unix(req.Time, 0)
	}

	result, err := h.promDAO.Query(c.Request.Context(), req.Query, timestamp)
	if err != nil {
		h.logger.Error("Failed to execute query", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Query failed",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// QueryRange executes a range query
func (h *MetricsHandler) QueryRange(c *gin.Context) {
	var req struct {
		Query string `json:"query" binding:"required"`
		Start int64  `json:"start" binding:"required"`
		End   int64  `json:"end" binding:"required"`
		Step  int    `json:"step"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	h.logger.Info("Executing metric range query",
		zap.String("query", req.Query),
		zap.Int64("start", req.Start),
		zap.Int64("end", req.End),
	)

	start := time.Unix(req.Start, 0)
	end := time.Unix(req.End, 0)
	step := time.Duration(req.Step) * time.Second
	if req.Step == 0 {
		step = 15 * time.Second
	}

	result, err := h.promDAO.QueryRange(c.Request.Context(), req.Query, start, end, step)
	if err != nil {
		h.logger.Error("Failed to execute range query", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Query failed",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetLabels returns all label names
func (h *MetricsHandler) GetLabels(c *gin.Context) {
	labels, err := h.promDAO.GetLabels(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to fetch labels", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch labels",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": labels,
	})
}

// GetLabelValues returns values for a specific label
func (h *MetricsHandler) GetLabelValues(c *gin.Context) {
	labelName := c.Param("name")

	values, err := h.promDAO.GetLabelValues(c.Request.Context(), labelName)
	if err != nil {
		h.logger.Error("Failed to fetch label values", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch label values",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"label":  labelName,
		"values": values,
	})
}

