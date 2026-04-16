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

// UserRepository handles database operations for users with caching.
type UserRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB, cache cache.Cache) *UserRepository {
	return &UserRepository{db: db, cache: cache}
}

// userCacheKey returns the cache key for a user by ID.
func userCacheKey(id uint) string {
	return fmt.Sprintf("user:%d", id)
}

// FindByID retrieves a user by ID. Results are cached for 10 minutes.
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	// Check cache
	cacheKey := userCacheKey(id)
	if r.cache != nil {
		cached, err := r.cache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var user models.User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				return &user, nil
			}
		}
	}

	var user models.User
	if err := r.db.WithContext(ctx).Preload("Resume").First(&user, id).Error; err != nil {
		return nil, err
	}

	// Cache the result for 10 minutes
	if r.cache != nil {
		if data, err := json.Marshal(user); err == nil {
			_ = r.cache.Set(ctx, cacheKey, string(data), 10*time.Minute)
		}
	}

	return &user, nil
}

// FindByEmail retrieves a user by email address.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create inserts a new user into the database.
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Update saves changes to a user and invalidates the cache.
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Invalidate cache
	if r.cache != nil {
		_ = r.cache.Delete(ctx, userCacheKey(user.ID))
	}

	return nil
}
