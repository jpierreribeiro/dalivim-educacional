package repository

import (
	"dalivim/internal/models"
	"time"

	"gorm.io/gorm"
)

type semesterRepository struct {
	db *gorm.DB
}

func NewSemesterRepository(db *gorm.DB) SemesterRepository {
	return &semesterRepository{db: db}
}

func (r *semesterRepository) Create(semester *models.Semester) error {
	return r.db.Create(semester).Error
}

func (r *semesterRepository) FindActive() (*models.Semester, error) {
	var semester models.Semester
	now := time.Now()
	err := r.db.Where("start_date <= ? AND end_date >= ?", now, now).First(&semester).Error
	if err != nil {
		return nil, err
	}
	return &semester, nil
}

func (r *semesterRepository) FindAll() ([]models.Semester, error) {
	var semesters []models.Semester
	err := r.db.Order("year desc, period desc").Find(&semesters).Error
	return semesters, err
}

func (r *semesterRepository) FindByYearAndPeriod(year, period int) (*models.Semester, error) {
	var semester models.Semester
	err := r.db.Where("year = ? AND period = ?", year, period).First(&semester).Error
	if err != nil {
		return nil, err
	}
	return &semester, nil
}
