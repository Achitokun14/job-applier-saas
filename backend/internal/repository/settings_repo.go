package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"job-applier-backend/internal/cache"
	"job-applier-backend/internal/models"
)

// SettingsRepository handles database operations for user settings with caching.
type SettingsRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewSettingsRepository creates a new SettingsRepository.
func NewSettingsRepository(db *gorm.DB, cache cache.Cache) *SettingsRepository {
	return &SettingsRepository{db: db, cache: cache}
}

// settingsCacheKey returns the cache key for settings by user ID.
func settingsCacheKey(userID uint) string {
	return fmt.Sprintf("settings:%d", userID)
}

// GetByUserID retrieves settings for a user. Results are cached for 30 minutes.
func (r *SettingsRepository) GetByUserID(ctx context.Context, userID uint) (*models.Settings, error) {
	// Check cache
	cacheKey := settingsCacheKey(userID)
	if r.cache != nil {
		cached, err := r.cache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var settings models.Settings
			if err := json.Unmarshal([]byte(cached), &settings); err == nil {
				return &settings, nil
			}
		}
	}

	var settings models.Settings
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error; err != nil {
		return nil, err
	}

	// Cache the result for 30 minutes
	if r.cache != nil {
		if data, err := json.Marshal(settings); err == nil {
			_ = r.cache.Set(ctx, cacheKey, string(data), 30*time.Minute)
		}
	}

	return &settings, nil
}

// Upsert creates or updates settings for a user with write-through cache invalidation.
func (r *SettingsRepository) Upsert(ctx context.Context, settings *models.Settings) error {
	var existing models.Settings
	result := r.db.WithContext(ctx).Where("user_id = ?", settings.UserID).First(&existing)

	if result.Error != nil {
		// Does not exist, create
		if err := r.db.WithContext(ctx).Create(settings).Error; err != nil {
			return fmt.Errorf("failed to create settings: %w", err)
		}
	} else {
		// Exists, update
		settings.ID = existing.ID
		if err := r.db.WithContext(ctx).Save(settings).Error; err != nil {
			return fmt.Errorf("failed to update settings: %w", err)
		}
	}

	// Invalidate cache (write-through)
	if r.cache != nil {
		_ = r.cache.Delete(ctx, settingsCacheKey(settings.UserID))
	}

	return nil
}
