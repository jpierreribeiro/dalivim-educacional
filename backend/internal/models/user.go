package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"unique;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Name      string    `json:"name"`
	Role      string    `gorm:"not null;check:role IN ('professor', 'student')" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

func (User) TableName() string {
	return "users"
}
