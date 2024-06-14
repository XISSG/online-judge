package request

import (
	"github.com/xissg/online-judge/internal/model/common"
)

type Question struct {
	// "标题"
	Title string `json:"title" validate:"max=512,omitempty"`
	// "内容"
	Content string `json:"content" validate:"max=1024,omitempty"`
	// "标签列表json数组"
	Tag []string `json:"tags" validate:"max=128,omitempty"`
	// "题目答案,数组"
	Answer []string `json:"answer" validate:"max=512,omitempty"`
	// "判题用例json数组,输入用例"
	JudgeCase []string `json:"judge_case" validate:"max=1024,omitempty"`
	// "判题配置json对象,内存限制，时间限制"
	JudgeConfig common.Config `json:"judge_config" validate:"max=64,omitempty"`
}

type UpdateQuestion struct {
	ID int `json:"id" validate:"max=64,required"`
	// "标题"
	Title string `json:"title" validate:"max=512,omitempty"`
	// "内容"
	Content string `json:"content" validate:"max=1024,omitempty"`
	// "标签列表json数组"
	Tag []string `json:"tag" validate:"max=128,omitempty"`
	// "题目答案,数组"
	Answer []string `json:"answer" validate:"max=512,omitempty"`
	// "判题用例json数组,输入用例"
	JudgeCase []string `json:"judge_case" validate:"max=1024,omitempty"`
	// "判题配置json对象,内存限制，时间限制"
	JudgeConfig []common.Config `json:"judge_config" validate:"max=64,omitempty"`
	// "题目提交数
	SubmitNum int `json:"submit_num" validate:"omitempty"`
	// "题目通过数"
	AcceptNum int `json:"accept_num" validate:"omitempty"`
	// "创建用户id"
	UserId int `json:"user_id" validate:"max=64,required"`
}
