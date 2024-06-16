package response

import (
	"github.com/xissg/online-judge/internal/model/common"
)

type Question struct {
	ID int `json:"id"`
	// "标题"
	Title string `json:"title,omitempty" `
	// "内容"
	Content string `json:"content,omitempty" `
	// "标签列表json数组"
	Tags []string `json:"tags,omitempty" `
	// "题目答案"
	Answer []string `json:"answer,omitempty" `
	// "题目提交数
	SubmitNum int `json:"submit_num,omitempty" `
	// "题目通过数"
	AcceptNum int `json:"accept_num,omitempty"`
	// "判题用例json数组"
	JudgeCase []string `json:"judge_case,omitempty"`
	// "判题配置json对象"
	JudgeConfig *common.Config `json:"judge_config,omitempty" `
	// "创建用户id"
	UserId int `json:"user_id,omitempty"`
	// "创建时间"
	CreateTime string `json:"create_time,omitempty"`
	// "更新时间"
	UpdateTime string `json:"update_time,omitempty"`
}
