package models

import (
	"encoding/json"
	"time"
)

type Submission struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	ActivityID           uint      `gorm:"not null;index" json:"activityId"`
	StudentID            uint      `gorm:"not null;index" json:"studentId"`
	StudentName          string    `json:"studentName"`
	StudentEmail         string    `json:"studentEmail"`
	Code                 string    `gorm:"type:text" json:"code"`
	AuthorshipScore      float64   `json:"authorshipScore"`
	Confidence           string    `json:"confidence"`
	Signals              string    `gorm:"type:text" json:"-"`
	SignalsArray         []string  `gorm:"-" json:"signals"`
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
	PasteEventDetails    string    `gorm:"type:text" json:"pasteEventDetails"`
	CreatedAt            time.Time `json:"createdAt"`
}

func (Submission) TableName() string {
	return "submissions"
}

func (s *Submission) MarshalSignals() error {
	data, err := json.Marshal(s.SignalsArray)
	if err != nil {
		return err
	}
	s.Signals = string(data)
	return nil
}

func (s *Submission) UnmarshalSignals() error {
	if s.Signals == "" {
		s.SignalsArray = []string{}
		return nil
	}
	return json.Unmarshal([]byte(s.Signals), &s.SignalsArray)
}
