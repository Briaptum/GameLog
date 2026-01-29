package controllers

import (
	"log"
	"net/http"
	"strconv"

	"website-dummy/internal/config"
	"website-dummy/internal/models"

	"github.com/gin-gonic/gin"
)

type TestimonialController struct{}

// NewTestimonialController creates a new testimonial controller
func NewTestimonialController() *TestimonialController {
	return &TestimonialController{}
}

// GetTestimonials returns a list of testimonials
// Public endpoint with optional filtering
func (tc *TestimonialController) GetTestimonials(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var testimonials []models.Testimonial
	query := db.Model(&models.Testimonial{})

	// Filter by for_user if provided
	if forUser := c.Query("for_user"); forUser != "" {
		query = query.Where("for_user = ?", forUser)
	}

	// Filter by show_on_front_page if provided
	if showOnFrontPage := c.Query("show_on_front_page"); showOnFrontPage != "" {
		if showOnFrontPage == "true" {
			query = query.Where("show_on_front_page = ?", true)
		} else if showOnFrontPage == "false" {
			query = query.Where("show_on_front_page = ?", false)
		}
	}

	// Order by order field (ascending), then created_at (descending)
	query = query.Order("\"order\" ASC, created_at DESC")

	if err := query.Find(&testimonials).Error; err != nil {
		log.Printf("Error fetching testimonials: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch testimonials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"testimonials": testimonials,
		"count":        len(testimonials),
	})
}

// GetTestimonialsAdmin returns all testimonials for admin
func (tc *TestimonialController) GetTestimonialsAdmin(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var testimonials []models.Testimonial
	query := db.Model(&models.Testimonial{})

	// Pagination
	page := 1
	pageSize := 50
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	// Order by order field (ascending), then created_at (descending)
	query = query.Order("\"order\" ASC, created_at DESC")

	if err := query.Find(&testimonials).Error; err != nil {
		log.Printf("Error fetching testimonials: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch testimonials",
		})
		return
	}

	// Get total count
	var total int64
	db.Model(&models.Testimonial{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"testimonials": testimonials,
		"count":        len(testimonials),
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	})
}

// GetTestimonial returns a single testimonial by ID
func (tc *TestimonialController) GetTestimonial(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	id := c.Param("id")
	var testimonial models.Testimonial

	if err := db.First(&testimonial, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Testimonial not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"testimonial": testimonial,
	})
}

// CreateTestimonial handles POST requests to create a new testimonial
func (tc *TestimonialController) CreateTestimonial(c *gin.Context) {
	var req struct {
		Content         string  `json:"content" binding:"required"`
		FromUser        string  `json:"from_user" binding:"required"`
		ForUser         *string `json:"for_user"`
		ShowOnFrontPage bool    `json:"show_on_front_page"`
		Order           int     `json:"order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request. Content and from_user are required.",
		})
		return
	}

	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	testimonial := models.Testimonial{
		Content:         req.Content,
		FromUser:        req.FromUser,
		ForUser:         req.ForUser,
		ShowOnFrontPage: req.ShowOnFrontPage,
		Order:           req.Order,
	}

	if err := db.Create(&testimonial).Error; err != nil {
		log.Printf("Error creating testimonial: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create testimonial",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Testimonial created successfully",
		"testimonial": testimonial,
	})
}

// UpdateTestimonial handles PUT requests to update a testimonial
func (tc *TestimonialController) UpdateTestimonial(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	id := c.Param("id")
	var testimonial models.Testimonial
	if err := db.First(&testimonial, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Testimonial not found",
		})
		return
	}

	var req struct {
		Content         *string `json:"content"`
		FromUser        *string `json:"from_user"`
		ForUser         *string `json:"for_user"`
		ShowOnFrontPage *bool   `json:"show_on_front_page"`
		Order           *int    `json:"order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Update fields
	if req.Content != nil {
		testimonial.Content = *req.Content
	}
	if req.FromUser != nil {
		testimonial.FromUser = *req.FromUser
	}
	if req.ForUser != nil {
		// Allow setting to NULL by sending empty string
		if *req.ForUser == "" {
			testimonial.ForUser = nil
		} else {
			testimonial.ForUser = req.ForUser
		}
	}
	if req.ShowOnFrontPage != nil {
		testimonial.ShowOnFrontPage = *req.ShowOnFrontPage
	}
	if req.Order != nil {
		testimonial.Order = *req.Order
	}

	if err := db.Save(&testimonial).Error; err != nil {
		log.Printf("Error updating testimonial: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update testimonial",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Testimonial updated successfully",
		"testimonial": testimonial,
	})
}

// DeleteTestimonial handles DELETE requests to delete a testimonial
func (tc *TestimonialController) DeleteTestimonial(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	id := c.Param("id")
	var testimonial models.Testimonial
	if err := db.First(&testimonial, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Testimonial not found",
		})
		return
	}

	if err := db.Delete(&testimonial, id).Error; err != nil {
		log.Printf("Error deleting testimonial: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete testimonial",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Testimonial deleted successfully",
	})
}
