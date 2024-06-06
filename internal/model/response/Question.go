package response

import "time"

type Question struct {
	ID int `json:"id" gorm:"column id;type varchar(256); primaryKey"`
	// "标题"
	Title string `json:"title" gorm:"column title; type varchar(512)"`
	// "内容"
	Content string `json:"content" gorm:"column content; type text"`
	// "标签列表json数组"
	Tag []string `json:"tags" gorm:"column tag; type varchar(1024)"`
	// "题目答案"
	Answer []string `json:"answer" gorm:"column answer; type text"`
	// "题目提交数
	SubmitNum int `json:"submit_num" gorm:"column submit_num; type int; not null;default: 0"`
	// "题目通过数"
	AcceptNum int `json:"accept_num" gorm:"column accept_num; type int; not null;default: 0"`
	// "判题用例json数组"
	JudgeCase []string `json:"judge_case" gorm:"column judge_case; type text"`
	// "判题配置json对象"
	JudgeConfig []Config `json:"judge_config" gorm:"column judge_config; type text"`
	// "点赞数"
	ThumNum int `json:"thum_num" gorm:"column thum_num; type int; not null;default: 0"`
	// "创建用户id"
	UserId int `json:"user_id" gorm:"index; column user_id;type varchar(256); not null"`
	// "创建时间"
	CreateTime string `json:"create_time" gorm:"column create_time; type datetime"`
	// "更新时间"
	UpdateTime string `json:"update_time" gorm:"column update_time; type datetime"`
}

type Config struct {
	TimeLimit   time.Duration `json:"time_limit"`
	MemoryLimit int           `json:"memory_limit"`
}
