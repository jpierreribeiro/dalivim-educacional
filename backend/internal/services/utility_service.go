package services

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func generateToken(userID uint) string {
	// Simple token generation (use JWT in production)
	return hex.EncodeToString([]byte(time.Now().String()))
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

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple auth middleware (implement proper JWT validation in production)
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Mock user ID (in production, extract from JWT)
		c.Set("userID", uint(1))
		c.Next()
	}
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0.0
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return 0
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
