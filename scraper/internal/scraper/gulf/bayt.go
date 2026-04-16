package gulf

import (
	"github.com/auto-job-applier/scraper/internal/models"
	"github.com/auto-job-applier/scraper/internal/scraper"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type BaytScraper struct {
	scraper.BaseScraper
}

func NewBaytScraper() *BaytScraper {
	return &BaytScraper{
		BaseScraper: scraper.BaseScraper{
			Config: models.SourceConfig{
				Name:       "bayt.com",
				BaseURL:    "https://www.bayt.com",
				Enabled:   true,
				Priority:  1,
				RateLimit: 30,
				Cron:      "*/30 * * * *",
			},
		},
	}
}

func (s *BaytScraper) FetchJobs() ([]models.Job, error) {
	jobs := []models.Job{}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		colly.AllowedDomains("bayt.com"),
	)

	c.OnHTML("div.job-card, article.job, div.job-item", func(e *colly.HTMLElement) {
		title := e.ChildText("h3.job-title, h2.title, a.job-title")
		company := e.ChildText("span.company-name, .company")
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

		// Determine country from location
		loc := strings.ToLower(location)
		switch {
		case strings.Contains(loc, "uae") || strings.Contains(loc, "dubai"):
			job.Country = "AE"
		case strings.Contains(loc, "saudi") || strings.Contains(loc, "riyadh"):
			job.Country = "SA"
		case strings.Contains(loc, "qatar") || strings.Contains(loc, "doha"):
			job.Country = "QA"
		case strings.Contains(loc, "kuwait"):
			job.Country = "KW"
		case strings.Contains(loc, "jordan"):
			job.Country = "JO"
		default:
			job.Country = "AE"
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

	if err := c.Visit(s.Config.BaseURL + "/en/jobs"); err != nil {
		return nil, err
	}

	c.Wait()
	return jobs, nil
}