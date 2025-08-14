package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestRouter sets up Gin in test mode and returns the router
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return setupRouter()
}

// TestMain runs before all tests to set up the environment
func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Run tests
	m.Run()
}

// Helper function to reset data before each test
func setupTestData() {
	resetUsers()
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "Gin REST API is running", response["message"])
}

// TestRootEndpoint tests the root endpoint
func TestRootEndpoint(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Welcome to Gin Simple REST API", response["message"])
	assert.Equal(t, "1.0.0", response["version"])
}

// TestGetAllUsers tests getting all users
func TestGetAllUsers(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, float64(3), response["count"])
	
	data, exists := response["data"].([]interface{})
	assert.True(t, exists)
	assert.Len(t, data, 3)
	
	// Check first user
	firstUser := data[0].(map[string]interface{})
	assert.Equal(t, float64(1), firstUser["id"])
	assert.Equal(t, "John Doe", firstUser["name"])
	assert.Equal(t, "john@example.com", firstUser["email"])
}

// TestGetUserByID tests getting a user by ID
func TestGetUserByID(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	// Test getting existing user
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	userData := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), userData["id"])
	assert.Equal(t, "John Doe", userData["name"])
	assert.Equal(t, "john@example.com", userData["email"])
}

// TestGetUserByIDNotFound tests getting a non-existent user
func TestGetUserByIDNotFound(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])
}

// TestGetUserByIDInvalid tests getting a user with invalid ID
func TestGetUserByIDInvalid(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/users/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])
}

// TestCreateUser tests creating a new user
func TestCreateUser(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	newUser := User{
		Name:  "Alice Cooper",
		Email: "alice@example.com",
	}
	
	jsonData, _ := json.Marshal(newUser)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "User created successfully", response["message"])
	
	userData := response["data"].(map[string]interface{})
	assert.Equal(t, float64(4), userData["id"]) // Should be assigned ID 4
	assert.Equal(t, "Alice Cooper", userData["name"])
	assert.Equal(t, "alice@example.com", userData["email"])
}

// TestCreateUserInvalidJSON tests creating a user with invalid JSON
func TestCreateUserInvalidJSON(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "invalid character")
}

// TestUpdateUser tests updating an existing user
func TestUpdateUser(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	updatedUser := User{
		Name:  "John Updated",
		Email: "john.updated@example.com",
	}
	
	jsonData, _ := json.Marshal(updatedUser)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	assert.Equal(t, "User updated successfully", response["message"])
	
	userData := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), userData["id"])
	assert.Equal(t, "John Updated", userData["name"])
	assert.Equal(t, "john.updated@example.com", userData["email"])
}

// TestUpdateUserNotFound tests updating a non-existent user
func TestUpdateUserNotFound(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	updatedUser := User{
		Name:  "Non Existent",
		Email: "nonexistent@example.com",
	}
	
	jsonData, _ := json.Marshal(updatedUser)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])
}

// TestUpdateUserInvalidID tests updating a user with invalid ID
func TestUpdateUserInvalidID(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	updatedUser := User{
		Name:  "Test User",
		Email: "test@example.com",
	}
	
	jsonData, _ := json.Marshal(updatedUser)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/users/invalid", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])
}

// TestDeleteUser tests deleting an existing user
func TestDeleteUser(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	// First verify the user exists
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/v1/users/3", nil)
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Delete the user
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/3", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User deleted successfully", response["message"])

	// Verify the user is deleted
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/v1/users/3", nil)
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

// TestDeleteUserNotFound tests deleting a non-existent user
func TestDeleteUserNotFound(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User not found", response["error"])
}

// TestDeleteUserInvalidID tests deleting a user with invalid ID
func TestDeleteUserInvalidID(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/v1/users/invalid", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid user ID", response["error"])
}

// TestCompleteUserLifecycle tests the complete CRUD lifecycle
func TestCompleteUserLifecycle(t *testing.T) {
	setupTestData()
	router := setupTestRouter()

	// 1. Get initial user count
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)
	
	var initialResponse map[string]interface{}
	err := json.Unmarshal(w1.Body.Bytes(), &initialResponse)
	assert.NoError(t, err)
	initialCount := int(initialResponse["count"].(float64))

	// 2. Create a new user
	newUser := User{Name: "Lifecycle Test", Email: "lifecycle@example.com"}
	jsonData, _ := json.Marshal(newUser)
	
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonData))
	req2.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusCreated, w2.Code)
	
	var createResponse map[string]interface{}
	err = json.Unmarshal(w2.Body.Bytes(), &createResponse)
	assert.NoError(t, err)
	
	userData := createResponse["data"].(map[string]interface{})
	userID := int(userData["id"].(float64))

	// 3. Verify user count increased
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
	
	var afterCreateResponse map[string]interface{}
	err = json.Unmarshal(w3.Body.Bytes(), &afterCreateResponse)
	assert.NoError(t, err)
	afterCreateCount := int(afterCreateResponse["count"].(float64))
	assert.Equal(t, initialCount+1, afterCreateCount)

	// 4. Update the user
	updatedUser := User{Name: "Lifecycle Updated", Email: "lifecycle.updated@example.com"}
	updateData, _ := json.Marshal(updatedUser)
	
	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("PUT", "/api/v1/users/"+strconv.Itoa(userID), bytes.NewBuffer(updateData))
	req4.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)

	// 5. Verify the update
	w5 := httptest.NewRecorder()
	req5, _ := http.NewRequest("GET", "/api/v1/users/"+strconv.Itoa(userID), nil)
	router.ServeHTTP(w5, req5)
	assert.Equal(t, http.StatusOK, w5.Code)
	
	var getResponse map[string]interface{}
	err = json.Unmarshal(w5.Body.Bytes(), &getResponse)
	assert.NoError(t, err)
	
	retrievedUser := getResponse["data"].(map[string]interface{})
	assert.Equal(t, "Lifecycle Updated", retrievedUser["name"])
	assert.Equal(t, "lifecycle.updated@example.com", retrievedUser["email"])

	// 6. Delete the user
	w6 := httptest.NewRecorder()
	req6, _ := http.NewRequest("DELETE", "/api/v1/users/"+strconv.Itoa(userID), nil)
	router.ServeHTTP(w6, req6)
	assert.Equal(t, http.StatusOK, w6.Code)

	// 7. Verify user count returned to initial
	w7 := httptest.NewRecorder()
	req7, _ := http.NewRequest("GET", "/api/v1/users", nil)
	router.ServeHTTP(w7, req7)
	assert.Equal(t, http.StatusOK, w7.Code)
	
	var finalResponse map[string]interface{}
	err = json.Unmarshal(w7.Body.Bytes(), &finalResponse)
	assert.NoError(t, err)
	finalCount := int(finalResponse["count"].(float64))
	assert.Equal(t, initialCount, finalCount)
}
