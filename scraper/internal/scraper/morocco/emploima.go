package morocco

import (
	"github.com/auto-job-applier/scraper/internal/models"
	"github.com/auto-job-applier/scraper/internal/scraper"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type EmploiMaScraper struct {
	scraper.BaseScraper
}

func NewEmploiMaScraper() *EmploiMaScraper {
	return &EmploiMaScraper{
		BaseScraper: scraper.BaseScraper{
			Config: models.SourceConfig{
				Name:       "emploi.ma",
				BaseURL:    "https://www.emploi.ma",
				Enabled:   true,
				Priority:  3,
				RateLimit: 30,
				Cron:      "*/30 * * * *",
			},
		},
	}
}

func (s *EmploiMaScraper) FetchJobs() ([]models.Job, error) {
	jobs := []models.Job{}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		colly.AllowedDomains("emploi.ma"),
	)

	c.OnHTML("div.job-listing, article.job, div.offer", func(e *colly.HTMLElement) {
		title := e.ChildText("h2.title, h3.title, a.position")
		company := e.ChildText("span.company, .company-name")
		location := e.ChildText("span.location, .city")

		if title == "" {
			return
		}

		job := models.Job{
			Source:     s.GetName(),
			ExternalID: e.Attr("id"),
			Title:     strings.TrimSpace(title),
			Company:   strings.TrimSpace(company),
			Location:  strings.TrimSpace(location),
			Country:   "MA",
			PostedAt:  func() *time.Time { t := time.Now(); return &t }(),
		}

		applyLink := e.ChildAttr("a.position, a.apply", "href")
		if applyLink != "" && !strings.HasPrefix(applyLink, "http") {
			applyLink = s.Config.BaseURL + applyLink
		}
		if applyLink != "" {
			job.ApplyURL = applyLink
		}

		jobs = append(jobs, job)
	})

	if err := c.Visit(s.Config.BaseURL + "/offres-emploi-maroc"); err != nil {
		return nil, err
	}

	c.Wait()
	return jobs, nil
}