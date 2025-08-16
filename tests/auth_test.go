package tests

import (
	"bytes"
	"encoding/json"
	"gin-simple-app/internal/models"
	"gin-simple-app/pkg/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Health check successful", response.Message)
}

func TestRootEndpoint(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Welcome", response.Message)
}

func TestRegisterUser(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	registerData := models.RegisterRequest{
		Name:            "Test User",
		Email:           "test@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
		Phone:           "+1-555-0123",
		Address:         nil,
	}

	jsonData, _ := json.Marshal(registerData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "User registered successfully", response.Message)

	// Check that the response contains user data and tokens
	data := response.Data.(map[string]interface{})
	assert.Contains(t, data, "user")
	assert.Contains(t, data, "access_token")
	assert.Contains(t, data, "token_type")
	assert.Equal(t, "Bearer", data["token_type"])

	user := data["user"].(map[string]interface{})
	assert.Equal(t, "Test User", user["name"])
	assert.Equal(t, "test@example.com", user["email"])
}

func TestRegisterUserPasswordMismatch(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	registerData := models.RegisterRequest{
		Name:            "Test User",
		Email:           "test@example.com",
		Password:        "password123",
		ConfirmPassword: "password456", // Different password
		Phone:           "+1-555-0123",
		Address:         nil,
	}

	jsonData, _ := json.Marshal(registerData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "passwords do not match", response.Error)
}

func TestRegisterUserDuplicateEmail(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	// Try to register with an email that already exists in test data
	registerData := models.RegisterRequest{
		Name:            "Another John",
		Email:           "john@example.com", // This email already exists in test data
		Password:        "password123",
		ConfirmPassword: "password123",
		Phone:           "+1-555-9999",
		Address:         nil,
	}

	jsonData, _ := json.Marshal(registerData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "user with this email already exists", response.Error)
}

func TestLoginUser(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	// Use one of the sample users (password for sample data is "password123")
	loginData := map[string]string{
		"email":    "john@example.com",
		"password": "password123",
	}

	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Login successful", response.Message)

	data := response.Data.(map[string]interface{})
	assert.Contains(t, data, "user")
	assert.Contains(t, data, "access_token")
	assert.Contains(t, data, "token_type")
	assert.Equal(t, "Bearer", data["token_type"])

	user := data["user"].(map[string]interface{})
	assert.Equal(t, "John Doe", user["name"])
	assert.Equal(t, "john@example.com", user["email"])
}

func TestLoginUserInvalidCredentials(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	loginData := map[string]string{
		"email":    "john@example.com",
		"password": "wrongpassword",
	}

	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "invalid email or password", response.Error)
}

func TestGetCurrentUser(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	// Login to get token
	token, err := app.loginUser("john@example.com", "password123")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response response.APIResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "User information retrieved", response.Message)

	userData := response.Data.(map[string]interface{})
	assert.Equal(t, "John Doe", userData["name"])
	assert.Equal(t, "john@example.com", userData["email"])
}

func TestGetCurrentUserUnauthorized(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Authorization header required", response.Error)
}
