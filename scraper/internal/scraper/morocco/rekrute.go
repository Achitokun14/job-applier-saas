package morocco

import (
	"github.com/auto-job-applier/scraper/internal/models"
	"github.com/auto-job-applier/scraper/internal/scraper"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type RekruteScraper struct {
	scraper.BaseScraper
}

func NewRekruteScraper() *RekruteScraper {
	return &RekruteScraper{
		BaseScraper: scraper.BaseScraper{
			Config: models.SourceConfig{
				Name:       "rekrute.com",
				BaseURL:    "https://www.rekrute.com",
				Enabled:   true,
				Priority:  2,
				RateLimit: 30,
				Cron:      "*/30 * * * *",
			},
		},
	}
}

func (s *RekruteScraper) FetchJobs() ([]models.Job, error) {
	jobs := []models.Job{}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		colly.AllowedDomains("rekrute.com"),
	)

	c.OnHTML("div.job-post, article.job, div.offer-item", func(e *colly.HTMLElement) {
		title := e.ChildText("h2.title, h3.title, a.job-title")
		company := e.ChildText("span.company-name, .company, .employer")
		location := e.ChildText("span.location, .city, .place")

		if title == "" {
			return
		}

		job := models.Job{
			Source:     s.GetName(),
			ExternalID: e.Attr("data-id"),
			Title:     strings.TrimSpace(title),
			Company:   strings.TrimSpace(company),
			Location:  strings.TrimSpace(location),
			Country:   "MA",
			PostedAt:  func() *time.Time { t := time.Now(); return &t }(),
		}

		applyLink := e.ChildAttr("a.apply, a.job-link", "href")
		if applyLink != "" && !strings.HasPrefix(applyLink, "http") {
			applyLink = s.Config.BaseURL + applyLink
		}
		if applyLink != "" {
			job.ApplyURL = applyLink
		}

		jobs = append(jobs, job)
	})

	url := s.Config.BaseURL + "/jobs.php"
	if err := c.Visit(url); err != nil {
		return nil, err
	}

	c.Wait()
	return jobs, nil
}