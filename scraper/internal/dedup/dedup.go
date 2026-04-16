package dedup

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/auto-job-applier/scraper/internal/models"
)

type Deduplicator struct {
	seen map[string]bool
}

func New() *Deduplicator {
	return &Deduplicator{
		seen: make(map[string]bool),
	}
}

func (d *Deduplicator) GenerateKey(job *models.Job) string {
	normalized := strings.ToLower(job.Title) + "|" + 
		strings.ToLower(job.Company) + "|" + 
		strings.ToLower(job.Location)
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

func (d *Deduplicator) IsDuplicate(job *models.Job) bool {
	key := d.GenerateKey(job)
	if d.seen[key] {
		return true
	}
	d.seen[key] = true
	return false
}

func (d *Deduplicator) Reset() {
	d.seen = make(map[string]bool)
}

func (d *Deduplicator) FilterDuplicate(jobs []models.Job) []models.Job {
	result := []models.Job{}
	for _, job := range jobs {
		if !d.IsDuplicate(&job) {
			result = append(result, job)
		}
	}
	return result
}