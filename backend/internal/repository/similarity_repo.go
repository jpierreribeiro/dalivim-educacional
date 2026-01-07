package repository

import (
	"dalivim/internal/models"

	"gorm.io/gorm"
)

type similarityRepository struct {
	db *gorm.DB
}

func NewSimilarityRepository(db *gorm.DB) SimilarityRepository {
	return &similarityRepository{db: db}
}

func (r *similarityRepository) CreateDetection(detection *models.SimilarityDetection) error {
	return r.db.Create(detection).Error
}

func (r *similarityRepository) CreateCluster(cluster *models.SimilarityCluster) error {
	return r.db.Create(cluster).Error
}

func (r *similarityRepository) UpdateDetection(detection *models.SimilarityDetection) error {
	return r.db.Save(detection).Error
}

func (r *similarityRepository) FindByActivityID(activityID uint) ([]models.SimilarityDetection, error) {
	var detections []models.SimilarityDetection
	err := r.db.Where("activity_id = ?", activityID).
		Preload("Submission1").
		Preload("Submission2").
		Order("similarity_score desc").
		Find(&detections).Error
	return detections, err
}

func (r *similarityRepository) FindClustersByActivityID(activityID uint) ([]models.SimilarityCluster, error) {
	var clusters []models.SimilarityCluster
	err := r.db.Where("activity_id = ?", activityID).
		Order("avg_similarity desc").
		Find(&clusters).Error
	return clusters, err
}

func (r *similarityRepository) FindByClusterID(clusterID uint) ([]models.SimilarityDetection, error) {
	var detections []models.SimilarityDetection
	err := r.db.Where("cluster_id = ?", clusterID).Find(&detections).Error
	return detections, err
}

func (r *similarityRepository) FindSuspiciousByActivityID(activityID uint) ([]models.SimilarityDetection, error) {
	var detections []models.SimilarityDetection
	err := r.db.Where("activity_id = ? AND is_suspicious = ?", activityID, true).
		Preload("Submission1").
		Preload("Submission2").
		Order("similarity_score desc").
		Find(&detections).Error
	return detections, err
}
