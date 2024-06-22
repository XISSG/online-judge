package service

import (
	"github.com/xissg/online-judge/internal/repository/ai"
)

type AIService interface {
	SendMessage(message string) error
	ReceiveMessage() (string, error)
}

type aiService struct {
	client *ai.AIClient
}

func NewAIService(client *ai.AIClient) AIService {
	return &aiService{
		client: client,
	}
}

// 发送消息
func (s *aiService) SendMessage(message string) error {
	return s.client.SendMessage(message)
}

// 获取接收到的消息
func (s *aiService) ReceiveMessage() (string, error) {
	messages, err := s.client.ReadMessage()
	if err != nil {
		return "", err
	}
	res := ""
	for i := range messages {
		if messages[i].Payload == nil {
			continue
		}
		res += messages[i].Payload.Choices.Text[0].Content
	}

	return res, nil
}
