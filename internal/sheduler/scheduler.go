package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"path/filepath"
)

type Scheduler struct {
	crontab *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		crontab: cron.New(),
	}
}

func (s *Scheduler) Start() {
	// 添加定时任务
	_, err := s.crontab.AddFunc("@monthly", s.cronJob)
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	s.crontab.Start()
	log.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.crontab.Stop()
	log.Println("Scheduler stopped")
}

// 定时任务的业务逻辑
func (s *Scheduler) cronJob() {
	log.Println("Running cron job task")
	curPath, _ := os.Getwd()
	path := filepath.Join(curPath, "public/pictures/avatar")
	ConcurrencyCrawl(30, path, 5)
}
