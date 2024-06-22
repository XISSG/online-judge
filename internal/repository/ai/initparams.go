package ai

import "github.com/xissg/online-judge/internal/model/request"

type Option func(ctx *request.AI)

func initAIParams(appid string, option ...Option) *request.AI {
	data := &request.AI{
		Header: &request.Header{
			AppID: appid,
		},
		Parameter: &request.Parameter{
			Chat: &request.Chat{
				Domain:      "generalv3.5",
				Temperature: 0.8,
				TopK:        4,
				MaxTokens:   4096,
			},
		},
		Payload: &request.Payload{
			Message: &request.Message{
				Text: []*request.Text{},
			},
		},
	}

	for _, opt := range option {
		opt(data)
	}
	return data
}

func WithDomain(domain string) Option {
	switch domain {
	case "general", "generalv2", "generalv3,", "generalv3.5":
	default:
		domain = "general"
	}
	return func(ctx *request.AI) {
		ctx.Parameter.Chat.Domain = domain
	}
}

// topk灵活度，1-6默认为4,
func WithTopK(topK int) Option {
	if topK <= 0 || topK > 6 {
		topK = 4
	}
	return func(ctx *request.AI) {
		ctx.Parameter.Chat.TopK = topK
	}
}

// temperature随机性0-1
func WithTemperature(temperature float64) Option {
	if temperature <= 0.0 || temperature > 1 {
		temperature = 0.8
	}
	return func(ctx *request.AI) {
		ctx.Parameter.Chat.Temperature = temperature
	}
}

// 最大token数最大为8192
func WithMaxTokens(maxTokens int) Option {
	if maxTokens <= 0 || maxTokens > 8192 {
		maxTokens = 4096
	}
	return func(ctx *request.AI) {
		ctx.Parameter.Chat.MaxTokens = maxTokens
	}
}
