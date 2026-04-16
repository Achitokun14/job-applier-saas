package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	Email            string         `gorm:"uniqueIndex;not null" json:"email"`
	Password         string         `gorm:"not null" json:"-"`
	Name             string         `json:"name"`
	Role             string         `gorm:"default:'user'" json:"role"`
	EmailVerified    bool           `gorm:"default:false" json:"email_verified"`
	TwoFactorSecret  string         `json:"-"`
	TwoFactorEnabled bool           `gorm:"default:false" json:"two_factor_enabled"`
	Resume           Resume         `gorm:"foreignKey:UserID" json:"resume,omitempty"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}

type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	TokenHash string    `gorm:"uniqueIndex" json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type Resume struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserID        uint      `gorm:"uniqueIndex" json:"user_id"`
	PersonalInfo  string    `gorm:"type:text" json:"personal_info"`
	Education     string    `gorm:"type:text" json:"education"`
	Experience    string    `gorm:"type:text" json:"experience"`
	Skills        string    `gorm:"type:text" json:"skills"`
	Projects      string    `gorm:"type:text" json:"projects"`
	Achievements  string    `gorm:"type:text" json:"achievements"`
	Certifications string   `gorm:"type:text" json:"certifications"`
	Languages     string    `gorm:"type:text" json:"languages"`
	PDFPath       string    `json:"pdf_path"`
}

type Job struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	ExternalID  string    `gorm:"index" json:"external_id"`
	Title       string    `gorm:"index" json:"title"`
	Company     string    `gorm:"index" json:"company"`
	Location    string    `json:"location"`
	Description string    `gorm:"type:text" json:"description"`
	URL         string    `json:"url"`
	Source      string    `gorm:"index" json:"source"`
	Remote      bool      `json:"remote"`
	Salary      string    `json:"salary"`
}

type Application struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UserID          uint      `gorm:"index" json:"user_id"`
	JobID           uint      `gorm:"index" json:"job_id"`
	Job             Job       `gorm:"foreignKey:JobID" json:"job"`
	Status          string    `json:"status"`
	ResumePDF       string    `json:"resume_pdf"`
	CoverPDF        string    `json:"cover_pdf"`
	Notes           string    `gorm:"type:text" json:"notes"`
	AppliedAt       time.Time `json:"applied_at"`
	AutoApplyMethod string    `json:"auto_apply_method,omitempty"`
	ScreenshotPath  string    `json:"screenshot_path,omitempty"`
	ConfirmationID  string    `json:"confirmation_id,omitempty"`
}

type Subscription struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	UserID               uint      `gorm:"uniqueIndex" json:"user_id"`
	StripeCustomerID     string    `gorm:"index" json:"stripe_customer_id"`
	StripeSubscriptionID string    `gorm:"index" json:"stripe_subscription_id"`
	Tier                 string    `gorm:"default:'free'" json:"tier"` // free, pro, enterprise
	Status               string    `gorm:"default:'active'" json:"status"` // active, canceled, past_due
	CurrentPeriodStart   time.Time `json:"current_period_start"`
	CurrentPeriodEnd     time.Time `json:"current_period_end"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type UsageRecord struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	ResourceType string    `json:"resource_type"` // application, resume_gen, cover_letter_gen, api_call
	Count        int       `json:"count"`
	PeriodStart  time.Time `json:"period_start"`
	PeriodEnd    time.Time `json:"period_end"`
	CreatedAt    time.Time `json:"created_at"`
}

type Settings struct {
	ID                uint      `gorm:"primarykey" json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	UserID            uint      `gorm:"uniqueIndex" json:"user_id"`
	LLMProvider       string    `json:"llm_provider"`
	LLMModel          string    `json:"llm_model"`
	LLMAPIKey         string    `json:"llm_api_key,omitempty"`
	JobSearchRemote   bool      `json:"job_search_remote"`
	JobSearchHybrid   bool      `json:"job_search_hybrid"`
	JobSearchOnsite   bool      `json:"job_search_onsite"`
	ExperienceLevel   string    `json:"experience_level"`
	JobTypes          string    `json:"job_types"`
	Positions         string    `json:"positions"`
	Locations         string    `json:"locations"`
	Distance          int       `json:"distance"`
	CompanyBlacklist  string    `json:"company_blacklist"`
	TitleBlacklist    string    `json:"title_blacklist"`
}
