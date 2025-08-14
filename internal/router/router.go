package router

import (
	"gin-simple-app/internal/handlers"

	"github.com/gin-gonic/gin"
)

// Router holds all the handlers and provides methods to set up routes
type Router struct {
	userHandler   *handlers.UserHandler
	healthHandler *handlers.HealthHandler
}

// NewRouter creates a new router with all handlers
func NewRouter(userHandler *handlers.UserHandler, healthHandler *handlers.HealthHandler) *Router {
	return &Router{
		userHandler:   userHandler,
		healthHandler: healthHandler,
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
		// User routes
		users := v1.Group("/users")
		{
			users.GET("", r.userHandler.GetUsers)
			users.GET("/:id", r.userHandler.GetUserByID)
			users.POST("", r.userHandler.CreateUser)
			users.PUT("/:id", r.userHandler.UpdateUser)
			users.DELETE("/:id", r.userHandler.DeleteUser)
		}
	}

	return engine
}
