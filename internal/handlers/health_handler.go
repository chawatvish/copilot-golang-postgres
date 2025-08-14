package handlers

import (
	"gin-simple-app/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health-related HTTP requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	healthData := gin.H{
		"status":  "ok",
		"message": "Gin REST API is running",
	}
	response.Success(c, http.StatusOK, "Health check successful", healthData)
}

// Root handles GET /
func (h *HealthHandler) Root(c *gin.Context) {
	rootData := gin.H{
		"message": "Welcome to Gin Simple REST API",
		"version": "1.0.0",
	}
	response.Success(c, http.StatusOK, "Welcome", rootData)
}
