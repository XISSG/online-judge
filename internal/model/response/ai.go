package response

type AI struct {
	Header  *Header  `json:"header"`
	Payload *Payload `json:"payload"`
}

type Header struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Sid     string `json:"sid"`
	Status  int    `json:"status"`
}
type Payload struct {
	Choices *Choices `json:"choices"`
	Usage   *Usage   `json:"usage"`
}
type Choices struct {
	Status int            `json:"status"`
	Seq    int            `json:"seq"`
	Text   []*PayLoadText `json:"text"`
}
type PayLoadText struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Index   int    `json:"index"`
}

type Usage struct {
	Text *UsageText `json:"text"`
}

// 最后一次结果返回
type UsageText struct {
	QuestionTokens   int `json:"question_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

//# 接口为流式返回，此示例为最后一次返回结果，开发者需要将接口多次返回的结果进行拼接展示
//{
//    "header":{
//        "code":0,
//        "message":"Success",
//        "sid":"cht000cb087@dx18793cd421fb894542",
//        "status":2
//    },
//    "payload":{
//        "choices":{
//            "status":2,
//            "seq":0,
//            "text":[
//                {
//                    "content":"我可以帮助你的吗？",
//                    "role":"assistant",
//                    "index":0
//                }
//            ]
//        },
//        "usage":{
//            "text":{
//                "question_tokens":4,
//                "prompt_tokens":5,
//                "completion_tokens":9,
//                "total_tokens":14
//            }
//        }
//    }
//}
