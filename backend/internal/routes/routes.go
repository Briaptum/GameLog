package routes

import (
	"website-dummy/internal/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// Initialize controllers
	healthController := controllers.NewHealthController()

	// Public routes
	r.GET("/api/health", healthController.HealthCheck)
}
