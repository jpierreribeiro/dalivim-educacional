package models

import "time"

type Activity struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProfessorID uint      `json:"professorId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	TimeLimit   int       `json:"timeLimit"` // in minutes
	InviteToken string    `gorm:"unique" json:"inviteToken"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
