package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/auto-job-applier/scraper/internal/handler"
	"github.com/auto-job-applier/scraper/internal/models"
	"github.com/auto-job-applier/scraper/internal/scraper/morocco"
	"github.com/auto-job-applier/scraper/internal/scraper/africa"
	"github.com/auto-job-applier/scraper/internal/scraper/gulf"
)

func main() {
	// Initialize scrapers
	scrapers := []interface {
		FetchJobs() ([]models.Job, error)
		GetName() string
	}{
		morocco.NewCareerlinkScraper(),
		morocco.NewRekruteScraper(),
		morocco.NewEmploiMaScraper(),
		africa.NewJobbermanScraper(),
		gulf.NewBaytScraper(),
	}

	// Collect all jobs
	allJobs := []models.Job{}
	for _, s := range scrapers {
		log.Printf("Fetching jobs from %s...", s.GetName())
		jobs, err := s.FetchJobs()
		if err != nil {
			log.Printf("Error fetching from %s: %v", s.GetName(), err)
			continue
		}
		allJobs = append(allJobs, jobs...)
		log.Printf("Found %d jobs from %s", len(jobs), s.GetName())
	}

	// Initialize handler with jobs
	h := handler.NewJobHandler()
	h.SetJobs(allJobs)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", h.Health)
		r.Get("/jobs", h.ListJobs)
		r.Get("/jobs/{id}", h.GetJob)
		r.Get("/sources", h.ListSources)
	})

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Scraper API starting on port %s", port)
	log.Printf("Total jobs loaded: %d", len(allJobs))

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}