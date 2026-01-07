package models

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Name     string `json:"name"`
	Role     string `gorm:"not null;check:role IN ('professor', 'student')" json:"role"`

	// CAMPOS PARA SISTEMA DE SEMESTRES
	CurrentSemester  int `json:"currentSemester"`  // Semestre atual do aluno (1-10)
	EnrollmentYear   int `json:"enrollmentYear"`   // Ano de matrícula (ex: 2024)
	EnrollmentPeriod int `json:"enrollmentPeriod"` // Período de matrícula (1 ou 2)

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

// UpdateCurrentSemester calculates and updates the student's current semester based on enrollment date
func (u *User) UpdateCurrentSemester(currentYear, currentPeriod int) {
	if u.Role != "student" || u.EnrollmentYear == 0 {
		return
	}

	// Calculate total semesters since enrollment
	yearDiff := currentYear - u.EnrollmentYear
	semestersPassed := yearDiff * 2

	// Adjust for period difference
	if currentPeriod < u.EnrollmentPeriod {
		semestersPassed--
	} else if currentPeriod > u.EnrollmentPeriod {
		semestersPassed++
	}

	// Add 1 because we start counting from semester 1
	u.CurrentSemester = semestersPassed + 1
}
