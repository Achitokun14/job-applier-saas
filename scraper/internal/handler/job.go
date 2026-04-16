package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/auto-job-applier/scraper/internal/models"
)

type JobHandler struct {
	jobs []models.Job
}

func NewJobHandler() *JobHandler {
	return &JobHandler{
		jobs: []models.Job{},
	}
}

func (h *JobHandler) SetJobs(jobs []models.Job) {
	h.jobs = jobs
}

func (h *JobHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *JobHandler) ListJobs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	country := r.URL.Query().Get("country")
	remoteStr := r.URL.Query().Get("remote")
	page := stringsToInt(r.URL.Query().Get("page"), 1)
	limit := stringsToInt(r.URL.Query().Get("limit"), 20)

	jobs := h.jobs

	if query != "" {
		q := strings.ToLower(query)
		filtered := []models.Job{}
		for _, j := range jobs {
			if strings.Contains(strings.ToLower(j.Title), q) ||
			   strings.Contains(strings.ToLower(j.Company), q) ||
			   strings.Contains(strings.ToLower(j.Description), q) {
				filtered = append(filtered, j)
			}
		}
		jobs = filtered
	}

	if country != "" {
		filtered := []models.Job{}
		for _, j := range jobs {
			if j.Country == country {
				filtered = append(filtered, j)
			}
		}
		jobs = filtered
	}

	if remoteStr != "" {
		remote := remoteStr == "true"
		filtered := []models.Job{}
		for _, j := range jobs {
			if j.Remote == remote {
				filtered = append(filtered, j)
			}
		}
		jobs = jobs[:0]
		for _, j := range filtered {
			jobs = append(jobs, j)
		}
	}

	total := len(jobs)
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	resp := models.JobSearchResponse{
		Jobs:     jobs[start:end],
		Total:    total,
		Page:     page,
		Limit:    limit,
		Sources:  []string{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	for _, job := range h.jobs {
		if job.ExternalID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(job)
			return
		}
	}

	http.Error(w, "Job not found", http.StatusNotFound)
}

func (h *JobHandler) ListSources(w http.ResponseWriter, r *http.Request) {
	sources := map[string]bool{}
	for _, job := range h.jobs {
		sources[job.Source] = true
	}

	sourceList := []string{}
	for s := range sources {
		sourceList = append(sourceList, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{"sources": sourceList})
}

func stringsToInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			return defaultVal
		}
	}
	return n
}