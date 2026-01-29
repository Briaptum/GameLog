package models

import (
	"time"
)

// Rental represents a rental property
type Rental struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Order           int       `gorm:"column:\"order\";not null" json:"order"`
	UnparsedAddress string    `gorm:"type:text;not null" json:"unparsed_address"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName specifies the table name for the Rental model
func (Rental) TableName() string {
	return "rentals"
}

