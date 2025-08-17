package tests

import (
	"gin-simple-app/internal/handlers"
	"gin-simple-app/internal/models"
	"gin-simple-app/internal/repository"
	"gin-simple-app/internal/router"
	"gin-simple-app/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

// TestApp holds the application components for testing
type TestApp struct {
	router      *gin.Engine
	userRepo    *repository.InMemoryUserRepository
	userService services.UserService
	authService services.AuthService
}

// setupTestApp initializes the application for testing
func setupTestApp() *TestApp {
	gin.SetMode(gin.TestMode)

	// Initialize components with in-memory repository for testing
	userRepo := repository.NewInMemoryUserRepository()
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, "test-secret-key-for-testing", 24*time.Hour)
	
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	healthHandler := handlers.NewHealthHandler()
	
	appRouter := router.NewRouter(userHandler, authHandler, healthHandler, authService, userService)

	return &TestApp{
		router:      appRouter.SetupRoutes(),
		userRepo:    userRepo,
		userService: userService,
		authService: authService,
	}
}

// resetTestData resets the test data before each test
func (app *TestApp) resetTestData() {
	app.userRepo.Reset()
}

// loginUser helper function to get a JWT token for testing
func (app *TestApp) loginUser(email, password string) (string, error) {
	loginReq := models.LoginRequest{
		Email:    email,
		Password: password,
	}
	
	loginResponse, err := app.authService.Login(loginReq)
	if err != nil {
		return "", err
	}
	
	return loginResponse.AccessToken, nil
}
