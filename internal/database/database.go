package database

import (
	"gin-simple-app/internal/config"
	"gin-simple-app/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// Connect initializes the database connection
func Connect(cfg *config.DatabaseConfig) error {
	var err error

	// Configure GORM logger
	gormLogger := logger.Default.LogMode(logger.Info)

	// Connect to database
	DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	
	if err != nil {
		return err
	}

	log.Println("Database connection established")

	// Auto-migrate the schema
	if err := AutoMigrate(); err != nil {
		return err
	}

	// Seed initial data
	if err := SeedData(); err != nil {
		return err
	}

	return nil
}

// AutoMigrate runs auto migration for all models
func AutoMigrate() error {
	log.Println("Running database migrations...")
	
	err := DB.AutoMigrate(
		&models.User{},
	)
	
	if err != nil {
		return err
	}
	
	log.Println("Database migrations completed")
	return nil
}

// SeedData seeds initial data into the database
func SeedData() error {
	log.Println("Checking for initial data...")

	// Check if users already exist
	var count int64
	DB.Model(&models.User{}).Count(&count)
	
	if count > 0 {
		log.Println("Data already exists, skipping seed")
		return nil
	}

	// Seed initial users
	address1 := "123 Main St, New York, NY 10001"
	address2 := "456 Oak Ave, Los Angeles, CA 90210"
	phone1 := "+1-555-0101"
	phone2 := "+1-555-0102"
	phone3 := "+1-555-0103"
	
	users := []models.User{
		{Name: "John Doe", Email: "john@example.com", Phone: &phone1, Address: &address1},
		{Name: "Jane Smith", Email: "jane@example.com", Phone: &phone2, Address: &address2},
		{Name: "Bob Johnson", Email: "bob@example.com", Phone: &phone3, Address: nil}, // No address
	}

	result := DB.Create(&users)
	if result.Error != nil {
		return result.Error
	}

	log.Printf("Seeded %d users into database", len(users))
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
