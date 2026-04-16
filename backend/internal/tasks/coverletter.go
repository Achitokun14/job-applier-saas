package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"job-applier-backend/internal/models"
	"job-applier-backend/internal/services"
)

// CoverLetterPayload contains all data needed to generate a cover letter in the background.
type CoverLetterPayload struct {
	UserID         uint   `json:"user_id"`
	ApplicationID  uint   `json:"application_id,omitempty"`
	ResumeText     string `json:"resume_text"`
	JobDescription string `json:"job_description"`
	CompanyName    string `json:"company_name"`
	JobTitle       string `json:"job_title"`
	LLMProvider    string `json:"llm_provider"`
	LLMModel       string `json:"llm_model"`
	LLMAPIKey      string `json:"llm_api_key"`
}

// NewCoverLetterTask creates an asynq task for cover letter generation.
func NewCoverLetterTask(payload CoverLetterPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal cover letter payload: %w", err)
	}

	return asynq.NewTask(
		TypeCoverLetterGenerate,
		data,
		asynq.MaxRetry(3),
		asynq.Timeout(120*time.Second),
		asynq.Queue("critical"),
	), nil
}

// CoverLetterHandler processes cover letter generation tasks.
type CoverLetterHandler struct {
	db           *gorm.DB
	pythonClient *services.PythonClient
}

// NewCoverLetterHandler creates a new CoverLetterHandler.
func NewCoverLetterHandler(db *gorm.DB, pythonClient *services.PythonClient) *CoverLetterHandler {
	return &CoverLetterHandler{
		db:           db,
		pythonClient: pythonClient,
	}
}

// ProcessTask handles a cover letter generation task.
func (h *CoverLetterHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload CoverLetterPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal cover letter payload: %w", err)
	}

	log.Printf("Processing cover letter generation for user %d", payload.UserID)

	result, err := h.pythonClient.GenerateCoverLetter(ctx, services.CoverLetterRequest{
		ResumeText:     payload.ResumeText,
		JobDescription: payload.JobDescription,
		CompanyName:    payload.CompanyName,
		JobTitle:       payload.JobTitle,
		LLMProvider:    payload.LLMProvider,
		LLMModel:       payload.LLMModel,
		LLMAPIKey:      payload.LLMAPIKey,
	})
	if err != nil {
		return fmt.Errorf("python service cover letter generation failed: %w", err)
	}

	log.Printf("Cover letter generated for user %d: pdf=%s, words=%d", payload.UserID, result.PDFPath, result.WordCount)

	// If an application ID was provided, update the application's cover letter PDF field.
	if payload.ApplicationID != 0 {
		if err := h.db.Model(&models.Application{}).Where("id = ? AND user_id = ?", payload.ApplicationID, payload.UserID).Update("cover_pdf", result.PDFPath).Error; err != nil {
			return fmt.Errorf("update application cover pdf: %w", err)
		}
	}

	return nil
}
