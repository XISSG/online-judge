package service

import (
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/repository/ai"
)

type AIService interface {
}

type aiService struct {
	client  *ai.AIClient
	req     *request.AI
	history []*request.Text
}

// flexibility 1-6, randomness 0-1
func NewAIService(client *ai.AIClient, appId string, roleStr string, flexibility int, randomness float64) *aiService {
	req := client.RoleSetting(flexibility, randomness, appId, roleStr)
	return &aiService{
		client:  client,
		req:     req,
		history: make([]*request.Text, 0),
	}
}
func (s *aiService) SendMessage(message string) error {
	sendMsg := &request.Text{
		Role:    constant.USER_ROLE,
		Content: message,
	}
	s.req.Payload.Message.Text = append(s.req.Payload.Message.Text, sendMsg)
	s.history = append(s.history, s.req.Payload.Message.Text...)
	return s.client.SendMessage(s.req)
}

func (s *aiService) SendMessageWithHistory(message string) error {
	sendMsg := &request.Text{
		Role:    constant.USER_ROLE,
		Content: message,
	}
	s.req.Payload.Message.Text = s.history
	s.req.Payload.Message.Text = append(s.req.Payload.Message.Text, sendMsg)
	s.history = append(s.history, s.req.Payload.Message.Text...)
	return s.client.SendMessage(s.req)
}
func (s *aiService) ReceiveMessage() (string, error) {
	messages, err := s.client.ReadMessage()
	if err != nil {
		return "", err
	}
	var res string
	for i := range messages {
		res += messages[i].Payload.Choices.Text[0].Content
	}

	resMsg := &request.Text{
		Role:    constant.ASSISTANT_ROLE,
		Content: res,
	}

	s.history = append(s.history, resMsg)
	return res, nil
}
