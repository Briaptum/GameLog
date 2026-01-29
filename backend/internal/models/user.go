package models

import (
	"time"
)

// User represents a user in the administration system
type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name" gorm:"not null"`
	Email           string     `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash    string     `json:"-" gorm:"column:password_hash;not null"`
	Role            string     `json:"role" gorm:"not null;default:user"`
	IsActive        bool       `json:"is_active" gorm:"not null;default:true"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// UserRole constants for role-based access control
const (
	RoleAdmin     = "admin"
	RoleModerator = "moderator"
	RoleUser      = "user"
)

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsModerator checks if the user has moderator role or higher
func (u *User) IsModerator() bool {
	return u.Role == RoleAdmin || u.Role == RoleModerator
}
