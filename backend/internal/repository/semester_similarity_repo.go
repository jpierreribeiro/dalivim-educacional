package repository

import "dalivim/internal/models"

// SemesterRepository handles semester data operations
type SemesterRepository interface {
	Create(semester *models.Semester) error
	FindActive() (*models.Semester, error)
	FindAll() ([]models.Semester, error)
	FindByYearAndPeriod(year, period int) (*models.Semester, error)
}

// SimilarityRepository handles similarity detection data
type SimilarityRepository interface {
	CreateDetection(detection *models.SimilarityDetection) error
	CreateCluster(cluster *models.SimilarityCluster) error
	UpdateDetection(detection *models.SimilarityDetection) error
	FindByActivityID(activityID uint) ([]models.SimilarityDetection, error)
	FindClustersByActivityID(activityID uint) ([]models.SimilarityCluster, error)
	FindByClusterID(clusterID uint) ([]models.SimilarityDetection, error)
	FindSuspiciousByActivityID(activityID uint) ([]models.SimilarityDetection, error)
}
