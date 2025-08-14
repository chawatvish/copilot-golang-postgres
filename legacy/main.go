package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// User represents a simple user model
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory storage for demo purposes
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	{ID: 3, Name: "Bob Johnson", Email: "bob@example.com"},
}

var nextID = 4

// resetUsers resets the users slice to initial state (for testing)
func resetUsers() {
	users = []User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
		{ID: 3, Name: "Bob Johnson", Email: "bob@example.com"},
	}
	nextID = 4
}

// setupRouter configures and returns the Gin router
func setupRouter() *gin.Engine {
	// Create Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// API routes
	api := r.Group("/api/v1")
	{
		// User routes
		api.GET("/users", getUsers)
		api.GET("/users/:id", getUserByID)
		api.POST("/users", createUser)
		api.PUT("/users/:id", updateUser)
		api.DELETE("/users/:id", deleteUser)
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Gin REST API is running",
		})
	})

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin Simple REST API",
			"version": "1.0.0",
		})
	})

	return r
}

func main() {
	r := setupRouter()
	// Start server on port 8080
	r.Run(":8080")
}

// getUsers returns all users
func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"count": len(users),
	})
}

// getUserByID returns a user by ID
func getUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	for _, user := range users {
		if user.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"data": user,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "User not found",
	})
}

// createUser creates a new user
func createUser(c *gin.Context) {
	var newUser User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Assign new ID
	newUser.ID = nextID
	nextID++

	// Add to users slice
	users = append(users, newUser)

	c.JSON(http.StatusCreated, gin.H{
		"data":    newUser,
		"message": "User created successfully",
	})
}

// updateUser updates an existing user
func updateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var updatedUser User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find and update user
	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = id // Keep the original ID
			users[i] = updatedUser
			c.JSON(http.StatusOK, gin.H{
				"data":    updatedUser,
				"message": "User updated successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "User not found",
	})
}

// deleteUser deletes a user by ID
func deleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Find and delete user
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "User deleted successfully",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "User not found",
	})
}
