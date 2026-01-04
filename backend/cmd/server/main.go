package main

import (
	"log"

	"dalivim/internal/config"
	"dalivim/internal/database"
	handler "dalivim/internal/handlers"
	"dalivim/internal/repository"
	"dalivim/internal/router"
	"dalivim/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	// Migrate database
	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)
	telemetryRepo := repository.NewTelemetryRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	activityService := service.NewActivityService(activityRepo, userRepo)
	analysisService := service.NewAnalysisService()
	telemetryService := service.NewTelemetryService(
		telemetryRepo,
		submissionRepo,
		userRepo,
		analysisService,
	)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	activityHandler := handler.NewActivityHandler(activityService)
	telemetryHandler := handler.NewTelemetryHandler(telemetryService)

	// Setup router
	r := router.NewRouter(authHandler, activityHandler, telemetryHandler)
	engine := r.Setup()

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("🚀 Server starting on %s", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
