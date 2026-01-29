package controllers

import (
	"net/http"

	"website-dummy/internal/config"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck returns the health status of the API
func (hc *HealthController) HealthCheck(c *gin.Context) {
	response := gin.H{
		"status": "ok",
		"service": "manage-api",
	}

	// Check database connection
	db := config.GetDB()
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Ping(); err == nil {
				response["database"] = "connected"
			} else {
				response["database"] = "disconnected"
			}
		} else {
			response["database"] = "error"
		}
	} else {
		response["database"] = "not_initialized"
	}

	c.JSON(http.StatusOK, response)
}
