package models

import (
	"fmt"
	"time"
)

// Semester represents an academic semester
type Semester struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Year      int       `gorm:"not null;index:idx_semester" json:"year"`
	Period    int       `gorm:"not null;check:period IN (1,2);index:idx_semester" json:"period"` // 1 or 2
	StartDate time.Time `gorm:"not null" json:"startDate"`
	EndDate   time.Time `gorm:"not null" json:"endDate"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Semester) TableName() string {
	return "semesters"
}

// IsActive checks if the semester is currently active
func (s *Semester) IsActive() bool {
	now := time.Now()
	return now.After(s.StartDate) && now.Before(s.EndDate)
}

// GetDisplayName returns formatted semester name (e.g., "2024.1")
func (s *Semester) GetDisplayName() string {
	return fmt.Sprintf("%d.%d", s.Year, s.Period)
}
