package request

type AI struct {
	Header    *Header    `json:"header,omitempty"`
	Parameter *Parameter `json:"parameter,omitempty"`
	Payload   *Payload   `json:"payload,omitempty"`
}

type Header struct {
	AppID string `json:"app_id"`
	Uid   string `json:"uid,omitempty"`
}
type Parameter struct {
	Chat *Chat `json:"chat,omitempty"`
}

type Chat struct {
	Domain      string  `json:"domain,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
}

type Payload struct {
	Message *Message `json:"message,omitempty"`
}

type Message struct {
	Text []*Text `json:"text,omitempty"`
}
type Text struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

//# 参数构造示例如下
//{
//        "header": {
//            "app_id": "12345",
//            "uid": "12345"
//        },
//        "parameter": {
//            "chat": {
//                "domain": "generalv3.5",
//                "temperature": 0.5,
//                "max_tokens": 1024,
//            }
//        },
//        "payload": {
//            "message": {
//                # 如果想获取结合上下文的回答，需要开发者每次将历史问答信息一起传给服务端，如下示例
//                # 注意：text里面的所有content内容加一起的tokens需要控制在8192以内，开发者如有较长对话需求，需要适当裁剪历史信息
//                "text": [
//                    {"role":"system","content":"你现在扮演李白，你豪情万丈，狂放不羁；接下来请用李白的口吻和用户对话。"} #设置对话背景或者模型角色
//                    {"role": "user", "content": "你是谁"} # 用户的历史问题
//                    {"role": "assistant", "content": "....."}  # AI的历史回答结果
//                    # ....... 省略的历史对话
//                    {"role": "user", "content": "你会做什么"}  # 最新的一条问题，如无需上下文，可只传最新一条问题
//                ]
//        }
//    }
//}
