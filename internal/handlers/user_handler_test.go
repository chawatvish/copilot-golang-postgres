package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gin-simple-app/internal/models"
	"gin-simple-app/pkg/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) GetUserCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func setupGinTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewUserHandler(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)
	
	assert.NotNil(t, handler)
	assert.IsType(t, &UserHandler{}, handler)
}

func TestUserHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		serviceUsers   []models.User
		serviceCount   int64
		serviceErr     error
		expectedStatus int
		expectError    bool
	}{
		{
			name: "success",
			serviceUsers: []models.User{
				{ID: 1, Name: "John", Email: "john@example.com"},
				{ID: 2, Name: "Jane", Email: "jane@example.com"},
			},
			serviceCount:   2,
			serviceErr:     nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "service error",
			serviceUsers:   []models.User{},
			serviceCount:   0,
			serviceErr:     errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			mockService.On("GetAllUsers").Return(tt.serviceUsers, tt.serviceErr)
			
			handler := NewUserHandler(mockService)
			router := setupGinTest()
			router.GET("/users", handler.GetUsers)
			
			req, _ := http.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var apiResponse response.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
			assert.NoError(t, err)
			
			if tt.expectError {
				assert.False(t, apiResponse.Success)
			} else {
				assert.True(t, apiResponse.Success)
				assert.Equal(t, "Users retrieved successfully", apiResponse.Message)
				assert.Equal(t, len(tt.serviceUsers), *apiResponse.Count)
			}
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		serviceUser    *models.User
		serviceErr     error
		expectedStatus int
		expectError    bool
		errorMessage   string
	}{
		{
			name:   "success",
			userID: "1",
			serviceUser: &models.User{
				ID: 1, Name: "John", Email: "john@example.com",
			},
			serviceErr:     nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid user id",
			userID:         "invalid",
			serviceUser:    nil,
			serviceErr:     nil,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorMessage:   "Invalid user ID",
		},
		{
			name:           "user not found",
			userID:         "999",
			serviceUser:    nil,
			serviceErr:     errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
			expectError:    true,
			errorMessage:   "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			
			if tt.userID != "invalid" {
				userID := uint(1)
				if tt.userID == "999" {
					userID = 999
				}
				mockService.On("GetUserByID", userID).Return(tt.serviceUser, tt.serviceErr)
			}
			
			handler := NewUserHandler(mockService)
			router := setupGinTest()
			router.GET("/users/:id", handler.GetUserByID)
			
			req, _ := http.NewRequest("GET", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var apiResponse response.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
			assert.NoError(t, err)
			
			if tt.expectError {
				assert.False(t, apiResponse.Success)
				if tt.errorMessage != "" {
					assert.Equal(t, tt.errorMessage, apiResponse.Error)
				}
			} else {
				assert.True(t, apiResponse.Success)
				assert.Equal(t, "User retrieved successfully", apiResponse.Message)
			}
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		serviceUser    *models.User
		serviceErr     error
		expectedStatus int
		expectError    bool
		errorMessage   string
	}{
		{
			name: "success",
			requestBody: models.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    "+1234567890",
			},
			serviceUser: &models.User{
				ID: 1, Name: "John Doe", Email: "john@example.com",
			},
			serviceErr:     nil,
			expectedStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name:           "invalid json",
			requestBody:    "invalid json",
			serviceUser:    nil,
			serviceErr:     nil,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorMessage:   "Invalid JSON",
		},
		{
			name: "duplicate email",
			requestBody: models.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
				Phone:    "+1234567890",
			},
			serviceUser:    nil,
			serviceErr:     errors.New("user with this email already exists"),
			expectedStatus: http.StatusConflict,
			expectError:    true,
			errorMessage:   "User with this email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			
			if tt.name != "invalid json" {
				mockService.On("CreateUser", mock.AnythingOfType("models.CreateUserRequest")).Return(tt.serviceUser, tt.serviceErr)
			}
			
			handler := NewUserHandler(mockService)
			router := setupGinTest()
			router.POST("/users", handler.CreateUser)
			
			var body []byte
			if tt.name == "invalid json" {
				body = []byte(tt.requestBody.(string))
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}
			
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var apiResponse response.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
			assert.NoError(t, err)
			
			if tt.expectError {
				assert.False(t, apiResponse.Success)
				if tt.errorMessage != "" {
					assert.Equal(t, tt.errorMessage, apiResponse.Error)
				}
			} else {
				assert.True(t, apiResponse.Success)
				assert.Equal(t, "User created successfully", apiResponse.Message)
			}
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    interface{}
		serviceUser    *models.User
		serviceErr     error
		expectedStatus int
		expectError    bool
		errorMessage   string
	}{
		{
			name:   "success",
			userID: "1",
			requestBody: models.UpdateUserRequest{
				Name:  "John Updated",
				Email: "john.updated@example.com",
				Phone: "+9876543210",
			},
			serviceUser: &models.User{
				ID: 1, Name: "John Updated", Email: "john.updated@example.com",
			},
			serviceErr:     nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid user id",
			userID:         "invalid",
			requestBody:    models.UpdateUserRequest{},
			serviceUser:    nil,
			serviceErr:     nil,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorMessage:   "Invalid user ID",
		},
		{
			name:   "user not found",
			userID: "999",
			requestBody: models.UpdateUserRequest{
				Name:  "John Updated",
				Email: "john.updated@example.com",
				Phone: "+9876543210",
			},
			serviceUser:    nil,
			serviceErr:     errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
			expectError:    true,
			errorMessage:   "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			
			if tt.userID != "invalid" {
				userID := uint(1)
				if tt.userID == "999" {
					userID = 999
				}
				mockService.On("UpdateUser", userID, mock.AnythingOfType("models.UpdateUserRequest")).Return(tt.serviceUser, tt.serviceErr)
			}
			
			handler := NewUserHandler(mockService)
			router := setupGinTest()
			router.PUT("/users/:id", handler.UpdateUser)
			
			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/users/"+tt.userID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var apiResponse response.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
			assert.NoError(t, err)
			
			if tt.expectError {
				assert.False(t, apiResponse.Success)
				if tt.errorMessage != "" {
					assert.Equal(t, tt.errorMessage, apiResponse.Error)
				}
			} else {
				assert.True(t, apiResponse.Success)
				assert.Equal(t, "User updated successfully", apiResponse.Message)
			}
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		serviceErr     error
		expectedStatus int
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "success",
			userID:         "1",
			serviceErr:     nil,
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "invalid user id",
			userID:         "invalid",
			serviceErr:     nil,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			errorMessage:   "Invalid user ID",
		},
		{
			name:           "user not found",
			userID:         "999",
			serviceErr:     errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
			expectError:    true,
			errorMessage:   "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockUserService{}
			
			if tt.userID != "invalid" {
				userID := uint(1)
				if tt.userID == "999" {
					userID = 999
				}
				mockService.On("DeleteUser", userID).Return(tt.serviceErr)
			}
			
			handler := NewUserHandler(mockService)
			router := setupGinTest()
			router.DELETE("/users/:id", handler.DeleteUser)
			
			req, _ := http.NewRequest("DELETE", "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var apiResponse response.APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
			assert.NoError(t, err)
			
			if tt.expectError {
				assert.False(t, apiResponse.Success)
				if tt.errorMessage != "" {
					assert.Equal(t, tt.errorMessage, apiResponse.Error)
				}
			} else {
				assert.True(t, apiResponse.Success)
				assert.Equal(t, "User deleted successfully", apiResponse.Message)
			}
			
			mockService.AssertExpectations(t)
		})
	}
}
