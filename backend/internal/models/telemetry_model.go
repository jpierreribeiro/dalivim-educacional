package models

import "time"

type TelemetryData struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ActivityID uint      `json:"activityId"`
	StudentID  uint      `json:"studentId"`
	Timestamp  int64     `json:"timestamp"`
	IsFinal    bool      `json:"isFinal"`
	Features   string    `gorm:"type:jsonb" json:"features"`
	RawEvents  string    `gorm:"type:jsonb" json:"rawEvents"`
	CreatedAt  time.Time `json:"createdAt"`
}
