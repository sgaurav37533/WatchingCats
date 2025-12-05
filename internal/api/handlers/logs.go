package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gaurav/watchingcat/internal/dao"
	"go.uber.org/zap"
)

// LogsHandler handles log-related endpoints
type LogsHandler struct {
	esDAO  *dao.ElasticsearchDAO
	logger *zap.Logger
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler(esDAO *dao.ElasticsearchDAO, logger *zap.Logger) *LogsHandler {
	return &LogsHandler{
		esDAO:  esDAO,
		logger: logger,
	}
}

// GetLogs returns logs with pagination
func (h *LogsHandler) GetLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Use POST /api/v1/logs/search for log queries",
	})
}

// SearchLogs searches for logs with filters
func (h *LogsHandler) SearchLogs(c *gin.Context) {
	var req struct {
		Query     string `json:"query"`
		Service   string `json:"service"`
		Level     string `json:"level"`
		StartTime int64  `json:"startTime"`
		EndTime   int64  `json:"endTime"`
		From      int    `json:"from"`
		Size      int    `json:"size"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	h.logger.Info("Searching logs",
		zap.String("service", req.Service),
		zap.String("level", req.Level),
	)

	params := dao.LogSearchParams{
		Query:   req.Query,
		Service: req.Service,
		Level:   req.Level,
		From:    req.From,
		Size:    req.Size,
	}

	if req.StartTime > 0 {
		params.StartTime = time.Unix(req.StartTime, 0)
	}
	if req.EndTime > 0 {
		params.EndTime = time.Unix(req.EndTime, 0)
	}

	if params.Size == 0 {
		params.Size = 100
	}

	result, err := h.esDAO.SearchLogs(c.Request.Context(), params)
	if err != nil {
		h.logger.Error("Failed to search logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Search failed",
		})
		return
	}

	logs := make([]interface{}, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		logs[i] = hit.Source
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": result.Hits.Total.Value,
	})
}

// GetLogsByTrace retrieves logs for a specific trace
func (h *LogsHandler) GetLogsByTrace(c *gin.Context) {
	traceID := c.Param("traceId")

	h.logger.Info("Fetching logs for trace",
		zap.String("trace_id", traceID),
	)

	logs, err := h.esDAO.GetLogsByTraceID(c.Request.Context(), traceID)
	if err != nil {
		h.logger.Error("Failed to fetch logs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch logs",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trace_id": traceID,
		"logs":     logs,
		"total":    len(logs),
	})
}

