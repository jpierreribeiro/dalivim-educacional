package repository

import (
	"dalivim/internal/models"

	"gorm.io/gorm"
)

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(activity *models.Activity) error {
	return r.db.Create(activity).Error
}

func (r *activityRepository) FindByID(id uint) (*models.Activity, error) {
	var activity models.Activity
	err := r.db.First(&activity, id).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepository) FindByProfessorID(professorID uint) ([]models.Activity, error) {
	var activities []models.Activity
	err := r.db.Where("professor_id = ?", professorID).Find(&activities).Error
	return activities, err
}

func (r *activityRepository) FindByInviteToken(token string) (*models.Activity, error) {
	var activity models.Activity
	err := r.db.Where("invite_token = ?", token).First(&activity).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepository) CountSubmissions(activityID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Submission{}).Where("activity_id = ?", activityID).Count(&count).Error
	return count, err
}
