package response

type Submit struct {
	ID int `json:"id"`
	//"判题信息json对象(包含上面的枚举值)
	JudgeResult []string `json:"judge_info,omitempty"`
	//"判题状态（待判题,判题中,成功,失败)",
	Status string `json:"status,omitempty"`
	//"判题id"
	QuestionId int `json:"question_id,omitempty"`
	//"创建用户id"
	UserId int `json:"user_id,omitempty"`
	//"创建时间"
	CreateTime string `json:"create_time,omitempty"`
	//"更新时间"
	UpdateTime string `json:"update_time,omitempty"`
}
