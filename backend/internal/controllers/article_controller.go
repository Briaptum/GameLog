package controllers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"website-dummy/internal/config"
	"website-dummy/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleController struct{}

// NewArticleController creates a new article controller
func NewArticleController() *ArticleController {
	return &ArticleController{}
}

// GetArticles returns a list of articles
// Public endpoint: only returns public articles
// Admin endpoint: returns all articles with filtering
func (ac *ArticleController) GetArticles(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	var articles []models.Article
	query := db.Model(&models.Article{})

	// Check if this is an admin request (has auth middleware)
	_, isAdmin := c.Get("user")
	
	if !isAdmin {
		// Public endpoint: only return public articles
		query = query.Where("is_public = ?", true)
	} else {
		// Admin endpoint: allow filtering
		if isPublic := c.Query("is_public"); isPublic != "" {
			if isPublic == "true" {
				query = query.Where("is_public = ?", true)
			} else if isPublic == "false" {
				query = query.Where("is_public = ?", false)
			}
		}
	}

	// Category filter
	if category := c.Query("category"); category != "" {
		query = query.Where("categories::text LIKE ?", "%"+category+"%")
	}

	// Limit for public endpoint
	limit := c.Query("limit")
	if limit != "" && !isAdmin {
		if limitInt, err := strconv.Atoi(limit); err == nil {
			query = query.Limit(limitInt)
		}
	}

	// Order by published date descending
	query = query.Order("published_date DESC NULLS LAST, created_at DESC")

	if err := query.Find(&articles).Error; err != nil {
		log.Printf("Error fetching articles: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch articles",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"count":    len(articles),
	})
}

// GetArticle returns a single article by ID or slug
func (ac *ArticleController) GetArticle(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	idOrSlug := c.Param("id")
	var article models.Article

	// Check if this is an admin request
	_, isAdmin := c.Get("user")

	query := db.Preload("Images")
	
	// Try to parse as ID first
	if id, err := strconv.ParseUint(idOrSlug, 10, 32); err == nil {
		query = query.Where("id = ?", id)
	} else {
		// Treat as slug
		query = query.Where("slug = ?", idOrSlug)
	}

	if !isAdmin {
		// Public endpoint: only return if public
		query = query.Where("is_public = ?", true)
	}

	if err := query.First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Article not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"article": article,
	})
}

