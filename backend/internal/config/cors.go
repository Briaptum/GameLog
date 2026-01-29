package config

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS configures CORS middleware for the Gin router
func SetupCORS(r *gin.Engine) {
	config := cors.DefaultConfig()
	
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:5173",
		"http://frontend:3000", // Docker internal network
	}
	
	// Add SITE_URL to allowed origins if set (matches rolston pattern)
	if siteURL := os.Getenv("SITE_URL"); siteURL != "" {
		allowedOrigins = append(allowedOrigins, siteURL)
	}
	
	config.AllowOrigins = allowedOrigins
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}
	// Expose custom headers so frontend can read them
	config.ExposeHeaders = []string{
		"X-Cache-Generated-At",
	}
	config.AllowCredentials = true
	
	r.Use(cors.New(config))
}
