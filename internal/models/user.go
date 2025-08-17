package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID                uint           `json:"id" gorm:"primarykey"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
	Name              string         `json:"name" gorm:"not null" binding:"required"`
	Email             string         `json:"email" gorm:"uniqueIndex;not null" binding:"required,email"`
	Password          string         `json:"-" gorm:"not null"` // Hidden from JSON responses
	Phone             *string        `json:"phone" gorm:"type:text;default:null" binding:"required"`
	Address           *string        `json:"address,omitempty" gorm:"type:text"`
	IsActive          bool           `json:"is_active" gorm:"default:true"`
	IsEmailVerified   bool           `json:"is_email_verified" gorm:"default:false"`
	EmailVerificationToken *string   `json:"-" gorm:"type:text"` // Hidden from JSON
	PasswordResetToken     *string   `json:"-" gorm:"type:text"` // Hidden from JSON
	PasswordResetExpiry    *time.Time `json:"-"`                  // Hidden from JSON
	LastLoginAt           *time.Time `json:"last_login_at"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name    string  `json:"name" binding:"required"`
	Email   string  `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
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

// Authentication request models
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name            string  `json:"name" binding:"required"`
	Email           string  `json:"email" binding:"required,email"`
	Password        string  `json:"password" binding:"required,min=6"`
	ConfirmPassword string  `json:"confirm_password" binding:"required"`
	Phone           string  `json:"phone" binding:"required"`
	Address         *string `json:"address,omitempty"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// Response models
type LoginResponse struct {
	User        *User  `json:"user"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type UserResponse struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Phone           *string    `json:"phone"`
	Address         *string    `json:"address"`
	IsActive        bool       `json:"is_active"`
	IsEmailVerified bool       `json:"is_email_verified"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastLoginAt     *time.Time `json:"last_login_at"`
}
