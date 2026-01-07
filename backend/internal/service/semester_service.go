package service

import (
	"dalivim/internal/models"
	"dalivim/internal/repository"
	"time"
)

type SemesterService interface {
	CreateSemester(year, period int, startDate, endDate time.Time) (*models.Semester, error)
	GetActiveSemester() (*models.Semester, error)
	GetAllSemesters() ([]models.Semester, error)
	UpdateAllStudentSemesters() error
}

type semesterService struct {
	semesterRepo repository.SemesterRepository
	userRepo     repository.UserRepository
}

func NewSemesterService(
	semesterRepo repository.SemesterRepository,
	userRepo repository.UserRepository,
) SemesterService {
	return &semesterService{
		semesterRepo: semesterRepo,
		userRepo:     userRepo,
	}
}

func (s *semesterService) CreateSemester(year, period int, startDate, endDate time.Time) (*models.Semester, error) {
	semester := &models.Semester{
		Year:      year,
		Period:    period,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.semesterRepo.Create(semester); err != nil {
		return nil, err
	}

	return semester, nil
}

func (s *semesterService) GetActiveSemester() (*models.Semester, error) {
	return s.semesterRepo.FindActive()
}

func (s *semesterService) GetAllSemesters() ([]models.Semester, error) {
	return s.semesterRepo.FindAll()
}

// UpdateAllStudentSemesters recalculates the current semester for all students
func (s *semesterService) UpdateAllStudentSemesters() error {
	// Get active semester
	activeSemester, err := s.GetActiveSemester()
	if err != nil {
		return err
	}

	// Get all students
	students, err := s.userRepo.FindAllStudents()
	if err != nil {
		return err
	}

	// Update each student's current semester
	for _, student := range students {
		student.UpdateCurrentSemester(activeSemester.Year, activeSemester.Period)
		s.userRepo.Update(&student)
	}

	return nil
}
