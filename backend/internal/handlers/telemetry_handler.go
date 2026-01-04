package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func receiveTelemetry(c *gin.Context) {
	var req struct {
		ActivityID uint                   `json:"activityId"`
		StudentID  uint                   `json:"studentId"`
		Timestamp  int64                  `json:"timestamp"`
		IsFinal    bool                   `json:"isFinal"`
		Code       string                 `json:"code"`
		Features   map[string]interface{} `json:"features"`
		RawEvents  map[string]interface{} `json:"rawEvents"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Analyze behavior
	analysis := analyzeBehavior(req.Features)

	// Save telemetry data
	telemetryJSON, _ := json.Marshal(req.Features)
	eventsJSON, _ := json.Marshal(req.RawEvents)

	telemetry := TelemetryData{
		ActivityID: req.ActivityID,
		StudentID:  req.StudentID,
		Timestamp:  req.Timestamp,
		IsFinal:    req.IsFinal,
		Features:   string(telemetryJSON),
		RawEvents:  string(eventsJSON),
	}
	db.Create(&telemetry)

	// If final submission, create submission record
	if req.IsFinal {
		var student User
		db.First(&student, req.StudentID)

		pasteEventsJSON, _ := json.Marshal(req.RawEvents["pasteEvents"])

		submission := Submission{
			ActivityID:           req.ActivityID,
			StudentID:            req.StudentID,
			StudentName:          student.Name,
			StudentEmail:         student.Email,
			Code:                 req.Code,
			AuthorshipScore:      analysis.AuthorshipScore,
			Confidence:           analysis.Confidence,
			Signals:              analysis.Signals,
			AvgKeystrokeInterval: getFloat(req.Features, "avgKeystrokeInterval"),
			StdKeystrokeInterval: getFloat(req.Features, "stdKeystrokeInterval"),
			PasteEvents:          getInt(req.Features, "pasteEvents"),
			PasteCharRatio:       getFloat(req.Features, "pasteCharRatio"),
			DeleteRatio:          getFloat(req.Features, "deleteRatio"),
			FocusLossCount:       getInt(req.Features, "focusLossCount"),
			LinearEditingScore:   getFloat(req.Features, "linearEditingScore"),
			Burstiness:           getFloat(req.Features, "burstiness"),
			TimeToFirstRun:       getFloat(req.Features, "timeToFirstRun"),
			ExecutionCount:       getInt(req.Features, "executionCount"),
			TotalTime:            getFloat(req.Features, "totalTime"),
			KeystrokeCount:       getInt(req.Features, "totalKeystrokes"),
			PasteEventDetails:    string(pasteEventsJSON),
		}
		db.Create(&submission)
	}

	c.JSON(http.StatusOK, analysis)
}
