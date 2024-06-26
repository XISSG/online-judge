package entity

type Question struct {
	ID int `json:"id" gorm:"column id;type varchar(256); primaryKey"`
	// "标题"
	Title string `json:"title" gorm:"column title; type varchar(512)"`
	// "内容"
	Content string `json:"content" gorm:"column content; type text"`
	// "标签列表json数组"
	Tags string `json:"tag" gorm:"column tag; type varchar(1024)"`
	// "判题用例json数组"
	JudgeCase string `json:"judge_case" gorm:"column judge_case; type text"`
	// "判题配置json对象"
	Answer string `json:"answer" gorm:"column answer; type text"`
	// "题目提交数
	JudgeConfig string `json:"judge_config" gorm:"column judge_config; type text"`
	// "题目答案"
	SubmitNum int `json:"submit_num" gorm:"column submit_num; type int; not null;default: 0"`
	// "题目通过数"
	AcceptNum int `json:"accept_num" gorm:"column accept_num; type int; not null;default: 0"`

	// "创建用户id"
	UserId int `json:"user_id" gorm:"index; column user_id;type varchar(256); not null"`
	// "创建时间"
	CreateTime string `json:"create_time" gorm:"column create_time; type datetime"`
	// "更新时间"
	UpdateTime string `json:"update_time" gorm:"column update_time; type datetime"`
	// "是否删除"
	IsDelete int8 `json:"is_delete" gorm:"column is_delete; type int; default: 0"`
}

func (q Question) TableName() string {
	return "question"
}
