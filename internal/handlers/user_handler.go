package handlers

import (
	"gin-simple-app/internal/models"
	"gin-simple-app/internal/services"
	"gin-simple-app/pkg/response"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers handles GET /api/v1/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve users")
		return
	}
	
	response.SuccessWithCount(c, http.StatusOK, "Users retrieved successfully", users, len(users))
}

// GetUserByID handles GET /api/v1/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to retrieve user")
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Check if it's a JSON parsing error or validation error
		if err.Error() == "EOF" || strings.Contains(err.Error(), "invalid character") {
			response.BadRequest(c, "Invalid JSON")
		} else {
			response.ValidationError(c, err)
		}
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			response.Conflict(c, "User with this email already exists")
			return
		}
		response.InternalServerError(c, "Failed to create user")
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", user)
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, err := h.userService.UpdateUser(uint(id), req)
	if err != nil {
		if err.Error() == "user not found" {
			response.NotFound(c, "User not found")
			return
		}
		if err.Error() == "user with this email already exists" {
			response.BadRequest(c, "User with this email already exists")
			return
		}
		response.InternalServerError(c, "Failed to update user")
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		if err.Error() == "user not found" {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to delete user")
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
