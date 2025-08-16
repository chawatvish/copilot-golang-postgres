package repository

import (
	"gin-simple-app/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewInMemoryUserRepository(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	assert.NotNil(t, repo)
	assert.IsType(t, &InMemoryUserRepository{}, repo)
}

func TestInMemoryUserRepository_Create(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	user := &models.User{
		Name:    "Test User",
		Email:   "test@example.com",
		Phone:   stringPointer("+1234567890"),
		Address: stringPointer("123 Main St"),
	}
	
	err := repo.Create(user)
	
	assert.NoError(t, err)
	assert.NotEqual(t, uint(0), user.ID) // ID should be set
	assert.False(t, user.CreatedAt.IsZero()) // CreatedAt should be set
	assert.False(t, user.UpdatedAt.IsZero()) // UpdatedAt should be set
}

func TestInMemoryUserRepository_GetAll(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	// Create test users
	user1 := &models.User{Name: "Test1", Email: "test1@example.com"}
	user2 := &models.User{Name: "Test2", Email: "test2@example.com"}
	
	repo.Create(user1)
	repo.Create(user2)
	
	users, err := repo.GetAll()
	
	assert.NoError(t, err)
	assert.Len(t, users, 5) // 3 default users + 2 created users
}

func TestInMemoryUserRepository_GetByID(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	user := &models.User{Name: "Test User", Email: "testuser@example.com"}
	repo.Create(user)
	
	tests := []struct {
		name        string
		id          uint
		expectError bool
	}{
		{
			name:        "existing user",
			id:          user.ID,
			expectError: false,
		},
		{
			name:        "non-existing user",
			id:          999,
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := repo.GetByID(tt.id)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
				assert.Nil(t, foundUser)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, foundUser)
				assert.Equal(t, user.ID, foundUser.ID)
				assert.Equal(t, user.Name, foundUser.Name)
			}
		})
	}
}

func TestInMemoryUserRepository_GetByEmail(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	user := &models.User{Name: "Test User", Email: "testuser@example.com"}
	repo.Create(user)
	
	tests := []struct {
		name        string
		email       string
		expectError bool
	}{
		{
			name:        "existing email",
			email:       "testuser@example.com",
			expectError: false,
		},
		{
			name:        "non-existing email",
			email:       "nonexistent@example.com",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := repo.GetByEmail(tt.email)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
				assert.Nil(t, foundUser)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, foundUser)
				assert.Equal(t, user.Email, foundUser.Email)
			}
		})
	}
}

func TestInMemoryUserRepository_GetByPasswordResetToken(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	token := "reset-token-123"
	user := &models.User{
		Name:               "Test User",
		Email:              "testuser@example.com",
		PasswordResetToken: &token,
	}
	repo.Create(user)
	
	tests := []struct {
		name        string
		token       string
		expectError bool
	}{
		{
			name:        "existing token",
			token:       "reset-token-123",
			expectError: false,
		},
		{
			name:        "non-existing token",
			token:       "nonexistent-token",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundUser, err := repo.GetByPasswordResetToken(tt.token)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
				assert.Nil(t, foundUser)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, foundUser)
				assert.Equal(t, tt.token, *foundUser.PasswordResetToken)
			}
		})
	}
}

func TestInMemoryUserRepository_Update(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	user := &models.User{Name: "Test User", Email: "testuser@example.com"}
	repo.Create(user)
	
	// Update user
	user.Name = "Test User Updated"
	user.Email = "testuserupdated@example.com"
	
	err := repo.Update(user)
	
	assert.NoError(t, err)
	assert.True(t, user.UpdatedAt.After(user.CreatedAt))
	
	// Verify update
	updatedUser, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test User Updated", updatedUser.Name)
	assert.Equal(t, "testuserupdated@example.com", updatedUser.Email)
}

func TestInMemoryUserRepository_Update_NonExistentUser(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	user := &models.User{ID: 999, Name: "Test User", Email: "testuser@example.com"}
	
	err := repo.Update(user)
	
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestInMemoryUserRepository_Delete(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	user := &models.User{Name: "Test User", Email: "testuser@example.com"}
	repo.Create(user)
	
	// Delete user
	err := repo.Delete(user.ID)
	assert.NoError(t, err)
	
	// Verify deletion
	_, err = repo.GetByID(user.ID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestInMemoryUserRepository_Delete_NonExistentUser(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	err := repo.Delete(999)
	
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestInMemoryUserRepository_Count(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	// Get initial count (should include seeded data)
	count, err := repo.Count()
	assert.NoError(t, err)
	initialCount := count
	
	// Add a user
	user := &models.User{Name: "Test User", Email: "testuser@example.com"}
	repo.Create(user)
	
	// Count should increase
	count, err = repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, initialCount+1, count)
	
	// Delete user
	repo.Delete(user.ID)
	
	// Count should decrease
	count, err = repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, initialCount, count)
}

func TestInMemoryUserRepository_Reset(t *testing.T) {
	repo := NewInMemoryUserRepository()
	
	// Add a user
	user := &models.User{Name: "Test User", Email: "testuser@example.com"}
	repo.Create(user)
	
	// Reset should restore initial state
	repo.Reset()
	
	// Should have default users only
	users, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, users, 3) // Only default users
	
	// Added user should not exist
	_, err = repo.GetByID(user.ID)
	assert.Error(t, err)
}

// Helper function for string pointers
func stringPointer(s string) *string {
	return &s
}
