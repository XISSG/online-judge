package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xissg/online-judge/internal/config"
)

type RabbitMQClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	ctx     *context
}
type context struct {
	exchangeName string
	exchangeType string
	routingKey   string
	queueName    string
	consumerTag  string
}

func NewRabbitMQClient(cfg config.RabbitMQConfig) *RabbitMQClient {
	mqurl := fmt.Sprintf("amqp://%s:%s@%s:%v/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(mqurl)
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return &RabbitMQClient{
		conn:    conn,
		channel: ch,
		ctx: &context{
			exchangeName: cfg.ExchangeName,
			exchangeType: cfg.ExchangeType,
			routingKey:   cfg.RoutingKey,
			queueName:    cfg.QueueName,
			consumerTag:  cfg.ConsumerTag,
		},
	}
}
