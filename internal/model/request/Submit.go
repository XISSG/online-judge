package request

type Submit struct {
	// "编程语言"
	Language string `json:"language"`
	//"用户代码"
	Code string `json:"code"`
	//"判题id"
	QuestionId int `json:"question_id"`
}

type UpdateSubmit struct {
	ID int `json:"id"`
	//"判题信息json对象
	JudgeResult string `json:"judge_info"`
	//"判题状态（0-待判题,1-判题中,2-成功,3-失败)",
	Status string `json:"status" `
}
