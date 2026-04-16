package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ResumeRequest is the payload sent to the Python service for resume generation.
type ResumeRequest struct {
	PersonalInfo   string `json:"personal_info"`
	Education      string `json:"education"`
	Experience     string `json:"experience"`
	Skills         string `json:"skills"`
	Projects       string `json:"projects"`
	Achievements   string `json:"achievements"`
	Certifications string `json:"certifications"`
	Languages      string `json:"languages"`
	JobDescription string `json:"job_description"`
	Style          string `json:"style"`
	LLMProvider    string `json:"llm_provider"`
	LLMModel       string `json:"llm_model"`
	LLMAPIKey      string `json:"llm_api_key"`
}

// ResumeResult is the response from the Python service for resume generation.
type ResumeResult struct {
	PDFPath     string `json:"pdf_path"`
	HTMLContent string `json:"html_content"`
	WordCount   int    `json:"word_count"`
}

// CoverLetterRequest is the payload sent to the Python service for cover letter generation.
type CoverLetterRequest struct {
	ResumeText     string `json:"resume_text"`
	JobDescription string `json:"job_description"`
	CompanyName    string `json:"company_name"`
	JobTitle       string `json:"job_title"`
	LLMProvider    string `json:"llm_provider"`
	LLMModel       string `json:"llm_model"`
	LLMAPIKey      string `json:"llm_api_key"`
}

// CoverLetterResult is the response from the Python service for cover letter generation.
type CoverLetterResult struct {
	PDFPath     string `json:"pdf_path"`
	HTMLContent string `json:"html_content"`
	WordCount   int    `json:"word_count"`
}

// PythonClient is an HTTP client for the Python resume/cover letter service.
type PythonClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewPythonClient creates a new PythonClient with a 30-second timeout.
func NewPythonClient(baseURL string) *PythonClient {
	return &PythonClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// BaseURL returns the base URL of the Python service.
func (c *PythonClient) BaseURL() string {
	return c.baseURL
}

// GenerateResume sends a resume generation request to the Python service.
func (c *PythonClient) GenerateResume(ctx context.Context, payload ResumeRequest) (*ResumeResult, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal resume request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/generate-resume", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create resume request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call python service for resume: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read resume response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result ResumeResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal resume response: %w", err)
	}

	return &result, nil
}

// GenerateCoverLetter sends a cover letter generation request to the Python service.
func (c *PythonClient) GenerateCoverLetter(ctx context.Context, payload CoverLetterRequest) (*CoverLetterResult, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal cover letter request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/generate-cover-letter", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create cover letter request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call python service for cover letter: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read cover letter response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("python service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result CoverLetterResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal cover letter response: %w", err)
	}

	return &result, nil
}
