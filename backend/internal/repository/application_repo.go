package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"job-applier-backend/internal/cache"
	"job-applier-backend/internal/models"
)

// ApplicationRepository handles database operations for applications.
type ApplicationRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewApplicationRepository creates a new ApplicationRepository.
func NewApplicationRepository(db *gorm.DB, cache cache.Cache) *ApplicationRepository {
	return &ApplicationRepository{db: db, cache: cache}
}

// ListByUser returns a paginated list of applications for a user, optionally filtered by status.
// Preloads the Job relation.
func (r *ApplicationRepository) ListByUser(ctx context.Context, userID uint, page int, perPage int, status string) ([]models.Application, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	db := r.db.WithContext(ctx).Model(&models.Application{}).Where("user_id = ?", userID)

	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count applications: %w", err)
	}

	var applications []models.Application
	if err := r.db.WithContext(ctx).
		Preload("Job").
		Where("user_id = ?", userID).
		Scopes(func(tx *gorm.DB) *gorm.DB {
			if status != "" {
				return tx.Where("status = ?", status)
			}
			return tx
		}).
		Order("applied_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&applications).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list applications: %w", err)
	}

	return applications, total, nil
}

// Create creates a new application after checking for duplicate (user_id + job_id).
func (r *ApplicationRepository) Create(ctx context.Context, app *models.Application) error {
	var existing models.Application
	if err := r.db.WithContext(ctx).Where("user_id = ? AND job_id = ?", app.UserID, app.JobID).First(&existing).Error; err == nil {
		return fmt.Errorf("duplicate application: user %d already applied to job %d", app.UserID, app.JobID)
	}

	if err := r.db.WithContext(ctx).Create(app).Error; err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of an application owned by a specific user.
func (r *ApplicationRepository) UpdateStatus(ctx context.Context, id uint, userID uint, status string) error {
	result := r.db.WithContext(ctx).
		Model(&models.Application{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update application status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("application not found or not owned by user")
	}
	return nil
}

// Delete removes an application owned by a specific user.
func (r *ApplicationRepository) Delete(ctx context.Context, id uint, userID uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Application{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete application: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("application not found or not owned by user")
	}
	return nil
}
