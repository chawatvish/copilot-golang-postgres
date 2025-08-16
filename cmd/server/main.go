package main

import (
	"gin-simple-app/internal/config"
	"gin-simple-app/internal/database"
	"gin-simple-app/internal/handlers"
	"gin-simple-app/internal/repository"
	"gin-simple-app/internal/router"
	"gin-simple-app/internal/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize database connection
	if err := database.Connect(&cfg.Database); err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Println("Falling back to in-memory storage...")
		
		// Use in-memory repository as fallback
		runWithInMemoryRepository(cfg)
		return
	}

	// Use database repository
	runWithDatabaseRepository(cfg)
}

func runWithDatabaseRepository(cfg *config.Config) {
	log.Println("Using database repository (GORM + PostgreSQL)")
	
	// Initialize repository with database
	userRepo := repository.NewGormUserRepository(database.GetDB())

	// Initialize services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.GetJWTExpiry())

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	healthHandler := handlers.NewHealthHandler()

	// Initialize router
	appRouter := router.NewRouter(userHandler, authHandler, healthHandler, authService, userService)

	// Setup routes
	engine := appRouter.SetupRoutes()

	// Start server
	port := ":" + cfg.Server.Port
	log.Printf("Starting server on %s (Database mode)", port)
	log.Printf("JWT Secret: %s (first 10 chars)", cfg.JWT.Secret[:min(10, len(cfg.JWT.Secret))])
	log.Printf("JWT Expiry: %v hours", cfg.JWT.ExpireHour)
	
	if err := engine.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	// Cleanup database connection
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()
}

func runWithInMemoryRepository(cfg *config.Config) {
	log.Println("Using in-memory repository (fallback)")
	
	// Initialize repository with in-memory storage
	userRepo := repository.NewInMemoryUserRepository()

	// Initialize services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.GetJWTExpiry())

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	healthHandler := handlers.NewHealthHandler()

	// Initialize router
	appRouter := router.NewRouter(userHandler, authHandler, healthHandler, authService, userService)

	// Setup routes
	engine := appRouter.SetupRoutes()

	// Start server
	port := ":" + cfg.Server.Port
	log.Printf("Starting server on %s (In-memory mode)", port)
	log.Printf("JWT Secret: %s (first 10 chars)", cfg.JWT.Secret[:min(10, len(cfg.JWT.Secret))])
	log.Printf("JWT Expiry: %v hours", cfg.JWT.ExpireHour)
	
	if err := engine.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func init() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// Handle graceful shutdown
	// This is a basic example - in production, you'd want more sophisticated signal handling
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		log.Println("Gin Simple REST API with Authentication")
		log.Println("Environment Variables:")
		log.Println("  DB_HOST        - Database host (default: localhost)")
		log.Println("  DB_PORT        - Database port (default: 5432)")
		log.Println("  DB_USER        - Database user (default: postgres)")
		log.Println("  DB_PASSWORD    - Database password (default: password)")
		log.Println("  DB_NAME        - Database name (default: gin_app)")
		log.Println("  DB_SSLMODE     - SSL mode (default: disable)")
		log.Println("  PORT           - Server port (default: 8080)")
		log.Println("  SERVER_PORT    - Server port (alternative to PORT)")
		log.Println("  GIN_MODE       - Gin mode (default: debug)")
		log.Println("  JWT_SECRET     - JWT secret key (default: your-secret-key-change-this-in-production)")
		log.Println("  JWT_EXPIRE_HOUR - JWT expiry in hours (default: 24)")
		os.Exit(0)
	}
}