// CreateArticle handles POST requests to create a new article
func (ac *ArticleController) CreateArticle(c *gin.Context) {
	var req struct {
		Title         string     `json:"title" binding:"required"`
		Slug          *string    `json:"slug"`
		Excerpt       *string    `json:"excerpt"`
		Author        *string    `json:"author"`
		PublishedDate *time.Time `json:"published_date"`
		FeaturedImage *string    `json:"featured_image"`
		IsPublic      bool       `json:"is_public"`
		Categories    models.JSONB `json:"categories"`
		Content       string     `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request. Title and content are required.",
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

	// Generate slug if not provided
	slug := ""
	if req.Slug != nil && *req.Slug != "" {
		slug = *req.Slug
	} else {
		slug = models.GenerateSlug(req.Title)
	}

	// Ensure slug uniqueness
	originalSlug := slug
	counter := 1
	for {
		var existing models.Article
		if err := db.Where("slug = ?", slug).First(&existing).Error; err != nil {
			// Slug is unique, break
			break
		}
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++
	}

	article := models.Article{
		Title:         req.Title,
		Slug:          slug,
		Excerpt:       req.Excerpt,
		Author:        req.Author,
		PublishedDate: req.PublishedDate,
		FeaturedImage: req.FeaturedImage,
		IsPublic:      req.IsPublic,
		Categories:    req.Categories,
		Content:       req.Content,
	}

	if err := db.Create(&article).Error; err != nil {
		log.Printf("Error creating article: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create article",
			"details": err.Error(),
		})
		return
	}

	// Parse markdown content and extract image references
	ac.updateArticleImages(db, article.ID, req.Content)

	// Reload article with images
	db.Preload("Images").First(&article, article.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Article created successfully",
		"article": article,
	})
}

// UpdateArticle handles PUT requests to update an article
func (ac *ArticleController) UpdateArticle(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	id := c.Param("id")
	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Article not found",
		})
		return
	}

	var req struct {
		Title         *string     `json:"title"`
		Slug          *string     `json:"slug"`
		Excerpt       *string     `json:"excerpt"`
		Author        *string     `json:"author"`
		PublishedDate *time.Time  `json:"published_date"`
		FeaturedImage *string     `json:"featured_image"`
		IsPublic      *bool       `json:"is_public"`
		Categories    *models.JSONB `json:"categories"`
		Content       *string     `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Update fields
	if req.Title != nil {
		article.Title = *req.Title
		// Auto-generate slug if title changed and slug not provided
		if req.Slug == nil || *req.Slug == "" {
			article.Slug = models.GenerateSlug(*req.Title)
			// Ensure uniqueness
			originalSlug := article.Slug
			counter := 1
			for {
				var existing models.Article
				if err := db.Where("slug = ? AND id != ?", article.Slug, article.ID).First(&existing).Error; err != nil {
					break
				}
				article.Slug = fmt.Sprintf("%s-%d", originalSlug, counter)
				counter++
			}
		}
	}
	if req.Slug != nil {
		// Check uniqueness if slug is being changed
		if *req.Slug != article.Slug {
			var existing models.Article
			if err := db.Where("slug = ? AND id != ?", *req.Slug, article.ID).First(&existing).Error; err == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Slug already exists",
				})
				return
			}
		}
		article.Slug = *req.Slug
	}
	if req.Excerpt != nil {
		article.Excerpt = req.Excerpt
	}
	if req.Author != nil {
		article.Author = req.Author
	}
	if req.PublishedDate != nil {
		article.PublishedDate = req.PublishedDate
	}
	if req.FeaturedImage != nil {
		article.FeaturedImage = req.FeaturedImage
	}
	if req.IsPublic != nil {
		article.IsPublic = *req.IsPublic
	}
	if req.Categories != nil {
		article.Categories = *req.Categories
	}
	if req.Content != nil {
		article.Content = *req.Content
	}

	if err := db.Save(&article).Error; err != nil {
		log.Printf("Error updating article: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update article",
			"details": err.Error(),
		})
		return
	}

	// Update image references if content changed
	if req.Content != nil {
		ac.updateArticleImages(db, article.ID, *req.Content)
	}

	// Reload article with images
	db.Preload("Images").First(&article, article.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
		"article": article,
	})
}

// DeleteArticle handles DELETE requests to delete an article
func (ac *ArticleController) DeleteArticle(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	id := c.Param("id")
	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Article not found",
		})
		return
	}

	// Delete article (cascade will delete article_images)
	if err := db.Delete(&article, id).Error; err != nil {
		log.Printf("Error deleting article: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete article",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Article deleted successfully",
	})
}

