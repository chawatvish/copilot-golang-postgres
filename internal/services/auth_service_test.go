package services

import (
	"errors"
	"gin-simple-app/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestNewAuthService(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)
	
	assert.NotNil(t, service)
	assert.IsType(t, &AuthServiceImpl{}, service)
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name           string
		request        models.RegisterRequest
		emailExistsErr error
		createErr      error
		expectError    bool
		errorMessage   string
	}{
		{
			name: "success",
			request: models.RegisterRequest{
				Name:            "John Doe",
				Email:           "john@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				Phone:           "+1234567890",
				Address:         stringPointer("123 Main St"),
			},
			emailExistsErr: gorm.ErrRecordNotFound,
			createErr:      nil,
			expectError:    false,
		},
		{
			name: "password mismatch",
			request: models.RegisterRequest{
				Name:            "John Doe",
				Email:           "john@example.com",
				Password:        "password123",
				ConfirmPassword: "password456",
				Phone:           "+1234567890",
			},
			expectError:  true,
			errorMessage: "passwords do not match",
		},
		{
			name: "user already exists",
			request: models.RegisterRequest{
				Name:            "John Doe",
				Email:           "existing@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				Phone:           "+1234567890",
			},
			emailExistsErr: nil,
			expectError:    true,
			errorMessage:   "user with this email already exists",
		},
		{
			name: "repository create error",
			request: models.RegisterRequest{
				Name:            "John Doe",
				Email:           "john@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				Phone:           "+1234567890",
			},
			emailExistsErr: gorm.ErrRecordNotFound,
			createErr:      errors.New("database error"),
			expectError:    true,
			errorMessage:   "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			
			// Only mock repository calls if password matches
			if tt.request.Password == tt.request.ConfirmPassword {
				if tt.emailExistsErr == gorm.ErrRecordNotFound {
					mockRepo.On("GetByEmail", tt.request.Email).Return(nil, tt.emailExistsErr)
				} else {
					existingUser := &models.User{Email: tt.request.Email}
					mockRepo.On("GetByEmail", tt.request.Email).Return(existingUser, tt.emailExistsErr)
				}
				
				if tt.emailExistsErr == gorm.ErrRecordNotFound {
					mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(tt.createErr)
					// Mock the Update call for LastLoginAt
					if tt.createErr == nil {
						mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)
					}
				}
			}
			
			service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)
			response, err := service.Register(tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.AccessToken)
				assert.Equal(t, "Bearer", response.TokenType)
				assert.Equal(t, tt.request.Name, response.User.Name)
				assert.Equal(t, tt.request.Email, response.User.Email)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	// Create a test user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := &models.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: string(hashedPassword),
		IsActive: true,
	}

	tests := []struct {
		name          string
		request       models.LoginRequest
		getUserResp   *models.User
		getUserErr    error
		expectError   bool
		errorMessage  string
	}{
		{
			name: "success",
			request: models.LoginRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			getUserResp:  testUser,
			getUserErr:   nil,
			expectError:  false,
		},
		{
			name: "user not found",
			request: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			getUserResp:  nil,
			getUserErr:   gorm.ErrRecordNotFound,
			expectError:  true,
			errorMessage: "invalid email or password",
		},
		{
			name: "incorrect password",
			request: models.LoginRequest{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			getUserResp:  testUser,
			getUserErr:   nil,
			expectError:  true,
			errorMessage: "invalid email or password",
		},
		{
			name: "inactive user",
			request: models.LoginRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			getUserResp: &models.User{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: string(hashedPassword),
				IsActive: false,
			},
			getUserErr:   nil,
			expectError:  true,
			errorMessage: "account is deactivated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			mockRepo.On("GetByEmail", tt.request.Email).Return(tt.getUserResp, tt.getUserErr)
			
			// Mock Update call for successful login (to update last login)
			if !tt.expectError && tt.getUserResp != nil {
				mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil)
			}
			
			service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)
			response, err := service.Login(tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.AccessToken)
				assert.Equal(t, "Bearer", response.TokenType)
				assert.Equal(t, tt.getUserResp.Name, response.User.Name)
				assert.Equal(t, tt.getUserResp.Email, response.User.Email)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockRepo := &MockUserRepository{}
	service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)
	
	// Generate a valid token first
	user := &models.User{ID: 1, Email: "john@example.com"}
	authService := service.(*AuthServiceImpl)
	token, _, err := authService.generateJWTToken(user.ID, user.Email)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		token        string
		expectError  bool
		errorMessage string
	}{
		{
			name:        "valid token",
			token:       token,
			expectError: false,
		},
		{
			name:         "invalid token",
			token:        "invalid.token.here",
			expectError:  true,
			errorMessage: "malformed",
		},
		{
			name:         "empty token",
			token:        "",
			expectError:  true,
			errorMessage: "malformed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				
				// Extract user ID from claims
				userIDFloat, exists := (*claims)["user_id"]
				assert.True(t, exists)
				assert.Equal(t, float64(1), userIDFloat)
			}
		})
	}
}

