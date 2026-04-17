package repository

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"job-applier-backend/internal/cache"
	"job-applier-backend/internal/models"
)

// JobRepository handles database operations for jobs with caching.
type JobRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewJobRepository creates a new JobRepository.
func NewJobRepository(db *gorm.DB, cache cache.Cache) *JobRepository {
	return &JobRepository{db: db, cache: cache}
}

// isPostgres returns true if the underlying database is PostgreSQL.
func (r *JobRepository) isPostgres() bool {
	return r.db.Dialector.Name() == "postgres"
}

// searchCacheKey generates a deterministic cache key for a search query.
func searchCacheKey(query, source string, page int) string {
	raw := fmt.Sprintf("%s:%s:%d", query, source, page)
	hash := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("jobs:search:%x", hash[:8])
}

// searchResult is used for cache serialization of search results.
type searchResult struct {
	Jobs  []models.Job `json:"jobs"`
	Total int64        `json:"total"`
}

// Search finds jobs matching the query with pagination. Uses PostgreSQL full-text
// search when available, falls back to LIKE for SQLite. Results are cached for 5 minutes.
func (r *JobRepository) Search(ctx context.Context, query string, source string, page int) ([]models.Job, int64, error) {
	if page < 1 {
		page = 1
	}
	perPage := 20
	offset := (page - 1) * perPage

	// Check cache
	cacheKey := searchCacheKey(query, source, page)
	if r.cache != nil {
		cached, err := r.cache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var result searchResult
			if err := json.Unmarshal([]byte(cached), &result); err == nil {
				return result.Jobs, result.Total, nil
			}
		}
	}

	db := r.db.WithContext(ctx)
	countDB := r.db.WithContext(ctx).Model(&models.Job{})

	if query != "" {
		// Use LIKE search for both PostgreSQL and SQLite.
		// This is simpler, avoids tsquery syntax errors from user input,
		// and performance is fine for < 100k jobs.
		search := "%" + strings.ToLower(query) + "%"
		db = db.Where("LOWER(title) LIKE ? OR LOWER(company) LIKE ? OR LOWER(description) LIKE ?",
			search, search, search)
		countDB = countDB.Where("LOWER(title) LIKE ? OR LOWER(company) LIKE ? OR LOWER(description) LIKE ?",
			search, search, search)
	}

	if source != "" {
		db = db.Where("source = ?", source)
		countDB = countDB.Where("source = ?", source)
	}

	var total int64
	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count jobs: %w", err)
	}

	var jobs []models.Job
	if err := db.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&jobs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search jobs: %w", err)
	}

	// Cache the result for 5 minutes
	if r.cache != nil {
		result := searchResult{Jobs: jobs, Total: total}
		if data, err := json.Marshal(result); err == nil {
			_ = r.cache.Set(ctx, cacheKey, string(data), 5*time.Minute)
		}
	}

	return jobs, total, nil
}

// FindByID retrieves a job by its primary key.
func (r *JobRepository) FindByID(ctx context.Context, id uint) (*models.Job, error) {
	var job models.Job
	if err := r.db.WithContext(ctx).First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// FindByExternalID retrieves a job by its external_id.
func (r *JobRepository) FindByExternalID(ctx context.Context, externalID string) (*models.Job, error) {
	var job models.Job
	if err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// Upsert inserts a job or skips it if the external_id already exists.
func (r *JobRepository) Upsert(ctx context.Context, job *models.Job) error {
	var existing models.Job
	if err := r.db.WithContext(ctx).Where("external_id = ?", job.ExternalID).First(&existing).Error; err == nil {
		// Already exists, skip
		return nil
	}
	if err := r.db.WithContext(ctx).Create(job).Error; err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	// Invalidate search cache since new jobs were added
	if r.cache != nil {
		_ = r.cache.DeletePattern(ctx, "jobs:search:*")
	}

	return nil
}

// BulkUpsert inserts multiple jobs, skipping those whose external_id already exists.
// Returns the count of inserted and skipped jobs.
func (r *JobRepository) BulkUpsert(ctx context.Context, jobs []models.Job) (inserted int, skipped int, err error) {
	for _, job := range jobs {
		if job.ExternalID == "" || job.Title == "" {
			skipped++
			continue
		}

		var existing models.Job
		if err := r.db.WithContext(ctx).Where("external_id = ?", job.ExternalID).First(&existing).Error; err == nil {
			skipped++
			continue
		}

		if err := r.db.WithContext(ctx).Create(&job).Error; err != nil {
			skipped++
			continue
		}
		inserted++
	}

	// Invalidate search cache if any jobs were inserted
	if inserted > 0 && r.cache != nil {
		_ = r.cache.DeletePattern(ctx, "jobs:search:*")
	}

	return inserted, skipped, nil
}
