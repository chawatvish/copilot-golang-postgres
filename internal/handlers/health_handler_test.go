package handlers

import (
	"encoding/json"
	"gin-simple-app/pkg/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHealthHandler(t *testing.T) {
	handler := NewHealthHandler()
	
	assert.NotNil(t, handler)
	assert.IsType(t, &HealthHandler{}, handler)
}

func TestHealthHandler_HealthCheck(t *testing.T) {
	handler := NewHealthHandler()
	router := setupGinTest()
	router.GET("/health", handler.HealthCheck)
	
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var apiResponse response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
	assert.NoError(t, err)
	
	assert.True(t, apiResponse.Success)
	assert.Equal(t, "Health check successful", apiResponse.Message)
	
	data := apiResponse.Data.(map[string]interface{})
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "Gin REST API is running", data["message"])
}

func TestHealthHandler_Root(t *testing.T) {
	handler := NewHealthHandler()
	router := setupGinTest()
	router.GET("/", handler.Root)
	
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var apiResponse response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
	assert.NoError(t, err)
	
	assert.True(t, apiResponse.Success)
	assert.Equal(t, "Welcome", apiResponse.Message)
	
	data := apiResponse.Data.(map[string]interface{})
	assert.Equal(t, "Welcome to Gin Simple REST API", data["message"])
	assert.Equal(t, "1.0.0", data["version"])
}
