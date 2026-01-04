package repository

import "dalivim/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}

type ActivityRepository interface {
	Create(activity *models.Activity) error
	FindByID(id uint) (*models.Activity, error)
	FindByProfessorID(professorID uint) ([]models.Activity, error)
	FindByInviteToken(token string) (*models.Activity, error)
	CountSubmissions(activityID uint) (int64, error)
}

type SubmissionRepository interface {
	Create(submission *models.Submission) error
	FindByActivityID(activityID uint) ([]models.Submission, error)
	FindByID(id uint) (*models.Submission, error)
}

type TelemetryRepository interface {
	Create(telemetry *models.TelemetryData) error
	FindByActivityAndStudent(activityID, studentID uint) ([]models.TelemetryData, error)
}
