package repository

import (
	"dalivim/internal/models"

	"gorm.io/gorm"
)

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Create(submission *models.Submission) error {
	// Marshal signals before saving
	if err := submission.MarshalSignals(); err != nil {
		return err
	}
	return r.db.Create(submission).Error
}

func (r *submissionRepository) FindByActivityID(activityID uint) ([]models.Submission, error) {
	var submissions []models.Submission
	err := r.db.Where("activity_id = ?", activityID).Order("created_at desc").Find(&submissions).Error
	if err != nil {
		return nil, err
	}

	// Unmarshal signals for each submission
	for i := range submissions {
		submissions[i].UnmarshalSignals()
	}

	return submissions, nil
}

func (r *submissionRepository) FindByID(id uint) (*models.Submission, error) {
	var submission models.Submission
	err := r.db.First(&submission, id).Error
	if err != nil {
		return nil, err
	}

	// Unmarshal signals
	submission.UnmarshalSignals()

	return &submission, nil
}
