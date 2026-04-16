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

// AutoApplyPayload contains data needed to auto-apply to a job in the background.
type AutoApplyPayload struct {
	UserID            uint   `json:"user_id"`
	ApplicationID     uint   `json:"application_id"`
	JobURL            string `json:"job_url"`
	ApplyURL          string `json:"apply_url"`
	Source            string `json:"source"`
	ResumePDFPath     string `json:"resume_pdf_path"`
	CoverLetterPath   string `json:"cover_letter_pdf_path,omitempty"`
	UserName          string `json:"user_name"`
	UserEmail         string `json:"user_email"`
	UserPhone         string `json:"user_phone,omitempty"`
	LinkedInURL       string `json:"linkedin_url,omitempty"`
}

// NewAutoApplyTask creates an asynq task for auto-applying to a job.
func NewAutoApplyTask(payload AutoApplyPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal auto-apply payload: %w", err)
	}

	return asynq.NewTask(
		TypeAutoApply,
		data,
		asynq.MaxRetry(1),
		asynq.Timeout(120*time.Second),
		asynq.Queue("default"),
	), nil
}

// AutoApplyHandler processes auto-apply tasks by calling the Python service.
type AutoApplyHandler struct {
	db           *gorm.DB
	pythonClient *services.PythonClient
}

// NewAutoApplyHandler creates a new AutoApplyHandler.
func NewAutoApplyHandler(db *gorm.DB, pythonClient *services.PythonClient) *AutoApplyHandler {
	return &AutoApplyHandler{
		db:           db,
		pythonClient: pythonClient,
	}
}

// ProcessTask handles an auto-apply task.
func (h *AutoApplyHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload AutoApplyPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal auto-apply payload: %w", err)
	}

	log.Printf("Processing auto-apply for user %d, job URL: %s", payload.UserID, payload.JobURL)

	// Build request body for the Python service /api/auto-apply endpoint
	reqBody := map[string]interface{}{
		"job_url":         payload.JobURL,
		"apply_url":       payload.ApplyURL,
		"source":          payload.Source,
		"resume_pdf_path": payload.ResumePDFPath,
		"user_name":       payload.UserName,
		"user_email":      payload.UserEmail,
	}

	if payload.CoverLetterPath != "" {
		reqBody["cover_letter_pdf_path"] = payload.CoverLetterPath
	}
	if payload.UserPhone != "" {
		reqBody["user_phone"] = payload.UserPhone
	}
	if payload.LinkedInURL != "" {
		reqBody["linkedin_url"] = payload.LinkedInURL
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal auto-apply request: %w", err)
	}

	pythonURL := h.pythonClient.BaseURL() + "/api/auto-apply"
	httpClient := &http.Client{Timeout: 120 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pythonURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create auto-apply request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call python auto-apply service: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read auto-apply response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("python auto-apply service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Success              bool   `json:"success"`
		Method               string `json:"method"`
		ConfirmationID       string `json:"confirmation_id"`
		ScreenshotPath       string `json:"screenshot_path"`
		Error                string `json:"error"`
		RequiresConfirmation bool   `json:"requires_confirmation"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("unmarshal auto-apply response: %w", err)
	}

	// Determine the application status based on the result
	status := "failed"
	if result.RequiresConfirmation {
		status = "ready_to_submit"
	} else if result.Success {
		status = "submitted"
	}

	// Update the Application record
	updates := map[string]interface{}{
		"status":            status,
		"auto_apply_method": result.Method,
	}
	if result.ScreenshotPath != "" {
		updates["screenshot_path"] = result.ScreenshotPath
	}
	if result.ConfirmationID != "" {
		updates["confirmation_id"] = result.ConfirmationID
	}

	if payload.ApplicationID != 0 {
		if err := h.db.Model(&models.Application{}).Where("id = ? AND user_id = ?", payload.ApplicationID, payload.UserID).Updates(updates).Error; err != nil {
			return fmt.Errorf("update application record: %w", err)
		}
	}

	log.Printf(
		"Auto-apply complete for user %d: method=%s, success=%t, requires_confirmation=%t",
		payload.UserID, result.Method, result.Success, result.RequiresConfirmation,
	)

	return nil
}
