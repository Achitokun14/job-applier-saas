package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/glebarez/sqlite"
	"job-applier-backend/internal/models"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	if strings.HasPrefix(databaseURL, "sqlite:") {
		dbPath := strings.TrimPrefix(databaseURL, "sqlite:")
		dialector = sqlite.Open(dbPath)
	} else if strings.HasPrefix(databaseURL, "postgres://") || strings.HasPrefix(databaseURL, "postgresql://") {
		dialector = postgres.Open(databaseURL)
	} else {
		return nil, fmt.Errorf("unsupported database URL: %s", databaseURL)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	return db, nil
}

// AutoMigrate runs GORM auto-migrations for all models.
// GORM AutoMigrate only adds missing columns/tables, never drops data.
func AutoMigrate(db *gorm.DB, appEnv string) error {
	if appEnv == "production" {
		log.Println("INFO: Running AutoMigrate in production (safe: only adds columns/tables)")
	}

	return db.AutoMigrate(
		&models.User{},
		&models.Resume{},
		&models.Job{},
		&models.Application{},
		&models.Settings{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
		&models.Subscription{},
		&models.UsageRecord{},
	)
}
