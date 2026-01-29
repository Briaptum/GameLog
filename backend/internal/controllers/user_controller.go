package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

// NewUserController creates a new user controller
func NewUserController() *UserController {
	return &UserController{}
}

// GetName returns a hardcoded name for the frontend
func (uc *UserController) GetName(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name": "Tyler",
	})
}
