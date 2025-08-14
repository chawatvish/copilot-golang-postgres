package tests

import (
	"bytes"
	"encoding/json"
	"gin-simple-app/internal/handlers"
	"gin-simple-app/internal/repository"
	"gin-simple-app/internal/router"
	"gin-simple-app/internal/services"
	"gin-simple-app/pkg/response"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestApp holds the application components for testing
type TestApp struct {
	router   *gin.Engine
	userRepo *repository.InMemoryUserRepository
}

// setupTestApp initializes the application for testing
func setupTestApp() *TestApp {
	gin.SetMode(gin.TestMode)

	// Initialize components with in-memory repository for testing
	userRepo := repository.NewInMemoryUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	healthHandler := handlers.NewHealthHandler()
	appRouter := router.NewRouter(userHandler, healthHandler)

	return &TestApp{
		router:   appRouter.SetupRoutes(),
		userRepo: userRepo,
	}
}

// resetTestData resets the test data before each test
func (app *TestApp) resetTestData() {
	app.userRepo.Reset()
}

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

	data := response.Data.(map[string]interface{})
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "Gin REST API is running", data["message"])
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

	data := response.Data.(map[string]interface{})
	assert.Equal(t, "Welcome to Gin Simple REST API", data["message"])
	assert.Equal(t, "1.0.0", data["version"])
}

func TestGetAllUsers(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Users retrieved successfully", response.Message)
	assert.Equal(t, 3, *response.Count)

	data := response.Data.([]interface{})
	assert.Len(t, data, 3)

	// Check first user
	firstUser := data[0].(map[string]interface{})
	assert.Equal(t, float64(1), firstUser["id"])
	assert.Equal(t, "John Doe", firstUser["name"])
	assert.Equal(t, "john@example.com", firstUser["email"])
}

func TestGetUserByID(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "User retrieved successfully", response.Message)

	userData := response.Data.(map[string]interface{})
	assert.Equal(t, float64(1), userData["id"])
	assert.Equal(t, "John Doe", userData["name"])
	assert.Equal(t, "john@example.com", userData["email"])
}

func TestGetUserByIDNotFound(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "User not found", response.Error)
}

func TestGetUserByIDInvalid(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/invalid", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "Invalid user ID", response.Error)
}

func TestCreateUser(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	newUser := map[string]string{
		"name":  "Alice Cooper",
		"email": "alice@example.com",
	}

	jsonData, _ := json.Marshal(newUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "User created successfully", response.Message)

	userData := response.Data.(map[string]interface{})
	assert.Equal(t, float64(4), userData["id"]) // Should be assigned ID 4
	assert.Equal(t, "Alice Cooper", userData["name"])
	assert.Equal(t, "alice@example.com", userData["email"])
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	newUser := map[string]string{
		"name":  "John Smith",
		"email": "john@example.com", // This email already exists
	}

	jsonData, _ := json.Marshal(newUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "User with this email already exists", response.Error)
}

func TestCreateUserInvalidJSON(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "invalid character")
}

func TestCreateUserMissingFields(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	newUser := map[string]string{
		"name": "Alice Cooper",
		// missing email
	}

	jsonData, _ := json.Marshal(newUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error, "required")
}

func TestUpdateUser(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	updatedUser := map[string]string{
		"name":  "John Updated",
		"email": "john.updated@example.com",
	}

	jsonData, _ := json.Marshal(updatedUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "User updated successfully", response.Message)

	userData := response.Data.(map[string]interface{})
	assert.Equal(t, float64(1), userData["id"])
	assert.Equal(t, "John Updated", userData["name"])
	assert.Equal(t, "john.updated@example.com", userData["email"])
}

func TestUpdateUserNotFound(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	updatedUser := map[string]string{
		"name":  "Non Existent",
		"email": "nonexistent@example.com",
	}

	jsonData, _ := json.Marshal(updatedUser)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "User not found", response.Error)
}

func TestDeleteUser(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	// First verify the user exists
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/v1/users/3", nil)
	app.router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Delete the user
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/3", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "User deleted successfully", response.Message)

	// Verify the user is deleted
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/users/3", nil)
	app.router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestDeleteUserNotFound(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/999", nil)
	app.router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "User not found", response.Error)
}

func TestCompleteUserLifecycle(t *testing.T) {
	app := setupTestApp()
	app.resetTestData()

	// 1. Get initial user count
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/v1/users", nil)
	app.router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	var initialResponse response.APIResponse
	err := json.Unmarshal(w1.Body.Bytes(), &initialResponse)
	assert.NoError(t, err)
	initialCount := *initialResponse.Count

	// 2. Create a new user
	newUser := map[string]string{
		"name":  "Lifecycle Test",
		"email": "lifecycle@example.com",
	}
	jsonData, _ := json.Marshal(newUser)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req2.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)

	var createResponse response.APIResponse
	err = json.Unmarshal(w2.Body.Bytes(), &createResponse)
	assert.NoError(t, err)

	userData := createResponse.Data.(map[string]interface{})
	userID := uint(userData["id"].(float64))

	// 3. Verify user count increased
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/users", nil)
	app.router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	var afterCreateResponse response.APIResponse
	err = json.Unmarshal(w3.Body.Bytes(), &afterCreateResponse)
	assert.NoError(t, err)
	afterCreateCount := *afterCreateResponse.Count
	assert.Equal(t, initialCount+1, afterCreateCount)

	// 4. Update the user
	updatedUser := map[string]string{
		"name":  "Lifecycle Updated",
		"email": "lifecycle.updated@example.com",
	}
	updateData, _ := json.Marshal(updatedUser)

	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("PUT", "/api/v1/users/"+strconv.FormatUint(uint64(userID), 10), bytes.NewBuffer(updateData))
	req4.Header.Set("Content-Type", "application/json")
	app.router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)

	// 5. Verify the update
	w5 := httptest.NewRecorder()
	req5, _ := http.NewRequest("GET", "/api/v1/users/"+strconv.FormatUint(uint64(userID), 10), nil)
	app.router.ServeHTTP(w5, req5)
	assert.Equal(t, http.StatusOK, w5.Code)

	var getResponse response.APIResponse
	err = json.Unmarshal(w5.Body.Bytes(), &getResponse)
	assert.NoError(t, err)

	retrievedUser := getResponse.Data.(map[string]interface{})
	assert.Equal(t, "Lifecycle Updated", retrievedUser["name"])
	assert.Equal(t, "lifecycle.updated@example.com", retrievedUser["email"])

	// 6. Delete the user
	w6 := httptest.NewRecorder()
	req6, _ := http.NewRequest("DELETE", "/api/v1/users/"+strconv.FormatUint(uint64(userID), 10), nil)
	app.router.ServeHTTP(w6, req6)
	assert.Equal(t, http.StatusOK, w6.Code)

	// 7. Verify user count returned to initial
	w7 := httptest.NewRecorder()
	req7, _ := http.NewRequest("GET", "/api/v1/users", nil)
	app.router.ServeHTTP(w7, req7)
	assert.Equal(t, http.StatusOK, w7.Code)

	var finalResponse response.APIResponse
	err = json.Unmarshal(w7.Body.Bytes(), &finalResponse)
	assert.NoError(t, err)
	finalCount := *finalResponse.Count
	assert.Equal(t, initialCount, finalCount)
}
