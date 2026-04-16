# Job Scraping Service Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a job scraping microservice in Go that aggregates jobs from 82+ sources (Morocco → Africa → Gulf → Global → Remote → Tech)

**Architecture:** Go-based scraper with modular source packages, each implementing a standard interface.Scheduler runs cron jobs per source,parser normalizes HTML to Job struct,deduplication prevents cross-source duplicates, REST API exposes to backend.

**Tech Stack:** Go, Colly (scraping), Cron (scheduler), GORM (storage), Chi (API)

---

## Task Structure

### Task 1: Project Setup & Base Interface

**Files:**
- Create: `scraper/Makefile`
- Create: `scraper/go.mod`
- Create: `scraper/internal/models/job.go`
- Create: `scraper/internal/scraper/scraper.go`

**Step 1: Create go.mod with dependencies**

```yaml
module github.com/auto-job-applier/scraper

go 1.21

require (
	github.com/chromedp/chromedp v0.9.5
	github.com/go-chi/chi/v5 v5.0.10
	github.com/gocolly/colly v1.2.0
	github.com/robfig/cron/v3 v3.0.1
	gorm.io/gorm v1.25.5
	gorm.io/driver/sqlite v1.5.4
)
```

**Step 2: Define Job model**

```go
// internal/models/job.go
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
	Currency   string    `size:3" json:"currency"`
	JobType    string    `json:"job_type"` // full-time, part-time, contract
	Remote    bool      `json:"remote"`
	PostedAt  *time.Time`json:"posted_at"`
	ApplyURL   string    `json:"apply_url"`
	ApplyEmail *string  `json:"apply_email"`
	Description string  `type:text" json:"description"`
	Skills    string   `gorm:"size:500" json:"skills"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SourceConfig struct {
	Name       string `json:"name"`
	BaseURL    string `json:"base_url"`
	Enabled   bool   `json:"enabled"`
	Priority  int    `json:"priority"` // 1-10
	RateLimit int    `json:"rate_limit"` // requests per minute
	Cron      string `json:"cron"` // cron expression
}

// IScraper interface that all scrapers must implement
type IScraper interface {
	GetName() string
	GetConfig() SourceConfig
	FetchJobs() ([]Job, error)
}
```

**Step 3: Create base scraper interface**

```go
// internal/scraper/scraper.go
package scraper

import "github.com/auto-job-applier/scraper/internal/models"

type BaseScraper struct {
	Config SourceConfig
}

func (b *BaseScraper) GetName() string {
	return b.Config.Name
}

func (b *BaseScraper) GetConfig() models.SourceConfig {
	return b.Config
}

// FetchJobs returns empty slice - to be implemented by each source
func (b *BaseScraper) FetchJobs() ([]models.Job, error) {
	return []models.Job{}, nil
}
```

**Step 4: Run test to verify it compiles**

Run: `cd C:/Users/X1/job-applier-saas/scraper && go mod tidy`
Expected: SUCCESS

**Step 5: Commit**

```bash
cd C:/Users/X1/job-applier-saas
git add scraper/
git commit -m "feat: add scraper base interface and models"
```

---

### Task 2: Moroccan Job Sources

**Files:**
- Create: `scraper/internal/scraper/morocco/careerlink.go`
- Create: `scraper/internal/scraper/morocco/rekrute.go`
- Create: `scraper/internal/scraper/morocco/emploima.go`
- Create: `scraper/internal/scraper/morocco/dreamjob.go`

**Step 1: Write Careerlink scraper**

```go
// internal/scraper/morocco/careerlink.go
package morocco

import (
	"github.com/auto-job-applier/scraper/internal/models"
	"github.com/auto-job-applier/scraper/internal/scraper"
	"strings"
	"time"
)

type CareerlinkScraper struct {
	scraper.BaseScraper
}

func NewCareerlinkScraper() *CareerlinkScraper {
	return &CareerlinkScraper{
		BaseScraper: scraper.BaseScraper{
			Config: models.SourceConfig{
				Name:       "careerlink.ma",
				BaseURL:    "https://careerlink.ma",
				Enabled:   true,
				Priority:  1,
				RateLimit: 30,
				Cron:      "*/15 * * * *", // every 15 min
			},
		},
	}
}

func (s *CareerlinkScraper) FetchJobs() ([]models.Job, error) {
	jobs := []models.Job{}
	
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowedDomains("careerlink.ma"),
	)
	
	// Find all job cards
	c.OnHTML("div.job-listing-card", func(e *colly.HTMLElement) {
		job := models.Job{
			Source:    s.GetName(),
			ExternalID: e.Attr("data-job-id"),
			Title:    e.ChildText("h3.job-title"),
			Company: e.ChildText("span.company-name"),
			Location: e.ChildText("span.location"),
			PostedAt: parseDate(e.ChildText("span.posted-date")),
		}
		
		// Apply URL
		applyLink := e.ChildAttr("a.apply-btn", "href")
		if applyLink != "" {
			job.ApplyURL = s.Config.BaseURL + applyLink
		}
		
		jobs = append(jobs, job)
	})
	
	url := s.Config.BaseURL + "/en"
	if err := c.Visit(url); err != nil {
		return nil, err
	}
	
	c.Wait()
	return jobs, nil
}

func parseDate(dateStr string) *time.Time {
	// Parse relative dates like "2 days ago", "1 week ago"
	// ... implementation
	return nil
}
```

