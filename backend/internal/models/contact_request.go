package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// ContactRequest represents a contact form submission
type ContactRequest struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"not null"`
	Phone     *string        `json:"phone" gorm:"type:varchar(20)"`
	Message   string         `json:"message" gorm:"type:text;not null"`
	IPAddress *string        `json:"ip_address" gorm:"column:ip_address;type:varchar(45)"`
	UserAgent *string        `json:"user_agent" gorm:"column:user_agent;type:text"`
	Metadata  JSONB          `json:"metadata" gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName specifies the table name for the ContactRequest model
func (ContactRequest) TableName() string {
	return "contact_requests"
}

// JSONB is a custom type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = JSONB{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

