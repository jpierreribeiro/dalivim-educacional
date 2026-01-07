package models

import "time"

// SimilarityDetection stores pairwise similarity comparisons between submissions
type SimilarityDetection struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ActivityID      uint      `gorm:"not null;index:idx_similarity" json:"activityId"`
	SubmissionID1   uint      `gorm:"not null;index:idx_similarity" json:"submissionId1"`
	SubmissionID2   uint      `gorm:"not null;index:idx_similarity" json:"submissionId2"`
	StudentID1      uint      `gorm:"not null" json:"studentId1"`
	StudentID2      uint      `gorm:"not null" json:"studentId2"`
	SimilarityScore float64   `gorm:"not null" json:"similarityScore"` // 0.0 to 1.0
	Algorithm       string    `gorm:"not null" json:"algorithm"`       // "levenshtein", "cosine", "ast"
	IsSuspicious    bool      `gorm:"not null;index" json:"isSuspicious"`
	ClusterID       *uint     `gorm:"index" json:"clusterId,omitempty"` // Group of similar submissions
	CreatedAt       time.Time `json:"createdAt"`

	// Relations
	Submission1 Submission `gorm:"foreignKey:SubmissionID1" json:"submission1,omitempty"`
	Submission2 Submission `gorm:"foreignKey:SubmissionID2" json:"submission2,omitempty"`
}

func (SimilarityDetection) TableName() string {
	return "similarity_detections"
}

// SimilarityCluster groups submissions that are similar to each other
type SimilarityCluster struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ActivityID     uint      `gorm:"not null;index" json:"activityId"`
	ClusterSize    int       `gorm:"not null" json:"clusterSize"`
	AvgSimilarity  float64   `json:"avgSimilarity"`
	SuspicionLevel string    `json:"suspicionLevel"` // "low", "medium", "high"
	CreatedAt      time.Time `json:"createdAt"`
}

func (SimilarityCluster) TableName() string {
	return "similarity_clusters"
}

// SimilarityPair is a helper struct for similarity comparison
type SimilarityPair struct {
	Submission1     Submission `json:"submission1"`
	Submission2     Submission `json:"submission2"`
	SimilarityScore float64    `json:"similarityScore"`
	IsSuspicious    bool       `json:"isSuspicious"`
}
