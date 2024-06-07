package request

type Submit struct {
	// "编程语言"
	Language string `json:"language"`
	//"用户代码"
	Code string `json:"code"`
	//"判题id"
	QuestionId int `json:"question_id"`
}
