package handler

import (
	"net/http"
	"strconv"

	"dalivim/internal/service"

	"github.com/gin-gonic/gin"
)

type TelemetryHandler struct {
	telemetryService service.TelemetryService
}

func NewTelemetryHandler(telemetryService service.TelemetryService) *TelemetryHandler {
	return &TelemetryHandler{telemetryService: telemetryService}
}

type TelemetryRequest struct {
	ActivityID uint                   `json:"activityId" binding:"required"`
	StudentID  uint                   `json:"studentId" binding:"required"`
	Timestamp  int64                  `json:"timestamp" binding:"required"`
	IsFinal    bool                   `json:"isFinal"`
	Code       string                 `json:"code"`
	Features   map[string]interface{} `json:"features" binding:"required"`
	RawEvents  map[string]interface{} `json:"rawEvents" binding:"required"`
}

func (h *TelemetryHandler) Process(c *gin.Context) {
	var req TelemetryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	analysis, err := h.telemetryService.ProcessTelemetry(
		req.ActivityID,
		req.StudentID,
		req.Timestamp,
		req.IsFinal,
		req.Code,
		req.Features,
		req.RawEvents,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

func (h *TelemetryHandler) GetSubmissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	submissions, err := h.telemetryService.GetSubmissions(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, submissions)
}
