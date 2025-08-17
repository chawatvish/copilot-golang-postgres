package services

import (
	"errors"
	"gin-simple-app/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Helper function to create string pointers
func stringPointer(s string) *string {
	return &s
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByPasswordResetToken(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestNewUserService(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)
	
	assert.NotNil(t, service)
	assert.IsType(t, &UserServiceImpl{}, service)
}

func TestUserService_GetAllUsers(t *testing.T) {
	tests := []struct {
		name           string
		repoResponse   []models.User
		repoError      error
		expectedUsers  []models.User
		expectedError  error
	}{
		{
			name: "success",
			repoResponse: []models.User{
				{ID: 1, Name: "John", Email: "john@example.com"},
				{ID: 2, Name: "Jane", Email: "jane@example.com"},
			},
			repoError: nil,
			expectedUsers: []models.User{
				{ID: 1, Name: "John", Email: "john@example.com"},
				{ID: 2, Name: "Jane", Email: "jane@example.com"},
			},
			expectedError: nil,
		},
		{
			name:          "repository error",
			repoResponse:  []models.User{},
			repoError:     errors.New("database error"),
			expectedUsers: []models.User{},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			mockRepo.On("GetAll").Return(tt.repoResponse, tt.repoError)
			
			service := NewUserService(mockRepo)
			users, err := service.GetAllUsers()
			
			assert.Equal(t, tt.expectedUsers, users)
			assert.Equal(t, tt.expectedError, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		repoResponse  *models.User
		repoError     error
		expectedUser  *models.User
		expectedError error
	}{
		{
			name:   "success",
			userID: 1,
			repoResponse: &models.User{
				ID: 1, Name: "John", Email: "john@example.com",
			},
			repoError:    nil,
			expectedUser: &models.User{ID: 1, Name: "John", Email: "john@example.com"},
			expectedError: nil,
		},
		{
			name:          "user not found",
			userID:        999,
			repoResponse:  nil,
			repoError:     gorm.ErrRecordNotFound,
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
		{
			name:          "repository error",
			userID:        1,
			repoResponse:  nil,
			repoError:     errors.New("database error"),
			expectedUser:  nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			mockRepo.On("GetByID", tt.userID).Return(tt.repoResponse, tt.repoError)
			
			service := NewUserService(mockRepo)
			user, err := service.GetUserByID(tt.userID)
			
			assert.Equal(t, tt.expectedUser, user)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		request        models.CreateUserRequest
		emailExistsErr error
		createErr      error
		expectedUser   *models.User
		expectedError  string
	}{
		{
			name: "success",
			request: models.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    "+1234567890",
				Address:  stringPointer("123 Main St"),
			},
			emailExistsErr: gorm.ErrRecordNotFound, // Email doesn't exist
			createErr:      nil,
			expectedUser: &models.User{
				Name:    "John Doe",
				Email:   "john@example.com",
				Phone:   stringPointer("+1234567890"),
				Address: stringPointer("123 Main St"),
			},
			expectedError: "",
		},
		{
			name: "duplicate email",
			request: models.CreateUserRequest{
				Name:     "John Doe",
				Email:    "existing@example.com",
				Password: "password123",
			},
			emailExistsErr: nil, // Email exists
			createErr:      nil,
			expectedUser:   nil,
			expectedError:  "user with this email already exists",
		},
		{
			name: "repository create error",
			request: models.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			emailExistsErr: gorm.ErrRecordNotFound,
			createErr:      errors.New("database error"),
			expectedUser:   nil,
			expectedError:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			
			// Mock the email existence check
			if tt.emailExistsErr == gorm.ErrRecordNotFound {
				mockRepo.On("GetByEmail", tt.request.Email).Return(nil, tt.emailExistsErr)
			} else {
				existingUser := &models.User{Email: tt.request.Email}
				mockRepo.On("GetByEmail", tt.request.Email).Return(existingUser, tt.emailExistsErr)
			}
			
			// Mock the create operation only if email doesn't exist
			if tt.emailExistsErr == gorm.ErrRecordNotFound {
				mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(tt.createErr)
			}
			
			service := NewUserService(mockRepo)
			user, err := service.CreateUser(tt.request)
			
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Phone, user.Phone)
				assert.Equal(t, tt.expectedUser.Address, user.Address)
				// Password should be hashed, not plain text
				assert.NotEqual(t, tt.request.Password, user.Password)
				assert.NotEmpty(t, user.Password)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		request       models.UpdateUserRequest
		existingUser  *models.User
		getUserErr    error
		updateErr     error
		expectedError string
	}{
		{
			name:   "success",
			userID: 1,
			request: models.UpdateUserRequest{
				Name:    "Updated Name",
				Email:   "updated@example.com",
				Phone:   "+9876543210",
				Address: stringPointer("456 Oak St"),
			},
			existingUser: &models.User{
				ID: 1, Name: "Old Name", Email: "old@example.com",
			},
			getUserErr:    nil,
			updateErr:     nil,
			expectedError: "",
		},
		{
			name:   "user not found",
			userID: 999,
			request: models.UpdateUserRequest{
				Name: "Updated Name",
			},
			existingUser:  nil,
			getUserErr:    gorm.ErrRecordNotFound,
			updateErr:     nil,
			expectedError: "user not found",
		},
		{
			name:   "update error",
			userID: 1,
			request: models.UpdateUserRequest{
				Name: "Updated Name",
			},
			existingUser: &models.User{ID: 1, Name: "Old Name"},
			getUserErr:   nil,
			updateErr:    errors.New("database error"),
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			mockRepo.On("GetByID", tt.userID).Return(tt.existingUser, tt.getUserErr)
			
			if tt.getUserErr == nil {
				// Mock the email existence check (for email conflict detection)
				if tt.existingUser != nil && tt.request.Email != tt.existingUser.Email {
					mockRepo.On("GetByEmail", tt.request.Email).Return(nil, gorm.ErrRecordNotFound)
				}
				mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(tt.updateErr)
			}
			
			service := NewUserService(mockRepo)
			user, err := service.UpdateUser(tt.userID, tt.request)
			
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				if tt.request.Name != "" {
					assert.Equal(t, tt.request.Name, user.Name)
				}
				if tt.request.Email != "" {
					assert.Equal(t, tt.request.Email, user.Email)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		getUserErr    error
		deleteErr     error
		expectedError string
	}{
		{
			name:          "success",
			userID:        1,
			getUserErr:    nil,
			deleteErr:     nil,
			expectedError: "",
		},
		{
			name:          "user not found",
			userID:        999,
			getUserErr:    gorm.ErrRecordNotFound,
			deleteErr:     nil,
			expectedError: "user not found",
		},
		{
			name:          "delete error",
			userID:        1,
			getUserErr:    nil,
			deleteErr:     errors.New("database error"),
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			
			if tt.getUserErr == gorm.ErrRecordNotFound {
				mockRepo.On("GetByID", tt.userID).Return(nil, tt.getUserErr)
			} else {
				user := &models.User{ID: tt.userID}
				mockRepo.On("GetByID", tt.userID).Return(user, tt.getUserErr)
				mockRepo.On("Delete", tt.userID).Return(tt.deleteErr)
			}
			
			service := NewUserService(mockRepo)
			err := service.DeleteUser(tt.userID)
			
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserCount(t *testing.T) {
	tests := []struct {
		name          string
		repoCount     int64
		repoError     error
		expectedCount int64
		expectedError error
	}{
		{
			name:          "success",
			repoCount:     5,
			repoError:     nil,
			expectedCount: 5,
			expectedError: nil,
		},
		{
			name:          "repository error",
			repoCount:     0,
			repoError:     errors.New("database error"),
			expectedCount: 0,
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			mockRepo.On("Count").Return(tt.repoCount, tt.repoError)
			
			service := NewUserService(mockRepo)
			count, err := service.GetUserCount()
			
			assert.Equal(t, tt.expectedCount, count)
			assert.Equal(t, tt.expectedError, err)
			mockRepo.AssertExpectations(t)
		})
	}
}
