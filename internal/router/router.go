package router

import (
	"gin-simple-app/internal/handlers"
	"gin-simple-app/internal/services"

	"github.com/gin-gonic/gin"
)

// Router holds all the handlers and provides methods to set up routes
type Router struct {
	userHandler   *handlers.UserHandler
	authHandler   *handlers.AuthHandler
	healthHandler *handlers.HealthHandler
	authService   services.AuthService
	userService   services.UserService
}

// NewRouter creates a new router with all handlers
func NewRouter(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	healthHandler *handlers.HealthHandler,
	authService services.AuthService,
	userService services.UserService,
) *Router {
	return &Router{
		userHandler:   userHandler,
		authHandler:   authHandler,
		healthHandler: healthHandler,
		authService:   authService,
		userService:   userService,
	}
}

// SetupRoutes configures all routes and returns a Gin engine
func (r *Router) SetupRoutes() *gin.Engine {
	// Create Gin router with default middleware (logger and recovery)
	engine := gin.Default()

	// Health and root endpoints
	engine.GET("/", r.healthHandler.Root)
	engine.GET("/health", r.healthHandler.HealthCheck)

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Public authentication routes (no middleware required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/forgot-password", r.authHandler.ForgotPassword)
			auth.POST("/reset-password", r.authHandler.ResetPassword)
		}

		// Protected authentication routes (require authentication)
		authProtected := v1.Group("/auth")
		authProtected.Use(handlers.AuthMiddleware(r.authService))
		{
			authProtected.POST("/logout", r.authHandler.Logout)
			authProtected.POST("/refresh-token", r.authHandler.RefreshToken)
			authProtected.POST("/change-password", r.authHandler.ChangePassword)
		}

		// Enhanced protected routes that need full user context
		authEnhanced := v1.Group("/auth")
		authEnhanced.Use(handlers.EnhancedAuthMiddleware(r.authService, r.userService))
		{
			// Routes that need full user object in context
			authEnhanced.GET("/me", r.authHandler.Me)
			authEnhanced.GET("/profile", r.authHandler.Me) // Alternative endpoint with full user data
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(handlers.AuthMiddleware(r.authService))
		{
			users.GET("", r.userHandler.GetUsers)
			users.GET("/:id", r.userHandler.GetUserByID)
			users.POST("", r.userHandler.CreateUser)
			users.PUT("/:id", r.userHandler.UpdateUser)
			users.DELETE("/:id", r.userHandler.DeleteUser)
		}

		// Admin routes (require authentication + admin privileges)
		admin := v1.Group("/admin")
		admin.Use(handlers.AuthMiddleware(r.authService))
		admin.Use(handlers.AdminMiddleware())
		{
			// Add admin-only routes here
			admin.GET("/users", r.userHandler.GetUsers) // Admin view of all users
		}
	}

	return engine
}
