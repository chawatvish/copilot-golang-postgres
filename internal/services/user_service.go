package services

import (
	"errors"
	"gin-simple-app/internal/models"
	"gin-simple-app/internal/repository"

	"gorm.io/gorm"
)

// UserService defines the interface for user business logic
type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	CreateUser(req models.CreateUserRequest) (*models.User, error)
	UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id uint) error
	GetUserCount() (int64, error)
}

// UserServiceImpl implements UserService
type UserServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// GetAllUsers returns all users
func (s *UserServiceImpl) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

// GetUserByID returns a user by ID
func (s *UserServiceImpl) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// CreateUser creates a new user
func (s *UserServiceImpl) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	// Check if user with email already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	
	user := &models.User{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   &req.Phone,
		Address: req.Address,
	}
	
	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// UpdateUser updates an existing user
func (s *UserServiceImpl) UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	// Check if email is already taken by another user
	if req.Email != user.Email {
		existingUser, err := s.userRepo.GetByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, errors.New("user with this email already exists")
		}
	}
	
	// Update user fields
	user.Name = req.Name
	user.Email = req.Email
	user.Phone = &req.Phone
	user.Address = req.Address
	
	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// DeleteUser deletes a user by ID
func (s *UserServiceImpl) DeleteUser(id uint) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	
	return s.userRepo.Delete(id)
}

// GetUserCount returns the total number of users
func (s *UserServiceImpl) GetUserCount() (int64, error) {
	return s.userRepo.Count()
}
