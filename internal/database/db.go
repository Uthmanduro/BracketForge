package database

import (
	"fmt"
	"log"

	"github.com/uthmanduro/BracketForge/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DBURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(config.DBMaxConns)
	sqlDB.SetMaxIdleConns(config.DBMinConns)

	// Test the database connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")

	return db, nil
}