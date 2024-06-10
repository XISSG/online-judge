package service

import (
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/repository/rabbitmq"
	"testing"
)

func TestRabbitmqService_Publish(t *testing.T) {
	appconfig := config.LoadConfig()
	rabbitMqClient := rabbitmq.NewRabbitMQClient(appconfig.RabbitMQ)
	rabbitmqSvc := NewRabbitMqService(rabbitMqClient)
	rabbitmqSvc.Publish("hello world", appconfig.RabbitMQ)
}

func TestRabbitmqService_Consume(t *testing.T) {
	appconfig := config.LoadConfig()
	rabbitMqClient := rabbitmq.NewRabbitMQClient(appconfig.RabbitMQ)
	rabbitmqSvc := NewRabbitMqService(rabbitMqClient)
	rabbitmqSvc.Consume(appconfig.RabbitMQ)
}
