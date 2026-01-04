package service

import (
	"encoding/json"

	"dalivim/internal/models"
	"dalivim/internal/repository"
)

type TelemetryService interface {
	ProcessTelemetry(activityID, studentID uint, timestamp int64, isFinal bool, code string, features, rawEvents map[string]interface{}) (AnalysisResult, error)
	GetSubmissions(activityID uint) ([]models.Submission, error)
}

type telemetryService struct {
	telemetryRepo   repository.TelemetryRepository
	submissionRepo  repository.SubmissionRepository
	userRepo        repository.UserRepository
	analysisService AnalysisService
}

func NewTelemetryService(
	telemetryRepo repository.TelemetryRepository,
	submissionRepo repository.SubmissionRepository,
	userRepo repository.UserRepository,
	analysisService AnalysisService,
) TelemetryService {
	return &telemetryService{
		telemetryRepo:   telemetryRepo,
		submissionRepo:  submissionRepo,
		userRepo:        userRepo,
		analysisService: analysisService,
	}
}

func (s *telemetryService) ProcessTelemetry(
	activityID, studentID uint,
	timestamp int64,
	isFinal bool,
	code string,
	features, rawEvents map[string]interface{},
) (AnalysisResult, error) {
	// Analyze behavior
	analysis := s.analysisService.Analyze(features)

	// Save telemetry data
	featuresJSON, _ := json.Marshal(features)
	eventsJSON, _ := json.Marshal(rawEvents)

	telemetry := &models.TelemetryData{
		ActivityID: activityID,
		StudentID:  studentID,
		Timestamp:  timestamp,
		IsFinal:    isFinal,
		Features:   string(featuresJSON),
		RawEvents:  string(eventsJSON),
	}

	if err := s.telemetryRepo.Create(telemetry); err != nil {
		return analysis, err
	}

	// If final submission, create submission record
	if isFinal {
		student, _ := s.userRepo.FindByID(studentID)

		pasteEventsJSON, _ := json.Marshal(rawEvents["pasteEvents"])

		submission := &models.Submission{
			ActivityID:           activityID,
			StudentID:            studentID,
			StudentName:          student.Name,
			StudentEmail:         student.Email,
			Code:                 code,
			AuthorshipScore:      analysis.AuthorshipScore,
			Confidence:           analysis.Confidence,
			SignalsArray:         analysis.Signals,
			AvgKeystrokeInterval: getFloat(features, "avgKeystrokeInterval"),
			StdKeystrokeInterval: getFloat(features, "stdKeystrokeInterval"),
			PasteEvents:          getInt(features, "pasteEvents"),
			PasteCharRatio:       getFloat(features, "pasteCharRatio"),
			DeleteRatio:          getFloat(features, "deleteRatio"),
			FocusLossCount:       getInt(features, "focusLossCount"),
			LinearEditingScore:   getFloat(features, "linearEditingScore"),
			Burstiness:           getFloat(features, "burstiness"),
			TimeToFirstRun:       getFloat(features, "timeToFirstRun"),
			ExecutionCount:       getInt(features, "executionCount"),
			TotalTime:            getFloat(features, "totalTime"),
			KeystrokeCount:       getInt(features, "totalKeystrokes"),
			PasteEventDetails:    string(pasteEventsJSON),
		}

		s.submissionRepo.Create(submission)
	}

	return analysis, nil
}

func (s *telemetryService) GetSubmissions(activityID uint) ([]models.Submission, error) {
	return s.submissionRepo.FindByActivityID(activityID)
}
