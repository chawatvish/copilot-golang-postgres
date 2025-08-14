package repository

import (
	"errors"
	"gin-simple-app/internal/models"
	"sync"
	"time"

	"gorm.io/gorm"
)

// InMemoryUserRepository implements UserRepository using in-memory storage (for testing)
type InMemoryUserRepository struct {
	users  []models.User
	nextID uint
	mutex  sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository with sample data
func NewInMemoryUserRepository() *InMemoryUserRepository {
	now := time.Now()
	address1 := "123 Main St, New York, NY 10001"
	address2 := "456 Oak Ave, Los Angeles, CA 90210"
	phone1 := "+1-555-0101"
	phone2 := "+1-555-0102"
	phone3 := "+1-555-0103"
	
	return &InMemoryUserRepository{
		users: []models.User{
			{ID: 1, Name: "John Doe", Email: "john@example.com", Phone: &phone1, Address: &address1, CreatedAt: now, UpdatedAt: now},
			{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Phone: &phone2, Address: &address2, CreatedAt: now, UpdatedAt: now},
			{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", Phone: &phone3, Address: nil, CreatedAt: now, UpdatedAt: now},
		},
		nextID: 4,
	}
}

// GetAll returns all users
func (r *InMemoryUserRepository) GetAll() ([]models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	// Return a copy to prevent external modification
	usersCopy := make([]models.User, 0, len(r.users))
	for _, user := range r.users {
		if user.DeletedAt.Time.IsZero() { // Only non-deleted users
			usersCopy = append(usersCopy, user)
		}
	}
	return usersCopy, nil
}

// GetByID returns a user by ID
func (r *InMemoryUserRepository) GetByID(id uint) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, user := range r.users {
		if user.ID == id && user.DeletedAt.Time.IsZero() {
			userCopy := user
			return &userCopy, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

// GetByEmail returns a user by email
func (r *InMemoryUserRepository) GetByEmail(email string) (*models.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, user := range r.users {
		if user.Email == email && user.DeletedAt.Time.IsZero() {
			userCopy := user
			return &userCopy, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check for duplicate email
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email && existingUser.DeletedAt.Time.IsZero() {
			return errors.New("email already exists")
		}
	}
	
	now := time.Now()
	user.ID = r.nextID
	user.CreatedAt = now
	user.UpdatedAt = now
	r.nextID++
	r.users = append(r.users, *user)
	return nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(user *models.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	for i, existingUser := range r.users {
		if existingUser.ID == user.ID && existingUser.DeletedAt.Time.IsZero() {
			// Check for duplicate email (excluding current user)
			for _, otherUser := range r.users {
				if otherUser.Email == user.Email && otherUser.ID != user.ID && otherUser.DeletedAt.Time.IsZero() {
					return errors.New("email already exists")
				}
			}
			
			user.UpdatedAt = time.Now()
			user.CreatedAt = existingUser.CreatedAt // Preserve creation time
			r.users[i] = *user
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

// Delete deletes a user by ID (soft delete)
func (r *InMemoryUserRepository) Delete(id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	for i, user := range r.users {
		if user.ID == id && user.DeletedAt.Time.IsZero() {
			r.users[i].DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

// Count returns the total number of non-deleted users
func (r *InMemoryUserRepository) Count() (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	var count int64
	for _, user := range r.users {
		if user.DeletedAt.Time.IsZero() {
			count++
		}
	}
	return count, nil
}

// Reset resets the repository to initial state (for testing)
func (r *InMemoryUserRepository) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	now := time.Now()
	address1 := "123 Main St, New York, NY 10001"
	address2 := "456 Oak Ave, Los Angeles, CA 90210"
	phone1 := "+1-555-0101"
	phone2 := "+1-555-0102"
	phone3 := "+1-555-0103"
	
	r.users = []models.User{
		{ID: 1, Name: "John Doe", Email: "john@example.com", Phone: &phone1, Address: &address1, CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Phone: &phone2, Address: &address2, CreatedAt: now, UpdatedAt: now},
		{ID: 3, Name: "Bob Johnson", Email: "bob@example.com", Phone: &phone3, Address: nil, CreatedAt: now, UpdatedAt: now},
	}
	r.nextID = 4
}
