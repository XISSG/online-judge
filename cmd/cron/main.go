package main

import scheduler "github.com/xissg/online-judge/internal/sheduler"

func main() {
	cron := scheduler.NewScheduler()
	cron.Start()
	defer cron.Stop()
}
