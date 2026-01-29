package models

import (
	"time"
)

// Testimonial represents a customer testimonial
type Testimonial struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Content          string    `json:"content" gorm:"type:text;not null"`
	FromUser         string    `json:"from_user" gorm:"type:text;not null"`
	ForUser          *string   `json:"for_user" gorm:"type:text"`
	ShowOnFrontPage  bool      `json:"show_on_front_page" gorm:"not null;default:false"`
	Order            int       `json:"order" gorm:"not null;default:0"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName specifies the table name for the Testimonial model
func (Testimonial) TableName() string {
	return "testimonials"
}
