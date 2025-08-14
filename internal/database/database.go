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
	users := []models.User{
		{Name: "John Doe", Email: "john@example.com"},
		{Name: "Jane Smith", Email: "jane@example.com"},
		{Name: "Bob Johnson", Email: "bob@example.com"},
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
