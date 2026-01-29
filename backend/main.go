package main

import (
	"log"
	"os"

	"website-dummy/internal/config"
	"website-dummy/internal/middleware"
	"website-dummy/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.InitDatabase()

	// Initialize Gin router
	r := gin.New()

	// Add middleware
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// Configure CORS
	config.SetupCORS(r)

	// Load Auth0 configuration into Gin context
	r.Use(func(c *gin.Context) {
		c.Set("AUTH0_DOMAIN", os.Getenv("AUTH0_DOMAIN"))
		c.Set("AUTH0_CLIENT_ID", os.Getenv("AUTH0_CLIENT_ID"))
		c.Set("AUTH0_AUDIENCE", os.Getenv("AUTH0_AUDIENCE"))
		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(r)

	// Start server (always on port 8080 inside container)
	port := "8080"

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📚 API Documentation available at http://localhost:%s/api/health", port)
	
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
