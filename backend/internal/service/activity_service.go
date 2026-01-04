package service

import (
	"crypto/rand"
	"encoding/hex"

	"dalivim/internal/models"
	"dalivim/internal/repository"
)

type ActivityService interface {
	Create(professorID uint, title, description, language string, timeLimit int) (*models.Activity, error)
	GetByID(id uint) (*models.Activity, error)
	GetByProfessorID(professorID uint) ([]ActivityWithCount, error)
	JoinActivity(inviteToken string) (*models.Activity, *models.User, error)
}

type ActivityWithCount struct {
	models.Activity
	SubmissionCount int64 `json:"submissionCount"`
}

type activityService struct {
	activityRepo repository.ActivityRepository
	userRepo     repository.UserRepository
}

func NewActivityService(activityRepo repository.ActivityRepository, userRepo repository.UserRepository) ActivityService {
	return &activityService{
		activityRepo: activityRepo,
		userRepo:     userRepo,
	}
}

func (s *activityService) Create(professorID uint, title, description, language string, timeLimit int) (*models.Activity, error) {
	activity := &models.Activity{
		ProfessorID: professorID,
		Title:       title,
		Description: description,
		Language:    language,
		TimeLimit:   timeLimit,
		InviteToken: generateInviteToken(),
	}

	if err := s.activityRepo.Create(activity); err != nil {
		return nil, err
	}

	return activity, nil
}

func (s *activityService) GetByID(id uint) (*models.Activity, error) {
	return s.activityRepo.FindByID(id)
}

func (s *activityService) GetByProfessorID(professorID uint) ([]ActivityWithCount, error) {
	activities, err := s.activityRepo.FindByProfessorID(professorID)
	if err != nil {
		return nil, err
	}

	result := make([]ActivityWithCount, len(activities))
	for i, activity := range activities {
		count, _ := s.activityRepo.CountSubmissions(activity.ID)
		result[i] = ActivityWithCount{
			Activity:        activity,
			SubmissionCount: count,
		}
	}

	return result, nil
}

func (s *activityService) JoinActivity(inviteToken string) (*models.Activity, *models.User, error) {
	activity, err := s.activityRepo.FindByInviteToken(inviteToken)
	if err != nil {
		return nil, nil, err
	}

	// Create anonymous student
	student := &models.User{
		Email: generateAnonymousEmail(),
		Name:  "Anonymous Student",
		Role:  "student",
	}

	if err := s.userRepo.Create(student); err != nil {
		return nil, nil, err
	}

	return activity, student, nil
}

func generateInviteToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func generateAnonymousEmail() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "student_" + hex.EncodeToString(bytes) + "@anonymous.local"
}
