package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xissg/online-judge/internal/config"
)

type RabbitMQClient struct {
	*amqp.Connection
}

func NewRabbitMQClient(cfg config.RabbitMQConfig) *RabbitMQClient {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%v/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		panic(err)
	}
	return &RabbitMQClient{
		conn,
	}
}
