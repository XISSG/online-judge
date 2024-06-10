package scheduler

import (
	cron "github.com/robfig/cron/v3"
	"github.com/xissg/online-judge/internal/service"
	"log"
)

type Scheduler struct {
	cron  *cron.Cron
	judge service.JudgeService
	mq    service.RabbiMqService
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron: cron.New(),
	}
}

func (s *Scheduler) Start() {
	// 添加定时任务
	_, err := s.cron.AddFunc("@every 1h", s.userCleanupTask)
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	s.cron.Start()
	log.Println("Scheduler started")
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// 定时任务的业务逻辑
func (s *Scheduler) userCleanupTask() {
	log.Println("Running user cleanup task")
}
