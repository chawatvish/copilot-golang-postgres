package repository

import (
	"gin-simple-app/internal/models"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByPasswordResetToken(token string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	Count() (int64, error)
}

// GormUserRepository implements UserRepository using GORM
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GORM user repository
func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{
		db: db,
	}
}

// GetAll returns all users
func (r *GormUserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// GetByID returns a user by ID
func (r *GormUserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail returns a user by email
func (r *GormUserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPasswordResetToken returns a user by password reset token
func (r *GormUserRepository) GetByPasswordResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("password_reset_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create creates a new user
func (r *GormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// Update updates an existing user
func (r *GormUserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user by ID (soft delete)
func (r *GormUserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// Count returns the total number of users
func (r *GormUserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}
