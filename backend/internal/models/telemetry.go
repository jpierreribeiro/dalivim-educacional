package models

import "time"

type TelemetryData struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ActivityID uint      `gorm:"not null;index" json:"activityId"`
	StudentID  uint      `gorm:"not null;index" json:"studentId"`
	Timestamp  int64     `gorm:"not null;index" json:"timestamp"`
	IsFinal    bool      `gorm:"default:false" json:"isFinal"`
	Features   string    `gorm:"type:text" json:"features"`
	RawEvents  string    `gorm:"type:text" json:"rawEvents"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (TelemetryData) TableName() string {
	return "telemetry_data"
}
