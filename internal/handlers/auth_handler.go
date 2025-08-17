package handlers

import (
	"gin-simple-app/internal/models"
	"gin-simple-app/internal/services"
	"gin-simple-app/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Register(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", result)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", result)
}

// ForgotPassword handles forgot password request
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := h.authService.ForgotPassword(req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "If the email exists, a password reset link has been sent", nil)
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := h.authService.ResetPassword(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password reset successfully", nil)
}

// ChangePassword handles password change for authenticated users
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err := h.authService.ChangePassword(userID.(uint), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	result, err := h.authService.RefreshToken(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", result)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	err := h.authService.Logout(userID.(uint))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Logged out successfully", nil)
}

// Me returns the current authenticated user's information
func (h *AuthHandler) Me(c *gin.Context) {
	// Get user information from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Convert to UserResponse format to hide sensitive information
	u := user.(*models.User)
	userResponse := &models.UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Phone:           u.Phone,
		Address:         u.Address,
		IsActive:        u.IsActive,
		IsEmailVerified: u.IsEmailVerified,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
		LastLoginAt:     u.LastLoginAt,
	}

	response.Success(c, http.StatusOK, "User information retrieved", userResponse)
}

// AuthMiddleware validates JWT token and sets user context
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		// Check if token has Bearer prefix
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Extract user ID from claims
		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		userID := uint(userIDFloat)

		// Set user ID in context
		c.Set("user_id", userID)
		c.Set("email", (*claims)["email"])

		c.Next()
	}
}

// OptionalAuthMiddleware is like AuthMiddleware but doesn't require authentication
// It sets user context if token is provided and valid
func OptionalAuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := tokenParts[1]
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		if userIDFloat, ok := (*claims)["user_id"].(float64); ok {
			userID := uint(userIDFloat)
			c.Set("user_id", userID)
			c.Set("email", (*claims)["email"])
		}

		c.Next()
	}
}

// AdminMiddleware checks if user has admin privileges (you can extend User model to include roles)
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// This is a placeholder - you would implement role-based access control
		// For now, we'll just check if user is authenticated
		_, exists := c.Get("user_id")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Authentication required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// EnhancedAuthMiddleware validates JWT token, fetches and sets user object in context
func EnhancedAuthMiddleware(authService services.AuthService, userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		userIDFloat, ok := (*claims)["user_id"].(float64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		userID := uint(userIDFloat)

		// Fetch user from database
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "User not found")
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			response.Error(c, http.StatusUnauthorized, "Account is deactivated")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", user.Email)
		c.Set("user", user)

		c.Next()
	}
}
