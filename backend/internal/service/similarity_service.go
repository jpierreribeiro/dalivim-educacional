package service

import (
	"dalivim/internal/models"
	"dalivim/internal/repository"
	"strings"
	"unicode"
)

type SimilarityService interface {
	DetectSimilarities(activityID uint) ([]models.SimilarityDetection, error)
	GetSimilaritiesForActivity(activityID uint) ([]models.SimilarityDetection, error)
	GetClustersForActivity(activityID uint) ([]SimilarityCluster, error)
}

type SimilarityCluster struct {
	models.SimilarityCluster
	Submissions []models.Submission `json:"submissions"`
}

type similarityService struct {
	similarityRepo repository.SimilarityRepository
	submissionRepo repository.SubmissionRepository
}

func NewSimilarityService(
	similarityRepo repository.SimilarityRepository,
	submissionRepo repository.SubmissionRepository,
) SimilarityService {
	return &similarityService{
		similarityRepo: similarityRepo,
		submissionRepo: submissionRepo,
	}
}

// DetectSimilarities performs pairwise comparison of all submissions in an activity
func (s *similarityService) DetectSimilarities(activityID uint) ([]models.SimilarityDetection, error) {
	// Get all submissions for the activity
	submissions, err := s.submissionRepo.FindByActivityID(activityID)
	if err != nil {
		return nil, err
	}

	var detections []models.SimilarityDetection

	// Compare each pair of submissions
	for i := 0; i < len(submissions); i++ {
		for j := i + 1; j < len(submissions); j++ {
			sub1 := submissions[i]
			sub2 := submissions[j]

			// Calculate similarity score
			score := calculateCodeSimilarity(sub1.Code, sub2.Code)

			// Create detection record
			detection := models.SimilarityDetection{
				ActivityID:      activityID,
				SubmissionID1:   sub1.ID,
				SubmissionID2:   sub2.ID,
				StudentID1:      sub1.StudentID,
				StudentID2:      sub2.StudentID,
				SimilarityScore: score,
				Algorithm:       "levenshtein_normalized",
				IsSuspicious:    score > 0.75, // Threshold for suspicion
			}

			// Save detection
			if err := s.similarityRepo.CreateDetection(&detection); err != nil {
				continue // Log error but continue with other comparisons
			}

			detections = append(detections, detection)
		}
	}

	// Cluster similar submissions
	s.clusterSimilarities(activityID, detections)

	return detections, nil
}

func (s *similarityService) GetSimilaritiesForActivity(activityID uint) ([]models.SimilarityDetection, error) {
	return s.similarityRepo.FindByActivityID(activityID)
}

func (s *similarityService) GetClustersForActivity(activityID uint) ([]SimilarityCluster, error) {
	clusters, err := s.similarityRepo.FindClustersByActivityID(activityID)
	if err != nil {
		return nil, err
	}

	var result []SimilarityCluster
	for _, cluster := range clusters {
		// Get all submissions in this cluster
		detections, _ := s.similarityRepo.FindByClusterID(cluster.ID)

		submissionIDs := make(map[uint]bool)
		for _, det := range detections {
			submissionIDs[det.SubmissionID1] = true
			submissionIDs[det.SubmissionID2] = true
		}

		var submissions []models.Submission
		for id := range submissionIDs {
			sub, _ := s.submissionRepo.FindByID(id)
			if sub != nil {
				submissions = append(submissions, *sub)
			}
		}

		result = append(result, SimilarityCluster{
			SimilarityCluster: cluster,
			Submissions:       submissions,
		})
	}

	return result, nil
}

// clusterSimilarities groups similar submissions using simple connected components
func (s *similarityService) clusterSimilarities(activityID uint, detections []models.SimilarityDetection) {
	// Build adjacency list of suspicious similarities
	graph := make(map[uint][]uint)

	for _, det := range detections {
		if det.IsSuspicious {
			graph[det.SubmissionID1] = append(graph[det.SubmissionID1], det.SubmissionID2)
			graph[det.SubmissionID2] = append(graph[det.SubmissionID2], det.SubmissionID1)
		}
	}

	// Find connected components (clusters)
	visited := make(map[uint]bool)
	clusterNum := 0

	for node := range graph {
		if !visited[node] {
			// BFS to find all nodes in this cluster
			cluster := []uint{}
			queue := []uint{node}
			visited[node] = true

			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]
				cluster = append(cluster, current)

				for _, neighbor := range graph[current] {
					if !visited[neighbor] {
						visited[neighbor] = true
						queue = append(queue, neighbor)
					}
				}
			}

			// Only create cluster if it has 2+ submissions
			if len(cluster) >= 2 {
				clusterNum++
				s.saveCluster(activityID, uint(clusterNum), cluster, detections)
			}
		}
	}
}

func (s *similarityService) saveCluster(activityID uint, clusterID uint, submissionIDs []uint, detections []models.SimilarityDetection) {
	// Calculate average similarity within cluster
	var totalSim float64
	count := 0

	for _, det := range detections {
		if contains(submissionIDs, det.SubmissionID1) && contains(submissionIDs, det.SubmissionID2) {
			totalSim += det.SimilarityScore
			count++
		}
	}

	avgSim := totalSim / float64(count)

	// Determine suspicion level
	suspicionLevel := "low"
	if avgSim > 0.9 {
		suspicionLevel = "high"
	} else if avgSim > 0.8 {
		suspicionLevel = "medium"
	}

	cluster := &models.SimilarityCluster{
		ActivityID:     activityID,
		ClusterSize:    len(submissionIDs),
		AvgSimilarity:  avgSim,
		SuspicionLevel: suspicionLevel,
	}

	// Save cluster
	if err := s.similarityRepo.CreateCluster(cluster); err != nil {
		return
	}

	// Update detections with cluster ID
	for i := range detections {
		if contains(submissionIDs, detections[i].SubmissionID1) && contains(submissionIDs, detections[i].SubmissionID2) {
			detections[i].ClusterID = &cluster.ID
			s.similarityRepo.UpdateDetection(&detections[i])
		}
	}
}

// calculateCodeSimilarity uses normalized Levenshtein distance
func calculateCodeSimilarity(code1, code2 string) float64 {
	// Normalize code (remove whitespace, lowercase)
	norm1 := normalizeCode(code1)
	norm2 := normalizeCode(code2)

	// Calculate Levenshtein distance
	distance := levenshteinDistance(norm1, norm2)

	// Normalize to 0-1 similarity score
	maxLen := float64(max(len(norm1), len(norm2)))
	if maxLen == 0 {
		return 1.0
	}

	return 1.0 - (float64(distance) / maxLen)
}

func normalizeCode(code string) string {
	// Remove all whitespace and convert to lowercase
	var result strings.Builder
	for _, r := range code {
		if !unicode.IsSpace(r) {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

func levenshteinDistance(s1, s2 string) int {
	len1, len2 := len(s1), len(s2)

	// Create matrix
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// Initialize first row and column
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min3(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len1][len2]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func contains(slice []uint, val uint) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