**Step 2: Create Rekrute scraper**

Similar structure with rekrute.com specifics.

**Step 3: Create Emploi.ma scraper**

Similar structure with emploi.ma specifics.

**Step 4: Run test**

Run: `go build ./...`
Expected: SUCCESS

**Step 5: Commit**

```bash
git add scraper/internal/scraper/morocco/
git commit -m "feat: add Moroccan job scrapers"
```

---

### Task 3: African Job Sources

**Files:**
- Create: `scraper/internal/scraper/africa/jobberman.go`
- Create: `scraper/internal/scraper/africa/brightermonday.go`
- Create: `scraper/internal/scraper/africa/freshtalent.go`

Implement with same pattern for African sources.

---

### Task 4: Gulf / Middle East Sources

**Files:**
- Create: `scraper/internal/scraper/gulf/bayt.go`
- Create: `scraper/internal/scraper/gulf/gulftalent.go`
- Create: `scraper/internal/scraper/gulf/dubizzle.go`

---

### Task 5: Global Aggregators (Indeed, LinkedIn)

**Files:**
- Create: `scraper/internal/scraper/global/indeed.go`
- Create: `scraper/internal/scraper/global/linkedin.go`
- Create: `scraper/internal/scraper/global/jooble.go`

---

### Task 6: Remote-Only Sources

**Files:**
- Create: `scraper/internal/scraper/remote/weworkremotely.go`
- Create: `scraper/internal/scraper/remote/remoteok.go`
- Create: `scraper/internal/scraper/remote/flexjobs.go`

---

### Task 7: Scheduler & Deduplication

**Files:**
- Create: `scraper/internal/scheduler/scheduler.go`
- Create: `scraper/internal/dedup/dedup.go`

**Step 1: Scheduler**

```go
// internal/scheduler/scheduler.go
package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
)

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) AddJob(sourceName, cronExpr string, fn func()) error {
	_, err := s.cron.AddFunc(cronExpr, fn)
	if err != nil {
		return err
	}
	log.Printf("Added cron job: %s (%s)", sourceName, cronExpr)
	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Scheduler stopped")
}
```

**Step 2: Deduplication**

```go
// internal/dedup/dedup.go
package dedup

import (
	"github.com/auto-job-applier/scraper/internal/models"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type Deduplicator struct {
	seen map[string]bool
}

func New() *Deduplicator {
	return &Deduplicator{
		seen: make(map[string]bool),
	}
}

func (d *Deduplicator) GenerateKey(job *models.Job) string {
	// Normalize for consistent hashing
	normalized := strings.ToLower(job.Title) + "|" + 
		strings.ToLower(job.Company) + "|" + 
		strings.ToLower(job.Location)
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

func (d *Deduplicator) IsDuplicate(job *models.Job) bool {
	key := d.GenerateKey(job)
	if d.seen[key] {
		return true
	}
	d.seen[key] = true
	return false
}

func (d *Deduplicator) Reset() {
	d.seen = make(map[string]bool)
}
```

---

### Task 8: REST API

**Files:**
- Create: `scraper/cmd/api/main.go`
- Create: `scraper/internal/handler/job.go`

```go
// cmd/api/main.go
package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"os"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	// Routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler)
		r.Get("/jobs", jobHandler)
		r.Get("/jobs/{id}", jobDetailHandler)
		r.Get("/sources", sourcesHandler)
	})
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	
	log.Printf("API server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"status":"ok"}`))
}
```

---

### Task 9: Docker & Integration

**Files:**
- Create: `scraper/Dockerfile`
- Modify: `docker-compose.yml`

**Step 1: Create Dockerfile**

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scraper ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/scrape r .
EXPOSE 8081
CMD ["./scraper"]
```

**Step 2: Update docker-compose.yml** (shown as diff)

```yaml
  scraper:
    build:
      context: ./scraper
      dockerfile: Dockerfile
    ports:
      - "8081:8081"
    volumes:
      - ./scraper/data:/app/data
    environment:
      - DATABASE_URL=jobs.db
      - PORT=8081
```

---

### Task 10: Testing & Verification

**Step 1: Run unit tests**

Run: `cd scraper && go test ./... -v`

**Step 2: Start service**

Run: `cd scraper && go run cmd/api/main.go`

**Step 3: Verify endpoint**

Run: `curl http://localhost:8081/api/v1/health`

Expected: `{"status":"ok"}`

---

## Execution Option

**Plan complete and saved to `docs/plans/2026-04-13-job-scraping.md`. Two execution options:**

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

Which approach?