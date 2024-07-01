package service

import (
	"fmt"
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
	err := s.client.SendMessage(message)
	if err != nil {
		return fmt.Errorf("service layer: ai -> %w", err)
	}
	return nil
}

// 获取接收到的消息
func (s *aiService) ReceiveMessage() (string, error) {
	messages, err := s.client.ReadMessage()
	if err != nil {
		return "", fmt.Errorf("service layer: ai -> %w", err)
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