// UploadImage handles image uploads to Bunny CDN
func (ac *ArticleController) UploadImage(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No image file provided",
		})
		return
	}

	// Validate file size (2MB = 2 * 1024 * 1024 bytes)
	const maxSize = 2 * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds 2MB limit. Please use an image optimizer to reduce the file size.",
		})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	allowed := false
	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Allowed: jpg, jpeg, png, gif, webp",
		})
		return
	}

	// Get Bunny CDN configuration
	bunnyStorageZone := os.Getenv("BUNNY_STORAGE_ZONE")
	bunnyAccessKey := os.Getenv("BUNNY_STORAGE_ACCESS_KEY")
	bunnyCDNHostname := os.Getenv("BUNNY_CDN_HOSTNAME")
	bunnyStorageHostname := os.Getenv("BUNNY_STORAGE_HOSTNAME") // Region-specific storage endpoint (e.g., ny.storage.bunnycdn.com)

	if bunnyStorageZone == "" || bunnyAccessKey == "" || bunnyCDNHostname == "" {
		log.Printf("Bunny CDN configuration missing. Storage Zone: %s (has value: %v), Access Key: %s (has value: %v), CDN Hostname: %s (has value: %v)", 
			bunnyStorageZone, bunnyStorageZone != "", 
			"***", bunnyAccessKey != "", 
			bunnyCDNHostname, bunnyCDNHostname != "")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "CDN configuration not available. Please check BUNNY_STORAGE_ZONE, BUNNY_STORAGE_ACCESS_KEY, and BUNNY_CDN_HOSTNAME environment variables.",
		})
		return
	}
	
	// Use region-specific storage hostname if provided, otherwise default to storage.bunnycdn.com
	storageHostname := bunnyStorageHostname
	if storageHostname == "" {
		storageHostname = "storage.bunnycdn.com"
	}
	
	log.Printf("Uploading to Bunny CDN - Storage Zone: %s, Storage Hostname: %s, CDN Hostname: %s", bunnyStorageZone, storageHostname, bunnyCDNHostname)

	// Generate unique filename
	timestamp := time.Now().Format("20060102-150405")
	baseFilename := strings.TrimSuffix(file.Filename, ext)
	filename := fmt.Sprintf("%s-%s%s", timestamp, strings.ReplaceAll(baseFilename, " ", "-"), ext)
	
	// Path in Bunny CDN storage (images/uploads/filename)
	cdnPath := fmt.Sprintf("images/uploads/%s", filename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		log.Printf("Error opening uploaded file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read uploaded file",
		})
		return
	}
	defer src.Close()

	// Read file content into buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		log.Printf("Error reading file content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read file content",
		})
		return
	}

	// Upload to Bunny CDN (use region-specific endpoint if configured)
	cdnURL := fmt.Sprintf("https://%s/%s/%s", storageHostname, bunnyStorageZone, cdnPath)
	log.Printf("Uploading file to: %s (size: %d bytes)", cdnURL, buf.Len())
	req, err := http.NewRequest("PUT", cdnURL, &buf)
	if err != nil {
		log.Printf("Error creating Bunny CDN request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to prepare CDN upload",
		})
		return
	}

	req.Header.Set("AccessKey", bunnyAccessKey)
	
	// Determine content type from file extension
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		// Fallback to MIME type based on extension
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".webp":
			contentType = "image/webp"
		default:
			contentType = "application/octet-stream"
		}
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error uploading to Bunny CDN: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upload to CDN",
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Bunny CDN upload failed with status %d: %s", resp.StatusCode, string(body))
		log.Printf("Request URL: %s", cdnURL)
		log.Printf("Storage Zone: %s", bunnyStorageZone)
		log.Printf("CDN Path: %s", cdnPath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("CDN upload failed: %d - %s", resp.StatusCode, string(body)),
		})
		return
	}

	// Construct the public CDN URL
	// The pull zone must be configured to pull from the storage zone
	// In Bunny CDN dashboard: Pull Zones > Your Pull Zone > Origin > Set to Storage Zone
	publicURL := fmt.Sprintf("https://%s/%s", bunnyCDNHostname, cdnPath)
	log.Printf("File successfully uploaded to Bunny CDN Storage: %s", publicURL)
	log.Printf("Storage path: %s", cdnPath)
	log.Printf("Note: Ensure pull zone '%s' is configured to pull from storage zone '%s'", bunnyCDNHostname, bunnyStorageZone)

	// Return the CDN URL
	c.JSON(http.StatusOK, gin.H{
		"path": publicURL,
		"filename": filename,
	})
}

// GetArticleFiles returns all files linked to an article
func (ac *ArticleController) GetArticleFiles(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	articleID := c.Param("id")
	var images []models.ArticleImage

	if err := db.Where("article_id = ?", articleID).Find(&images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch article files",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": images,
		"count": len(images),
	})
}

// RemoveImage removes an image reference from an article (unlinks but keeps file)
func (ac *ArticleController) RemoveImage(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	articleID := c.Param("id")
	filePath := c.Query("path")
	
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File path is required",
		})
		return
	}

	// Ensure path starts with /
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}

	// Remove from database (only for this article)
	if err := db.Where("article_id = ? AND file_path = ?", articleID, filePath).Delete(&models.ArticleImage{}).Error; err != nil {
		log.Printf("Error removing image reference: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove image reference",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image removed from article successfully",
	})
}

