package service

import (
	"github.com/xissg/online-judge/internal/repository/rabbitmq"
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
	return err
}

type HandlerFunc func(string)

func (s *rabbitmqService) Consume(handler HandlerFunc) {
	err := s.rabbitMqClient.ExchangeDeclare()
	err = s.rabbitMqClient.QueueDeclareAndBind()
	msgs, err := s.rabbitMqClient.Consume()
	if err != nil {
		panic(err)
	}
	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			handler(string(msg.Body))
		}
	}()
	<-forever
}
