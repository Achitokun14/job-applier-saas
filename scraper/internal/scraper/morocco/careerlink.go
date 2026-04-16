package morocco

import (
	"github.com/auto-job-applier/scraper/internal/models"
	"github.com/auto-job-applier/scraper/internal/scraper"
	"strings"
	"time"

	"github.com/gocolly/colly"
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
				Cron:      "*/15 * * * *",
			},
		},
	}
}

func (s *CareerlinkScraper) FetchJobs() ([]models.Job, error) {
	jobs := []models.Job{}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		colly.AllowedDomains("careerlink.ma"),
	)

	c.OnHTML("div.job-card, div.listing-job, article.job-listing", func(e *colly.HTMLElement) {
		title := e.ChildText("h3.job-title, h2.job-title, a.job-link")
		company := e.ChildText("span.company-name, .company")
		location := e.ChildText("span.location, .job-location")
		
		if title == "" {
			return
		}

		job := models.Job{
			Source:     s.GetName(),
			ExternalID:  e.Attr("data-job-id"),
			Title:     strings.TrimSpace(title),
			Company:   strings.TrimSpace(company),
			Location:  strings.TrimSpace(location),
			Country:   "MA",
			PostedAt:  func() *time.Time { t := time.Now(); return &t }(),
		}

		applyLink := e.ChildAttr("a.apply-btn, a.job-link", "href")
		if applyLink != "" {
			if !strings.HasPrefix(applyLink, "http") {
				applyLink = s.Config.BaseURL + applyLink
			}
			job.ApplyURL = applyLink
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