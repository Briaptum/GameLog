package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"website-dummy/internal/config"
	"website-dummy/internal/models"
	"website-dummy/internal/services"

	"github.com/gin-gonic/gin"
)

type ContactRequestController struct {
	emailService *services.EmailService
}

// NewContactRequestController creates a new contact request controller
func NewContactRequestController() *ContactRequestController {
	return &ContactRequestController{
		emailService: services.NewEmailService(),
	}
}

// CreateContactRequest handles POST requests to create a new contact request
func (crc *ContactRequestController) CreateContactRequest(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Email   string `json:"email" binding:"required,email"`
		Phone   string `json:"phone"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request. Name, email, and message are required.",
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

	// Get IP address and user agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Rate limiting: Check if this IP has submitted more than 3 requests in the last hour
	var recentCount int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	if err := db.Model(&models.ContactRequest{}).
		Where("ip_address = ? AND created_at > ?", ipAddress, oneHourAgo).
		Count(&recentCount).Error; err != nil {
		log.Printf("Error checking rate limit: %v", err)
		// Continue anyway - don't block legitimate users if there's a DB error
	} else if recentCount >= 3 {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many requests. Please wait before submitting again.",
		})
		return
	}

	// Create request
	contactRequest := models.ContactRequest{
		Name:      req.Name,
		Email:     req.Email,
		Message:   req.Message,
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
		Metadata:  models.JSONB{},
	}

	// Add phone if provided
	if req.Phone != "" {
		contactRequest.Phone = &req.Phone
	}

	// Save to database
	if err := db.Create(&contactRequest).Error; err != nil {
		fmt.Printf("Database error saving contact request: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save contact request",
			"details": err.Error(),
		})
		return
	}

	// Send notification email (non-blocking, log errors but don't fail the request)
	go func() {
		// Recover from any panics in email sending
		defer func() {
			if r := recover(); r != nil {
				log.Printf("❌ PANIC in email sending goroutine for contact request ID %d: %v", contactRequest.ID, r)
			}
		}()
		
		log.Printf("📬 Starting email send process for contact request ID %d", contactRequest.ID)
		if err := crc.sendNotificationEmail(contactRequest); err != nil {
			log.Printf("❌ Failed to send notification email for contact request ID %d: %v", contactRequest.ID, err)
		} else {
			log.Printf("✅ Successfully sent notification email for contact request ID %d from %s", contactRequest.ID, contactRequest.Email)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Contact request received successfully",
		"id":      contactRequest.ID,
	})
}

// sendNotificationEmail sends an email notification about the new contact request
func (crc *ContactRequestController) sendNotificationEmail(contactRequest models.ContactRequest) error {
	if !crc.emailService.IsConfigured() {
		return nil // SMTP not configured, silently skip
	}

	subject := fmt.Sprintf("New Contact Form Submission from %s", contactRequest.Name)

	// Build HTML email body with proper structure
	messageHTML := strings.ReplaceAll(contactRequest.Message, "&", "&amp;")
	messageHTML = strings.ReplaceAll(messageHTML, "<", "&lt;")
	messageHTML = strings.ReplaceAll(messageHTML, ">", "&gt;")
	messageHTML = strings.ReplaceAll(messageHTML, "\n", "<br>")

	phoneHTML := ""
	if contactRequest.Phone != nil && *contactRequest.Phone != "" {
		phoneHTML = fmt.Sprintf(`<p style="margin: 10px 0;"><strong>Phone:</strong> <a href="tel:%s" style="color: #0066cc; text-decoration: none;">%s</a></p>`, *contactRequest.Phone, *contactRequest.Phone)
	}

	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 20px; background-color: #f9f9f9;">
	<div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; padding: 30px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
		<h2 style="color: #0066cc; margin-top: 0; margin-bottom: 20px; font-size: 24px;">New Contact Form Submission</h2>
		
		<div style="background-color: #f5f5f5; padding: 20px; border-radius: 5px; margin: 20px 0;">
			<p style="margin: 10px 0;"><strong>Name:</strong> %s</p>
			<p style="margin: 10px 0;"><strong>Email:</strong> <a href="mailto:%s" style="color: #0066cc; text-decoration: none;">%s</a></p>
			%s
			<p style="margin: 10px 0;"><strong>Message:</strong></p>
			<div style="background-color: white; padding: 15px; border-left: 3px solid #0066cc; margin: 10px 0; border-radius: 3px;">
				%s
			</div>
		</div>
		
		<div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #ddd; font-size: 12px; color: #666;">
			<p style="margin: 5px 0;"><strong>Submitted:</strong> %s</p>
			%s
		</div>
	</div>
</body>
</html>`,
		contactRequest.Name,
		contactRequest.Email,
		contactRequest.Email,
		phoneHTML,
		messageHTML,
		contactRequest.CreatedAt.Format("January 2, 2006 at 3:04 PM MST"),
		crc.formatMetadata(contactRequest),
	)

	return crc.emailService.SendNotificationEmail(subject, body)
}

// formatMetadata formats the metadata for display in the email
func (crc *ContactRequestController) formatMetadata(contactRequest models.ContactRequest) string {
	var parts []string
	if contactRequest.IPAddress != nil && *contactRequest.IPAddress != "" {
		parts = append(parts, fmt.Sprintf("<p><strong>IP Address:</strong> %s</p>", *contactRequest.IPAddress))
	}
	if contactRequest.UserAgent != nil && *contactRequest.UserAgent != "" {
		parts = append(parts, fmt.Sprintf("<p><strong>User Agent:</strong> %s</p>", *contactRequest.UserAgent))
	}
	if len(parts) > 0 {
		return strings.Join(parts, "")
	}
	return ""
}

// GetContactRequests returns all contact requests
func (crc *ContactRequestController) GetContactRequests(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var requests []models.ContactRequest

	// Order by created_at descending (newest first)
	if err := db.Order("created_at DESC").Find(&requests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch contact requests",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"requests": requests,
	})
}

// GetContactRequest returns a single contact request by ID
func (crc *ContactRequestController) GetContactRequest(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var request models.ContactRequest
	id := c.Param("id")

	if err := db.First(&request, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Contact request not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request": request,
	})
}

// DeleteContactRequest soft deletes a contact request by ID
func (crc *ContactRequestController) DeleteContactRequest(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	id := c.Param("id")
	var request models.ContactRequest

	// First check if the request exists
	if err := db.First(&request, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Contact request not found",
		})
		return
	}

	// Soft delete the request
	if err := db.Delete(&request, id).Error; err != nil {
		log.Printf("Error deleting contact request ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete contact request",
		})
		return
	}

	log.Printf("✅ Soft deleted contact request ID %s", id)
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact request deleted successfully",
		"id":      id,
	})
}

