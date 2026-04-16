package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"job-applier-backend/internal/models"
	"job-applier-backend/internal/services"
)

// ScrapePayload contains all data needed to trigger a job scrape.
type ScrapePayload struct {
	SearchTerm string   `json:"search_term"`
	Location   string   `json:"location,omitempty"`
	IsRemote   bool     `json:"is_remote,omitempty"`
	Distance   int      `json:"distance,omitempty"`
	Sites      []string `json:"sites,omitempty"`
}

// NewScrapeTask creates an asynq task for job scraping.
func NewScrapeTask(payload ScrapePayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal scrape payload: %w", err)
	}

	return asynq.NewTask(
		TypeScrapeJobs,
		data,
		asynq.MaxRetry(2),
		asynq.Timeout(300*time.Second),
		asynq.Queue("default"),
	), nil
}

// ScrapeHandler processes job scraping tasks.
type ScrapeHandler struct {
	db           *gorm.DB
	pythonClient *services.PythonClient
}

// NewScrapeHandler creates a new ScrapeHandler.
func NewScrapeHandler(db *gorm.DB, pythonClient *services.PythonClient) *ScrapeHandler {
	return &ScrapeHandler{
		db:           db,
		pythonClient: pythonClient,
	}
}

// ProcessTask handles a job scraping task.
func (h *ScrapeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload ScrapePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal scrape payload: %w", err)
	}

	log.Printf("Processing job scrape: search_term=%s, location=%s", payload.SearchTerm, payload.Location)

	// Build request for Python service
	scrapeReq := map[string]interface{}{
		"search_term":    payload.SearchTerm,
		"results_wanted": 50,
		"hours_old":      72,
	}

	if payload.Location != "" {
		scrapeReq["location"] = payload.Location
	}
	if payload.IsRemote {
		scrapeReq["is_remote"] = true
	}
	if payload.Distance > 0 {
		scrapeReq["distance"] = payload.Distance
	}
	if len(payload.Sites) > 0 {
		scrapeReq["sites"] = payload.Sites
	}

	body, err := json.Marshal(scrapeReq)
	if err != nil {
		return fmt.Errorf("marshal scrape request: %w", err)
	}

	// Call Python service /scrape-jobs endpoint
	pythonURL := h.pythonClient.BaseURL() + "/scrape-jobs"
	httpClient := &http.Client{Timeout: 300 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pythonURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create scrape request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call python scrape service: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read scrape response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("python scrape service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var scrapeResult struct {
		Jobs []struct {
			ExternalID  string `json:"external_id"`
			Source      string `json:"source"`
			Title       string `json:"title"`
			Company     string `json:"company"`
			Location    string `json:"location"`
			Description string `json:"description"`
			URL         string `json:"url"`
			Remote      bool   `json:"remote"`
			Salary      string `json:"salary"`
			PostedAt    string `json:"posted_at"`
		} `json:"jobs"`
		Total int `json:"total"`
	}

	if err := json.Unmarshal(respBody, &scrapeResult); err != nil {
		return fmt.Errorf("unmarshal scrape response: %w", err)
	}

	// Ingest results into database
	inserted := 0
	skipped := 0
	for _, j := range scrapeResult.Jobs {
		if j.ExternalID == "" || j.Title == "" {
			skipped++
			continue
		}

		var existing models.Job
		if err := h.db.Where("external_id = ?", j.ExternalID).First(&existing).Error; err == nil {
			skipped++
			continue
		}

		job := models.Job{
			ExternalID:  j.ExternalID,
			Source:      j.Source,
			Title:       j.Title,
			Company:     j.Company,
			Location:    j.Location,
			Description: j.Description,
			URL:         j.URL,
			Remote:      j.Remote,
			Salary:      j.Salary,
		}

		if err := h.db.Create(&job).Error; err != nil {
			skipped++
			continue
		}
		inserted++
	}

	log.Printf("Scrape complete: scraped=%d, inserted=%d, skipped=%d", scrapeResult.Total, inserted, skipped)
	return nil
}
