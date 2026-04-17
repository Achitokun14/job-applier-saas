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
func AutoMigrate(db *gorm.DB, appEnv string, databaseURL string) error {
	if appEnv == "production" {
		log.Println("INFO: Running AutoMigrate in production (safe: only adds columns/tables)")
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Resume{},
		&models.Job{},
		&models.Application{},
		&models.Settings{},
		&models.RefreshToken{},
		&models.PasswordResetToken{},
		&models.Subscription{},
		&models.UsageRecord{},
	); err != nil {
		return err
	}

	// Set up PostgreSQL full-text search trigger and index for jobs
	if strings.HasPrefix(databaseURL, "postgres") {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("WARN: Failed to get underlying sql.DB for FTS setup: %v", err)
			return nil
		}

		// Create trigger function (IF NOT EXISTS via OR REPLACE)
		_, _ = sqlDB.Exec(`
			CREATE OR REPLACE FUNCTION jobs_search_vector_update() RETURNS trigger AS $$
			BEGIN
				NEW.search_vector := to_tsvector('english',
					coalesce(NEW.title, '') || ' ' ||
					coalesce(NEW.company, '') || ' ' ||
					coalesce(NEW.description, '')
				);
				RETURN NEW;
			END
			$$ LANGUAGE plpgsql;
		`)

		// Create trigger (drop first to avoid duplicate)
		_, _ = sqlDB.Exec(`DROP TRIGGER IF EXISTS trg_jobs_search_vector ON jobs`)
		_, _ = sqlDB.Exec(`CREATE TRIGGER trg_jobs_search_vector BEFORE INSERT OR UPDATE ON jobs FOR EACH ROW EXECUTE FUNCTION jobs_search_vector_update()`)

		// Create GIN index
		_, _ = sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_jobs_search_vector ON jobs USING GIN(search_vector)`)

		// Backfill existing rows
		_, _ = sqlDB.Exec(`UPDATE jobs SET search_vector = to_tsvector('english', coalesce(title,'') || ' ' || coalesce(company,'') || ' ' || coalesce(description,'')) WHERE search_vector IS NULL`)
	}

	return nil
}
