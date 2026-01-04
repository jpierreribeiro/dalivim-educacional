package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Public routes
	r.POST("/api/auth/register", register)
	r.POST("/api/auth/login", login)
	r.POST("/api/activities/join/:inviteToken", joinActivity)

	// Student routes
	r.POST("/api/telemetry", receiveTelemetry)

	// Protected routes (require auth)
	authorized := r.Group("/api")
	authorized.Use(authMiddleware())
	{
		// Professor routes
		authorized.POST("/activities", createActivity)
		authorized.GET("/activities", getActivities)
		authorized.GET("/activities/:id", getActivity)
		authorized.GET("/activities/:id/submissions", getSubmissions)
	}
}