// DeleteImage deletes an image file from Bunny CDN and removes references
func (ac *ArticleController) DeleteImage(c *gin.Context) {
	db := config.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database connection not available",
		})
		return
	}

	// Get path from URL parameter
	// Frontend may send either a CDN URL or a path
	rawPath := c.Param("path")
	
	// Extract path from CDN URL if it's a full URL
	var cdnPath string
	if strings.HasPrefix(rawPath, "http://") || strings.HasPrefix(rawPath, "https://") {
		// Extract path from CDN URL (e.g., https://cdn.example.com/images/uploads/file.jpg)
		// Get everything after the hostname
		parts := strings.Split(rawPath, "/")
		if len(parts) > 3 {
			cdnPath = strings.Join(parts[3:], "/")
		} else {
			cdnPath = rawPath
		}
	} else {
		// Remove leading slash if present
		cdnPath = strings.TrimPrefix(rawPath, "/")
		// Ensure path starts with images/uploads/
		if !strings.HasPrefix(cdnPath, "images/uploads/") {
			cdnPath = "images/uploads/" + strings.TrimPrefix(cdnPath, "/")
		}
	}
	
	// Find the image record by filename and delete from database
	var existingImages []models.ArticleImage
	filename := filepath.Base(cdnPath)
	db.Where("file_path LIKE ?", "%"+filename+"%").Find(&existingImages)
	
	for _, img := range existingImages {
		if err := db.Where("id = ?", img.ID).Delete(&models.ArticleImage{}).Error; err != nil {
			log.Printf("Error deleting image reference: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete image reference",
			})
			return
		}
	}
	
	// Delete from Bunny CDN
	bunnyStorageZone := os.Getenv("BUNNY_STORAGE_ZONE")
	bunnyAccessKey := os.Getenv("BUNNY_STORAGE_ACCESS_KEY")
	bunnyStorageHostname := os.Getenv("BUNNY_STORAGE_HOSTNAME")
	
	// Use region-specific storage hostname if provided, otherwise default
	storageHostname := bunnyStorageHostname
	if storageHostname == "" {
		storageHostname = "storage.bunnycdn.com"
	}

	if bunnyStorageZone != "" && bunnyAccessKey != "" {
		cdnURL := fmt.Sprintf("https://%s/%s/%s", storageHostname, bunnyStorageZone, cdnPath)
		req, err := http.NewRequest("DELETE", cdnURL, nil)
		if err != nil {
			log.Printf("Error creating Bunny CDN delete request: %v", err)
		} else {
			req.Header.Set("AccessKey", bunnyAccessKey)
			
			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("Error deleting from Bunny CDN: %v", err)
			} else {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
					log.Printf("File successfully deleted from Bunny CDN: %s", cdnPath)
				} else {
					log.Printf("Bunny CDN delete returned status %d for: %s", resp.StatusCode, cdnPath)
				}
			}
		}
	} else {
		log.Printf("Bunny CDN configuration missing, skipping CDN deletion")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image deleted successfully",
	})
}

// updateArticleImages parses markdown content and updates article_images table
func (ac *ArticleController) updateArticleImages(db *gorm.DB, articleID uint, content string) {
	// Remove existing image references for this article
	db.Where("article_id = ?", articleID).Delete(&models.ArticleImage{})

	// Regex to find image references in markdown: ![alt](path){align=left/right} or ![alt](path)
	imageRegex := regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)(?:\{align=([^}]+)\})?`)
	matches := imageRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			imagePath := match[2]
			alignment := ""
			if len(match) >= 4 && match[3] != "" {
				alignment = match[3]
			}

			// Only track images in /images/uploads/
			if strings.Contains(imagePath, "/images/uploads/") {
				articleImage := models.ArticleImage{
					ArticleID: articleID,
					FilePath:  imagePath,
				}
				if alignment != "" {
					articleImage.Alignment = &alignment
				}
				db.Create(&articleImage)
			}
		}
	}
}

