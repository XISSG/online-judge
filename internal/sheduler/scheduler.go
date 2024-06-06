package scheduler

import (
	cron "github.com/robfig/cron/v3"
	"github.com/xissg/online-judge/internal/service"
	"log"
)

type Scheduler struct {
	cron        *cron.Cron
	userService service.UserService
}

func NewScheduler(userService service.UserService) *Scheduler {
	return &Scheduler{
		cron:        cron.New(),
		userService: userService,
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
	err := s.userService.CleanupInactiveUsers()
	if err != nil {
		log.Printf("Error running user cleanup task: %v", err)
	}
}
