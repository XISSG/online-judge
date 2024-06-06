package request

import "time"

type Question struct {
	// "标题"
	Title string `json:"title" gorm:"column title; type varchar(512)"`
	// "内容"
	Content string `json:"content" gorm:"column content; type text"`
	// "标签列表json数组"
	Tag []string `json:"tags" gorm:"column tag; type varchar(1024)"`
	// "题目答案,数组"
	Answer []string `json:"answer" gorm:"column answer; type text"`
	// "判题用例json数组,输入用例"
	JudgeCase []string `json:"judge_case" gorm:"column judge_case; type text"`
	// "判题配置json对象,内存限制，时间限制"
	JudgeConfig []Config `json:"judge_config" gorm:"column judge_config; type text"`
	// "创建用户id"
	UserId int `json:"user_id" gorm:"index; column user_id;type varchar(256); not null"`
}

type Config struct {
	TimeLimit   time.Duration `json:"time_limit"`
	MemoryLimit int           `json:"memory_limit"`
}
