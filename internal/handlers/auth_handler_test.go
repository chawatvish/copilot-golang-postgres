package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gin-simple-app/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the AuthService interface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req models.RegisterRequest) (*models.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*jwt.MapClaims), args.Error(1)
}

func (m *MockAuthService) ForgotPassword(req models.ForgotPasswordRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) ChangePassword(userID uint, req models.ChangePasswordRequest) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

func (m *MockAuthService) ResetPassword(req models.ResetPasswordRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) RefreshToken(userID uint) (*models.LoginResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginResponse), args.Error(1)
}

func (m *MockAuthService) Logout(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func TestNewAuthHandler(t *testing.T) {
	mockAuthService := &MockAuthService{}
	handler := NewAuthHandler(mockAuthService)
	
	assert.NotNil(t, handler)
	assert.Equal(t, mockAuthService, handler.authService)
}

func TestAuthHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "success",
			requestBody: models.RegisterRequest{
				Name:            "John Doe",
				Email:           "john@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				Phone:           "1234567890",
			},
			mockSetup: func(m *MockAuthService) {
				loginResponse := &models.LoginResponse{
					User: &models.User{
						ID:    1,
						Name:  "John Doe",
						Email: "john@example.com",
					},
					AccessToken: "mock-token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				}
				m.On("Register", mock.AnythingOfType("models.RegisterRequest")).Return(loginResponse, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid_json",
			requestBody: "invalid json",
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed as it should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing_required_fields",
			requestBody: map[string]interface{}{
				"name": "John Doe",
				// Missing email, password, confirm_password
			},
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed as validation should fail
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			requestBody: models.RegisterRequest{
				Name:            "John Doe",
				Email:           "john@example.com",
				Password:        "password123",
				ConfirmPassword: "password123",
				Phone:           "1234567890",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Register", mock.AnythingOfType("models.RegisterRequest")).Return(nil, errors.New("email already exists"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &MockAuthService{}
			tt.mockSetup(mockAuthService)

			handler := NewAuthHandler(mockAuthService)
			router := gin.New()
			router.POST("/auth/register", handler.Register)

			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				jsonBody, _ := json.Marshal(tt.requestBody)
				body.Write(jsonBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/register", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "success",
			requestBody: models.LoginRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				loginResponse := &models.LoginResponse{
					User: &models.User{
						ID:    1,
						Name:  "John Doe",
						Email: "john@example.com",
					},
					AccessToken: "mock-jwt-token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				}
				m.On("Login", mock.AnythingOfType("models.LoginRequest")).Return(loginResponse, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_json",
			requestBody: "invalid json",
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing_required_fields",
			requestBody: map[string]interface{}{
				"email": "john@example.com",
				// Missing password
			},
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid_credentials",
			requestBody: models.LoginRequest{
				Email:    "john@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("Login", mock.AnythingOfType("models.LoginRequest")).Return(nil, errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &MockAuthService{}
			tt.mockSetup(mockAuthService)

			handler := NewAuthHandler(mockAuthService)
			router := gin.New()
			router.POST("/auth/login", handler.Login)

			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				jsonBody, _ := json.Marshal(tt.requestBody)
				body.Write(jsonBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/login", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Me(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *gin.Context) {
				user := &models.User{
					ID:    1,
					Name:  "John Doe",
					Email: "john@example.com",
				}
				c.Set("user", user)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "user_not_in_context",
			setupContext: func(c *gin.Context) {
				// Don't set user in context
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &MockAuthService{}
			handler := NewAuthHandler(mockAuthService)

			router := gin.New()
			router.GET("/auth/me", func(c *gin.Context) {
				tt.setupContext(c)
				handler.Me(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAuthHandler_ForgotPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "success",
			requestBody: models.ForgotPasswordRequest{
				Email: "john@example.com",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("ForgotPassword", mock.AnythingOfType("models.ForgotPasswordRequest")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_json",
			requestBody: "invalid json",
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing_email",
			requestBody: map[string]interface{}{
				// Missing email
			},
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service_error",
			requestBody: models.ForgotPasswordRequest{
				Email: "nonexistent@example.com",
			},
			mockSetup: func(m *MockAuthService) {
				m.On("ForgotPassword", mock.AnythingOfType("models.ForgotPasswordRequest")).Return(errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &MockAuthService{}
			tt.mockSetup(mockAuthService)

			handler := NewAuthHandler(mockAuthService)
			router := gin.New()
			router.POST("/auth/forgot-password", handler.ForgotPassword)

			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				jsonBody, _ := json.Marshal(tt.requestBody)
				body.Write(jsonBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/forgot-password", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockAuthService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupContext   func(*gin.Context)
		mockSetup      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "success",
			requestBody: models.ChangePasswordRequest{
				CurrentPassword: "oldpassword",
				NewPassword:     "newpassword",
				ConfirmPassword: "newpassword",
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", uint(1))
			},
			mockSetup: func(m *MockAuthService) {
				m.On("ChangePassword", uint(1), mock.AnythingOfType("models.ChangePasswordRequest")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid_json",
			requestBody: "invalid json",
			setupContext: func(c *gin.Context) {
				c.Set("user_id", uint(1))
			},
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing_required_fields",
			requestBody: map[string]interface{}{
				"current_password": "oldpassword",
				// Missing new_password and confirm_password
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", uint(1))
			},
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "user_not_authenticated",
			requestBody: models.ChangePasswordRequest{
				CurrentPassword: "oldpassword",
				NewPassword:     "newpassword",
				ConfirmPassword: "newpassword",
			},
			setupContext: func(c *gin.Context) {
				// Don't set user_id in context
			},
			mockSetup: func(m *MockAuthService) {
				// No mock setup needed
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "service_error",
			requestBody: models.ChangePasswordRequest{
				CurrentPassword: "wrongpassword",
				NewPassword:     "newpassword",
				ConfirmPassword: "newpassword",
			},
			setupContext: func(c *gin.Context) {
				c.Set("user_id", uint(1))
			},
			mockSetup: func(m *MockAuthService) {
				m.On("ChangePassword", uint(1), mock.AnythingOfType("models.ChangePasswordRequest")).Return(errors.New("current password is incorrect"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthService := &MockAuthService{}
			tt.mockSetup(mockAuthService)

			handler := NewAuthHandler(mockAuthService)
			router := gin.New()
			router.POST("/auth/change-password", func(c *gin.Context) {
				tt.setupContext(c)
				handler.ChangePassword(c)
			})

			var body bytes.Buffer
			if str, ok := tt.requestBody.(string); ok {
				body.WriteString(str)
			} else {
				jsonBody, _ := json.Marshal(tt.requestBody)
				body.Write(jsonBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/auth/change-password", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockAuthService.AssertExpectations(t)
		})
	}
}
