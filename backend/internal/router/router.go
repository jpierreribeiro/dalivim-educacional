package router

import (
	"time"

	handler "dalivim/internal/handlers"
	"dalivim/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler      *handler.AuthHandler
	activityHandler  *handler.ActivityHandler
	telemetryHandler *handler.TelemetryHandler
}

func NewRouter(
	authHandler *handler.AuthHandler,
	activityHandler *handler.ActivityHandler,
	telemetryHandler *handler.TelemetryHandler,
) *Router {
	return &Router{
		authHandler:      authHandler,
		activityHandler:  activityHandler,
		telemetryHandler: telemetryHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	api := router.Group("/api")
	{
		// Auth
		api.POST("/auth/register", r.authHandler.Register)
		api.POST("/auth/login", r.authHandler.Login)

		// Activity (public)
		api.POST("/activities/join/:inviteToken", r.activityHandler.Join)

		// Telemetry (public)
		api.POST("/telemetry", r.telemetryHandler.Process)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Activities
		protected.POST("/activities", r.activityHandler.Create)
		protected.GET("/activities", r.activityHandler.GetAll)
		protected.GET("/activities/:id", r.activityHandler.GetByID)

		// Submissions
		protected.GET("/activities/:id/submissions", r.telemetryHandler.GetSubmissions)
	}

	return router
}
