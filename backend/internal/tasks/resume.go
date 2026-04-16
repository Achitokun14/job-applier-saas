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

// ResumePayload contains all data needed to generate a resume in the background.
type ResumePayload struct {
	UserID         uint   `json:"user_id"`
	ApplicationID  uint   `json:"application_id,omitempty"`
	Style          string `json:"style"`
	JobDescription string `json:"job_description"`
	LLMProvider    string `json:"llm_provider"`
	LLMModel       string `json:"llm_model"`
	LLMAPIKey      string `json:"llm_api_key"`
	PersonalInfo   string `json:"personal_info"`
	Education      string `json:"education"`
	Experience     string `json:"experience"`
	Skills         string `json:"skills"`
	Projects       string `json:"projects"`
	Achievements   string `json:"achievements"`
	Certifications string `json:"certifications"`
	Languages      string `json:"languages"`
}

// NewResumeTask creates an asynq task for resume generation.
func NewResumeTask(payload ResumePayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal resume payload: %w", err)
	}

	return asynq.NewTask(
		TypeResumeGenerate,
		data,
		asynq.MaxRetry(3),
		asynq.Timeout(120*time.Second),
		asynq.Queue("critical"),
	), nil
}

// ResumeHandler processes resume generation tasks.
type ResumeHandler struct {
	db           *gorm.DB
	pythonClient *services.PythonClient
}

// NewResumeHandler creates a new ResumeHandler.
func NewResumeHandler(db *gorm.DB, pythonClient *services.PythonClient) *ResumeHandler {
	return &ResumeHandler{
		db:           db,
		pythonClient: pythonClient,
	}
}

// ProcessTask handles a resume generation task.
func (h *ResumeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload ResumePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal resume payload: %w", err)
	}

	log.Printf("Processing resume generation for user %d", payload.UserID)

	result, err := h.pythonClient.GenerateResume(ctx, services.ResumeRequest{
		PersonalInfo:   payload.PersonalInfo,
		Education:      payload.Education,
		Experience:     payload.Experience,
		Skills:         payload.Skills,
		Projects:       payload.Projects,
		Achievements:   payload.Achievements,
		Certifications: payload.Certifications,
		Languages:      payload.Languages,
		JobDescription: payload.JobDescription,
		Style:          payload.Style,
		LLMProvider:    payload.LLMProvider,
		LLMModel:       payload.LLMModel,
		LLMAPIKey:      payload.LLMAPIKey,
	})
	if err != nil {
		return fmt.Errorf("python service resume generation failed: %w", err)
	}

	log.Printf("Resume generated for user %d: pdf=%s, words=%d", payload.UserID, result.PDFPath, result.WordCount)

	// Update the user's resume record with the generated PDF path.
	if err := h.db.Model(&models.Resume{}).Where("user_id = ?", payload.UserID).Update("pdf_path", result.PDFPath).Error; err != nil {
		return fmt.Errorf("update resume pdf path: %w", err)
	}

	// If an application ID was provided, update the application's resume PDF field.
	if payload.ApplicationID != 0 {
		if err := h.db.Model(&models.Application{}).Where("id = ? AND user_id = ?", payload.ApplicationID, payload.UserID).Update("resume_pdf", result.PDFPath).Error; err != nil {
			return fmt.Errorf("update application resume pdf: %w", err)
		}
	}

	return nil
}
