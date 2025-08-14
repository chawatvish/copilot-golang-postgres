package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `json:"name" gorm:"not null" binding:"required"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" binding:"required,email"`
	Phone     *string        `json:"phone" gorm:"type:text;default:null" binding:"required"`
	Address   *string        `json:"address,omitempty" gorm:"type:text"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name    string  `json:"name" binding:"required"`
	Email   string  `json:"email" binding:"required,email"`
	Phone   string  `json:"phone" binding:"required"`
	Address *string `json:"address,omitempty"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name    string  `json:"name" binding:"required"`
	Email   string  `json:"email" binding:"required,email"`
	Phone   string  `json:"phone" binding:"required"`
	Address *string `json:"address,omitempty"`
}
