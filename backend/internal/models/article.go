package models

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Article represents a blog article
type Article struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"type:text;not null" json:"title"`
	Slug          string         `gorm:"type:text;not null;uniqueIndex" json:"slug"`
	Excerpt       *string        `gorm:"type:text" json:"excerpt"`
	Author        *string        `gorm:"type:text" json:"author"`
	PublishedDate *time.Time     `gorm:"type:timestamptz" json:"published_date"`
	FeaturedImage *string        `gorm:"type:text" json:"featured_image"`
	IsPublic      bool           `gorm:"not null;default:false" json:"is_public"`
	Categories    JSONB          `gorm:"type:jsonb;default:'[]'::jsonb" json:"categories"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Images        []ArticleImage `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"images,omitempty"`
}

// TableName specifies the table name for the Article model
func (Article) TableName() string {
	return "articles"
}

// GenerateSlug generates a URL-friendly slug from a title
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.TrimSpace(slug)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove special characters, keep only alphanumeric and hyphens
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	slug = result.String()
	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	return slug
}

// ArticleImage represents a file/image linked to an article
type ArticleImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"not null;index" json:"article_id"`
	FilePath  string    `gorm:"type:text;not null;index" json:"file_path"`
	Alignment *string   `gorm:"type:text" json:"alignment"` // 'left', 'right', or null
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies the table name for the ArticleImage model
func (ArticleImage) TableName() string {
	return "article_images"
}

