package main

import (
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/scheduler"
	"github.com/xissg/online-judge/internal/scheduler/datasync"
	"github.com/xissg/online-judge/internal/scheduler/spider"
	"os"
	"path/filepath"
)

func main() {
	appConfig := config.LoadConfig()
	question := datasync.NewQuestionSync(appConfig.Elasticsearch, appConfig.Database)
	submit := datasync.NewSubmitSync(appConfig.Elasticsearch, appConfig.Database)

	jobs := []scheduler.Job{
		{
			ScheduleTime: "monthly",
			Handler:      CrawlAvatar,
		},
		{

			ScheduleTime: "weekly",
			Handler:      question.SyncData,
		},
		{
			ScheduleTime: "weekly",
			Handler:      submit.SyncData,
		},
	}

	cron := scheduler.NewScheduler()
	cron.Start(jobs)
}

// 爬虫任务
func CrawlAvatar() {
	curPath, _ := os.Getwd()
	path := filepath.Join(curPath, "public/pictures/avatar")
	crawler := spider.NewSpider(path)
	crawler.ConcurrencyCrawl(30, 5)
}
