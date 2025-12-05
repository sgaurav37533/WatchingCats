package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gaurav/watchingcat/internal/dao"
	"go.uber.org/zap"
)

// TracesHandler handles trace-related endpoints
type TracesHandler struct {
	jaegerDAO *dao.JaegerDAO
	logger    *zap.Logger
}

// NewTracesHandler creates a new traces handler
func NewTracesHandler(jaegerDAO *dao.JaegerDAO, logger *zap.Logger) *TracesHandler {
	return &TracesHandler{
		jaegerDAO: jaegerDAO,
		logger:    logger,
	}
}

// ListTraces lists traces based on query parameters
func (h *TracesHandler) ListTraces(c *gin.Context) {
	service := c.DefaultQuery("service", "frontend")
	limitStr := c.DefaultQuery("limit", "20")
	operation := c.Query("operation")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	h.logger.Info("Listing traces",
		zap.String("service", service),
		zap.Int("limit", limit),
	)

	params := dao.SearchParams{
		ServiceName: service,
		Operation:   operation,
		Limit:       limit,
	}

	traces, err := h.jaegerDAO.SearchTraces(c.Request.Context(), params)
	if err != nil {
		h.logger.Error("Failed to search traces", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch traces",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"traces": traces,
		"total":  len(traces),
	})
}

// GetTrace retrieves a single trace by ID
func (h *TracesHandler) GetTrace(c *gin.Context) {
	traceID := c.Param("id")

	h.logger.Info("Fetching trace",
		zap.String("trace_id", traceID),
	)

	trace, err := h.jaegerDAO.GetTrace(c.Request.Context(), traceID)
	if err != nil {
		h.logger.Error("Failed to fetch trace",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Trace not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trace": trace,
	})
}

// SearchTraces searches for traces with advanced filters
func (h *TracesHandler) SearchTraces(c *gin.Context) {
	var req struct {
		Service     string            `json:"service" binding:"required"`
		Operation   string            `json:"operation"`
		MinDuration string            `json:"minDuration"`
		MaxDuration string            `json:"maxDuration"`
		Limit       int               `json:"limit"`
		StartTime   int64             `json:"startTime"`
		EndTime     int64             `json:"endTime"`
		Tags        map[string]string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	h.logger.Info("Searching traces",
		zap.String("service", req.Service),
		zap.String("operation", req.Operation),
	)

	params := dao.SearchParams{
		ServiceName: req.Service,
		Operation:   req.Operation,
		MinDuration: req.MinDuration,
		MaxDuration: req.MaxDuration,
		Limit:       req.Limit,
		Start:       req.StartTime,
		End:         req.EndTime,
		Tags:        req.Tags,
	}

	if params.Limit == 0 {
		params.Limit = 20
	}

	traces, err := h.jaegerDAO.SearchTraces(c.Request.Context(), params)
	if err != nil {
		h.logger.Error("Failed to search traces", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"traces": traces,
		"total":  len(traces),
	})
}

