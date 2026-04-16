package models

import "time"

type Job struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ExternalID  string    `gorm:"uniqueIndex;size:255" json:"external_id"`
	Source     string    `gorm:"index;size:50" json:"source"`
	Title      string    `gorm:"index;size:255" json:"title"`
	Company    string    `gorm:"index;size:255" json:"company"`
	Location   string    `gorm:"index" json:"location"`
	Country    string    `gorm:"index;size:50" json:"country"`
	City       string    `gorm:"index;size:100" json:"city"`
	SalaryMin  *int     `json:"salary_min"`
	SalaryMax  *int     `json:"salary_max"`
	Currency   string    `gorm:"size:3" json:"currency"`
	JobType    string    `gorm:"size:50" json:"job_type"`
	Remote    bool      `json:"remote"`
	PostedAt  *time.Time `json:"posted_at"`
	ApplyURL   string    `json:"apply_url"`
	ApplyEmail *string  `json:"apply_email"`
	Description string  `gorm:"type:text" json:"description"`
	Skills    string   `gorm:"size:500" json:"skills"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SourceConfig struct {
	Name       string `json:"name"`
	BaseURL    string `json:"base_url"`
	Enabled   bool   `json:"enabled"`
	Priority  int    `json:"priority"`
	RateLimit int    `json:"rate_limit"`
	Cron      string `json:"cron"`
}

type JobSearchRequest struct {
	Query    string `json:"query"`
	Country string `json:"country"`
	City    string `json:"city"`
	Remote  *bool  `json:"remote"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
}

type JobSearchResponse struct {
	Jobs      []Job  `json:"jobs"`
	Total     int   `json:"total"`
	Page      int   `json:"page"`
	Limit    int   `json:"limit"`
	Sources  []string `json:"sources"`
}