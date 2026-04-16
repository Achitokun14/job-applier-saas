package scraper

import "github.com/auto-job-applier/scraper/internal/models"

type BaseScraper struct {
	Config models.SourceConfig
}

func (b *BaseScraper) GetName() string {
	return b.Config.Name
}

func (b *BaseScraper) GetConfig() models.SourceConfig {
	return b.Config
}

func (b *BaseScraper) FetchJobs() ([]models.Job, error) {
	return []models.Job{}, nil
}

type Scrapable interface {
	GetName() string
	GetConfig() models.SourceConfig
	FetchJobs() ([]models.Job, error)
}

func FetchAll(scrapers []Scrapable) ([]models.Job, error) {
	var allJobs []models.Job
	for _, s := range scrapers {
		jobs, err := s.FetchJobs()
		if err != nil {
			continue
		}
		allJobs = append(allJobs, jobs...)
	}
	return allJobs, nil
}