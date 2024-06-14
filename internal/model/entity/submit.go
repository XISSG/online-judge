package entity

type Submit struct {
	ID int `json:"id" gorm:"column id;type varchar(256); primaryKey"`
	// "编程语言"
	Language string `json:"language" gorm:"column language; type: varchar(128)"`
	//"用户代码"
	Code string `json:"code" gorm:"column code; type: text; not null"`
	//"判题信息json对象
	JudgeResult string `json:"judge_result" gorm:"column judge_result; type: text;"`
	//判题状态
	Status string `json:"status" gorm:"column status; type: int; default: 0; not null"`
	//"判题id"
	QuestionId int `json:"question_id" gorm:"index; column question_id; type: varchar(256); not null"`
	//"创建用户id"
	UserId int `json:"user_id" gorm:"index; column user_id; type: varchar(256); not null"`
	//"创建时间"
	CreateTime string `json:"create_time" gorm:"column create_time; type: datetime; not null"`
	//"更新时间"
	UpdateTime string `json:"update_time" gorm:"column update_time; type: datetime; not null"`
	//"是否删除",
	IsDelete int8 `json:"is_delete" gorm:"column is_delete; type: int; default: 0; not null"`
}
