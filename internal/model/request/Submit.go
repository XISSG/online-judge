package request

type Submit struct {
	// "编程语言"
	Language string `json:"language" gorm:"column language; type: varchar(128)"`
	//"用户代码"
	Code string `json:"code" gorm:"column code; type: text; not null"`
	//"判题id"
	QuestionId int `json:"question_id" gorm:"index; column question_id; type: varchar(256); not null"`
}
