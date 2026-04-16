package africa

import (
	"github.com/auto-job-applier/scraper/internal/models"
	"github.com/auto-job-applier/scraper/internal/scraper"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type JobbermanScraper struct {
	scraper.BaseScraper
}

func NewJobbermanScraper() *JobbermanScraper {
	return &JobbermanScraper{
		BaseScraper: scraper.BaseScraper{
			Config: models.SourceConfig{
				Name:       "jobberman.com",
				BaseURL:    "https://www.jobberman.com",
				Enabled:   true,
				Priority:  1,
				RateLimit: 30,
				Cron:      "*/30 * * * *",
			},
		},
	}
}

func (s *JobbermanScraper) FetchJobs() ([]models.Job, error) {
	jobs := []models.Job{}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowedDomains("jobberman.com", "jobberman.com.gh"),
	)

	c.OnHTML("div.job-card, article.job-listing, div.single-job", func(e *colly.HTMLElement) {
		title := e.ChildText("h3.job-title, h2.title")
		company := e.ChildText("span.company, .company-name")
		location := e.ChildText("span.location, .job-location")
		
		if title == "" {
			return
		}

		job := models.Job{
			Source:     s.GetName(),
			ExternalID: e.Attr("data-job-id"),
			Title:     strings.TrimSpace(title),
			Company:   strings.TrimSpace(company),
			Location:  strings.TrimSpace(location),
			PostedAt:  func() *time.Time { t := time.Now(); return &t }(),
		}

		if strings.Contains(location, "Ghana") {
			job.Country = "GH"
		} else if strings.Contains(location, "Nigeria") {
			job.Country = "NG"
		} else {
			job.Country = "GH"
		}

		applyLink := e.ChildAttr("a.apply-btn, a.job-link", "href")
		if applyLink != "" && !strings.HasPrefix(applyLink, "http") {
			applyLink = s.Config.BaseURL + applyLink
		}
		if applyLink != "" {
			job.ApplyURL = applyLink
		}

		jobs = append(jobs, job)
	})

	url := s.Config.BaseURL + "/gh/jobs"
	if err := c.Visit(url); err != nil {
		c.Visit(s.Config.BaseURL + "/job/search")
	}

	c.Wait()
	return jobs, nil
}