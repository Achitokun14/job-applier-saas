package scheduler

import (
	"log"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
	jobs map[string]func()
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
		jobs: make(map[string]func()),
	}
}

func (s *Scheduler) AddJob(sourceName, cronExpr string, fn func()) error {
	_, err := s.cron.AddFunc(cronExpr, fn)
	if err != nil {
		return err
	}
	s.jobs[sourceName] = fn
	log.Printf("Added cron job: %s (%s)", sourceName, cronExpr)
	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) ListJobs() map[string]func() {
	return s.jobs
}