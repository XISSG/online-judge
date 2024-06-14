package service

import (
	"fmt"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/repository/rabbitmq"
	"testing"
)

func TestRabbitmqService_Publish(t *testing.T) {
	appconfig := config.LoadConfig()
	rabbitMqClient := rabbitmq.NewRabbitMQClient(appconfig.RabbitMQ)
	rabbitmqSvc := NewRabbitMqService(rabbitMqClient)
	rabbitmqSvc.Publish("hello world")
}

func TestRabbitmqService_Consume(t *testing.T) {
	appconfig := config.LoadConfig()
	rabbitMqClient := rabbitmq.NewRabbitMQClient(appconfig.RabbitMQ)
	rabbitmqSvc := NewRabbitMqService(rabbitMqClient)
	rabbitmqSvc.Consume(run)
}

func run(id string) {
	fmt.Println(id)
}
