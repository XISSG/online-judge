package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

type Scheduler struct {
	crontab *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		crontab: cron.New(
			cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))),
			cron.WithSeconds()),
	}
}

type handlerFunc func()

// ScheduleTime为cron支持的表达式
type Job struct {
	ScheduleTime string
	Handler      handlerFunc
}

func (s *Scheduler) Start(jobs []Job) {
	// 添加定时任务
	for _, job := range jobs {
		_, err := s.crontab.AddFunc(job.ScheduleTime, job.Handler)
		if err != nil {
			log.Fatalf("Error adding cron job: %v", err)
		}
	}

	s.crontab.Start()
	log.Println("Scheduler started")
	//阻塞执行
	select {}
}

func (s *Scheduler) Stop() {
	s.crontab.Stop()
	log.Println("Scheduler stopped")
}
