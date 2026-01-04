package repository

import (
	"dalivim/internal/models"

	"gorm.io/gorm"
)

type telemetryRepository struct {
	db *gorm.DB
}

func NewTelemetryRepository(db *gorm.DB) TelemetryRepository {
	return &telemetryRepository{db: db}
}

func (r *telemetryRepository) Create(telemetry *models.TelemetryData) error {
	return r.db.Create(telemetry).Error
}

func (r *telemetryRepository) FindByActivityAndStudent(activityID, studentID uint) ([]models.TelemetryData, error) {
	var telemetry []models.TelemetryData
	err := r.db.Where("activity_id = ? AND student_id = ?", activityID, studentID).
		Order("timestamp asc").
		Find(&telemetry).Error
	return telemetry, err
}
