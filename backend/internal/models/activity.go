package models

import "time"

type Activity struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProfessorID    uint      `gorm:"not null;index" json:"professorId"`
	SemesterID     uint      `gorm:"not null;index" json:"semesterId"`
	TargetSemester int       `gorm:"not null" json:"targetSemester"` // Which student semester (1-10) this activity is for
	Title          string    `gorm:"not null" json:"title"`
	Description    string    `json:"description"`
	Language       string    `gorm:"not null" json:"language"`
	TimeLimit      int       `gorm:"not null" json:"timeLimit"`
	InviteToken    string    `gorm:"unique;not null;index" json:"inviteToken"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Relations
	Semester Semester `gorm:"foreignKey:SemesterID" json:"semester,omitempty"`
}

func (Activity) TableName() string {
	return "activities"
}
