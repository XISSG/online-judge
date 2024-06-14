package request

type Submit struct {
	// "编程语言"
	Language string `json:"language" validate:"max=64,required"`
	//"用户代码"
	Code string `json:"code" validate:"max=1024,required"`
	//"判题id"
	QuestionId int `json:"question_id" validate:"max=64,required"`
}

type UpdateSubmit struct {
	ID int `json:"id" validate:"max=64,required"`
	//"判题信息json对象
	JudgeResult string `json:"judge_result" validate:"max=1024,omitempty"`
	//"判题状态
	Status string `json:"status" validate:"max=64,omitempty"`
}
