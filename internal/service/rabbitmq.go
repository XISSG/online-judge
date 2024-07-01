package service

import (
	"fmt"
	"github.com/xissg/online-judge/internal/repository/rabbitmq"
	"log"
)

type RabbitMqService interface {
	Publish(message string) error
	Consume(handler HandlerFunc)
}

type rabbitmqService struct {
	rabbitMqClient *rabbitmq.RabbitMQClient
}

func NewRabbitMqService(rabbitMqClient *rabbitmq.RabbitMQClient) RabbitMqService {

	return &rabbitmqService{
		rabbitMqClient: rabbitMqClient,
	}
}

func (s *rabbitmqService) Publish(message string) error {
	err := s.rabbitMqClient.ExchangeDeclare()
	err = s.rabbitMqClient.Publish(message)
	defer s.rabbitMqClient.Close()
	if err != nil {
		return fmt.Errorf("service layer: rabbitmq -> %w", err)
	}
	return nil
}

type HandlerFunc func(string)

func (s *rabbitmqService) Consume(handler HandlerFunc) {
	err := s.rabbitMqClient.ExchangeDeclare()
	err = s.rabbitMqClient.QueueDeclareAndBind()
	msgs, err := s.rabbitMqClient.Consume()
	if err != nil {
		log.Println(err)
	}
	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			handler(string(msg.Body))
		}
	}()
	<-forever
}
