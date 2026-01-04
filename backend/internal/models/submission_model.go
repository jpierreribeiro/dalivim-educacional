package models

import "time"

type Submission struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	ActivityID           uint      `json:"activityId"`
	StudentID            uint      `json:"studentId"`
	StudentName          string    `json:"studentName"`
	StudentEmail         string    `json:"studentEmail"`
	Code                 string    `json:"code"`
	AuthorshipScore      float64   `json:"authorshipScore"`
	Confidence           string    `json:"confidence"`
	Signals              []string  `gorm:"type:text[]" json:"signals"`
	AvgKeystrokeInterval float64   `json:"avgKeystrokeInterval"`
	StdKeystrokeInterval float64   `json:"stdKeystrokeInterval"`
	PasteEvents          int       `json:"pasteEvents"`
	PasteCharRatio       float64   `json:"pasteCharRatio"`
	DeleteRatio          float64   `json:"deleteRatio"`
	FocusLossCount       int       `json:"focusLossCount"`
	LinearEditingScore   float64   `json:"linearEditingScore"`
	Burstiness           float64   `json:"burstiness"`
	TimeToFirstRun       float64   `json:"timeToFirstRun"`
	ExecutionCount       int       `json:"executionCount"`
	TotalTime            float64   `json:"totalTime"`
	KeystrokeCount       int       `json:"keystrokeCount"`
	PasteEventDetails    string    `gorm:"type:jsonb" json:"pasteEventDetails"`
	CreatedAt            time.Time `json:"createdAt"`
}
