package services

import (
	"time"

	"gorm.io/gorm"

	"job-applier-backend/internal/models"
)

// UsageLimits defines the maximum allowed usage per resource type for a subscription tier.
type UsageLimits struct {
	Applications    int
	ResumeGens      int
	CoverLetterGens int
	InterviewPreps  int
}

// TierLimits maps each subscription tier to its usage limits.
// A value of -1 means unlimited.
var TierLimits = map[string]UsageLimits{
	"free":       {Applications: 5, ResumeGens: 2, CoverLetterGens: 2, InterviewPreps: 0},
	"pro":        {Applications: 50, ResumeGens: -1, CoverLetterGens: -1, InterviewPreps: 5},
	"enterprise": {Applications: -1, ResumeGens: -1, CoverLetterGens: -1, InterviewPreps: -1},
}

// resourceLimitField returns the limit value for a given resource type and tier.
func resourceLimitField(tier string, resourceType string) int {
	limits, ok := TierLimits[tier]
	if !ok {
		limits = TierLimits["free"]
	}

	switch resourceType {
	case "application":
		return limits.Applications
	case "resume_gen":
		return limits.ResumeGens
	case "cover_letter_gen":
		return limits.CoverLetterGens
	case "interview_prep":
		return limits.InterviewPreps
	default:
		return 0
	}
}

// GetUsage returns the current period usage count for a user and resource type.
func GetUsage(db *gorm.DB, userID uint, resourceType string) int {
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	periodEnd := periodStart.AddDate(0, 1, 0)

	var record models.UsageRecord
	err := db.Where(
		"user_id = ? AND resource_type = ? AND period_start = ? AND period_end = ?",
		userID, resourceType, periodStart, periodEnd,
	).First(&record).Error

	if err != nil {
		return 0
	}
	return record.Count
}

// IncrementUsage increments the usage counter for a user and resource type in the current period.
func IncrementUsage(db *gorm.DB, userID uint, resourceType string) error {
	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	periodEnd := periodStart.AddDate(0, 1, 0)

	var record models.UsageRecord
	err := db.Where(
		"user_id = ? AND resource_type = ? AND period_start = ? AND period_end = ?",
		userID, resourceType, periodStart, periodEnd,
	).First(&record).Error

	if err != nil {
		// Record does not exist; create a new one.
		record = models.UsageRecord{
			UserID:       userID,
			ResourceType: resourceType,
			Count:        1,
			PeriodStart:  periodStart,
			PeriodEnd:    periodEnd,
		}
		return db.Create(&record).Error
	}

	// Increment existing record.
	return db.Model(&record).Update("count", record.Count+1).Error
}

// CheckLimit checks whether a user is still under the usage limit for a resource type.
// Returns true if the user can still use the resource, false if the limit is exceeded.
func CheckLimit(db *gorm.DB, userID uint, resourceType string) (bool, error) {
	// Get user's subscription tier.
	tier := getUserTier(db, userID)

	limit := resourceLimitField(tier, resourceType)

	// -1 means unlimited.
	if limit == -1 {
		return true, nil
	}

	// 0 means the feature is not available on this tier.
	if limit == 0 {
		return false, nil
	}

	currentUsage := GetUsage(db, userID, resourceType)
	return currentUsage < limit, nil
}

// GetUserTier returns the subscription tier for a user. Defaults to "free".
func getUserTier(db *gorm.DB, userID uint) string {
	var sub models.Subscription
	err := db.Where("user_id = ? AND status = ?", userID, "active").First(&sub).Error
	if err != nil {
		return "free"
	}
	if sub.Tier == "" {
		return "free"
	}
	return sub.Tier
}

// GetUserTierAndUsage returns the tier, limit, and current usage for a user and resource type.
// This is useful for building detailed error responses.
func GetUserTierAndUsage(db *gorm.DB, userID uint, resourceType string) (tier string, limit int, used int) {
	tier = getUserTier(db, userID)
	limit = resourceLimitField(tier, resourceType)
	used = GetUsage(db, userID, resourceType)
	return
}