func TestAuthService_ForgotPassword(t *testing.T) {
	testUser := &models.User{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}

	tests := []struct {
		name          string
		request       models.ForgotPasswordRequest
		getUserResp   *models.User
		getUserErr    error
		updateErr     error
		expectError   bool
		errorMessage  string
	}{
		{
			name: "success",
			request: models.ForgotPasswordRequest{
				Email: "john@example.com",
			},
			getUserResp: testUser,
			getUserErr:  nil,
			updateErr:   nil,
			expectError: false,
		},
		{
			name: "user not found",
			request: models.ForgotPasswordRequest{
				Email: "nonexistent@example.com",
			},
			getUserResp:  nil,
			getUserErr:   gorm.ErrRecordNotFound,
			expectError:  false, // ForgotPassword returns nil even if user not found
		},
		{
			name: "update error",
			request: models.ForgotPasswordRequest{
				Email: "john@example.com",
			},
			getUserResp:  testUser,
			getUserErr:   nil,
			updateErr:    errors.New("database error"),
			expectError:  true,
			errorMessage: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			mockRepo.On("GetByEmail", tt.request.Email).Return(tt.getUserResp, tt.getUserErr)
			
			if tt.getUserResp != nil {
				mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(tt.updateErr)
			}
			
			service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)
			err := service.ForgotPassword(tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_ChangePassword(t *testing.T) {
	// Create a test user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	testUser := &models.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: string(hashedPassword),
	}

	tests := []struct {
		name         string
		userID       uint
		request      models.ChangePasswordRequest
		getUserResp  *models.User
		getUserErr   error
		updateErr    error
		expectError  bool
		errorMessage string
	}{
		{
			name:   "success",
			userID: 1,
			request: models.ChangePasswordRequest{
				CurrentPassword: "oldpassword",
				NewPassword:     "newpassword123",
				ConfirmPassword: "newpassword123",
			},
			getUserResp: testUser,
			getUserErr:  nil,
			updateErr:   nil,
			expectError: false,
		},
		{
			name:   "password mismatch",
			userID: 1,
			request: models.ChangePasswordRequest{
				CurrentPassword: "oldpassword",
				NewPassword:     "newpassword123",
				ConfirmPassword: "differentpassword",
			},
			expectError:  true,
			errorMessage: "passwords do not match",
		},
		{
			name:   "wrong current password",
			userID: 1,
			request: models.ChangePasswordRequest{
				CurrentPassword: "wrongpassword",
				NewPassword:     "newpassword123",
				ConfirmPassword: "newpassword123",
			},
			getUserResp:  testUser,
			getUserErr:   nil,
			expectError:  true,
			errorMessage: "current password is incorrect",
		},
		{
			name:   "user not found",
			userID: 999,
			request: models.ChangePasswordRequest{
				CurrentPassword: "oldpassword",
				NewPassword:     "newpassword123",
				ConfirmPassword: "newpassword123",
			},
			getUserResp:  nil,
			getUserErr:   gorm.ErrRecordNotFound,
			expectError:  true,
			errorMessage: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{}
			
			// Only mock repository calls if passwords match
			if tt.request.NewPassword == tt.request.ConfirmPassword {
				mockRepo.On("GetByID", tt.userID).Return(tt.getUserResp, tt.getUserErr)
				
				if tt.getUserResp != nil && tt.request.CurrentPassword == "oldpassword" {
					mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(tt.updateErr)
				}
			}
			
			service := NewAuthService(mockRepo, "test-secret", 24*time.Hour)
			err := service.ChangePassword(tt.userID, tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}
