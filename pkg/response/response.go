package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Count   *int        `json:"count,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}

// SuccessWithCount sends a successful response with count
func SuccessWithCount(c *gin.Context, statusCode int, message string, data interface{}, count int) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Count:   &count,
	}
	c.JSON(statusCode, response)
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	c.JSON(statusCode, response)
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, err error) {
	response := APIResponse{
		Success: false,
		Error:   err.Error(),
	}
	c.JSON(http.StatusBadRequest, response)
}

// NotFound sends a not found error response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// BadRequest sends a bad request error response
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// InternalServerError sends an internal server error response
func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
