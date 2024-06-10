package rabbitmq

import (
	"fmt"
	"github.com/xissg/online-judge/internal/config"
	"testing"
)

func TestRabbitMQPublisher(t *testing.T) {
	appconfig := config.LoadConfig()
	client := NewRabbitMQClient(appconfig.RabbitMQ)
	err := client.ExchangeDeclare()
	if err != nil {
		panic(err)
	}
	err = client.Publish("hello world")
	if err != nil {
		panic(err)
	}
}

func TestRabbitMQConsumer(t *testing.T) {
	appConfig := config.LoadConfig()
	client := NewRabbitMQClient(appConfig.RabbitMQ)
	err := client.ExchangeDeclare()
	if err != nil {
		panic(err)
	}
	err = client.QueueDeclareAndBind()
	if err != nil {
		panic(err)
	}
	data, err := client.Consume()
	if err != nil {
		panic(err)
	}
	ch := make(chan struct{})
	go func() {
		for d := range data {
			fmt.Println(string(d.Body))
		}
	}()
	<-ch
}
