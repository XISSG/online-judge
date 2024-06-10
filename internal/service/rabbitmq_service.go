package service

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/xissg/online-judge/internal/repository/rabbitmq"
)

type RabbiMqService interface {
	Publish(message string) error
	Consume() (<-chan amqp091.Delivery, error)
}

type rabbitmqService struct {
	rabbitMqClient *rabbitmq.RabbitMQClient
}

func NewRabbitMqService(rabbitMqClient *rabbitmq.RabbitMQClient) RabbiMqService {
	return &rabbitmqService{
		rabbitMqClient: rabbitMqClient,
	}
}

func (s *rabbitmqService) Publish(message string) error {
	err := s.rabbitMqClient.ExchangeDeclare()
	err = s.rabbitMqClient.Publish(message)
	err = s.rabbitMqClient.Close()
	return err
}

func (s *rabbitmqService) Consume() (<-chan amqp091.Delivery, error) {
	err := s.rabbitMqClient.ExchangeDeclare()
	err = s.rabbitMqClient.QueueDeclareAndBind()
	res, err := s.rabbitMqClient.Consume()
	s.rabbitMqClient.Close()
	return res, err
}
