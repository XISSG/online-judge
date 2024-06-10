package request

import (
	"github.com/xissg/online-judge/internal/model/common"
)

type Question struct {
	// "标题"
	Title string `json:"title"`
	// "内容"
	Content string `json:"content"`
	// "标签列表json数组"
	Tag []string `json:"tags"`
	// "题目答案,数组"
	Answer []string `json:"answer"`
	// "判题用例json数组,输入用例"
	JudgeCase []string `json:"judge_case"`
	// "判题配置json对象,内存限制，时间限制"
	JudgeConfig common.Config `json:"judge_config"`
	// "创建用户id"
	UserId int `json:"user_id"`
}

type UpdateQuestion struct {
	ID int `json:"id"`
	// "标题"
	Title string `json:"title"`
	// "内容"
	Content string `json:"content"`
	// "标签列表json数组"
	Tag []string `json:"tags"`
	// "题目答案,数组"
	Answer []string `json:"answer"`
	// "判题用例json数组,输入用例"
	JudgeCase []string `json:"judge_case"`
	// "判题配置json对象,内存限制，时间限制"
	JudgeConfig []common.Config `json:"judge_config"`
	// "题目提交数
	SubmitNum int `json:"submit_num"`
	// "题目通过数"
	AcceptNum int `json:"accept_num"`
	// "创建用户id"
	UserId int `json:"user_id"`
}
