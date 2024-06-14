package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/constant"
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

	//根据不同的发布模式进行相应配置
	ctx := &context{
		exchangeName: cfg.ExchangeName,
		routingKey:   cfg.RoutingKey,
		queueName:    cfg.QueueName,
		consumerTag:  cfg.ConsumerTag,
	}

	switch cfg.PublishType {
	case constant.PUBLISH_PATTERN:
		ctx.routingKey = ""
		ctx.exchangeType = constant.FANOUT_TYPE
	case constant.ROUTING_PATTERN:
		ctx.queueName = ""
		ctx.exchangeType = constant.DIRECT_TYPE
	case constant.TOPIC_PATTERN:
		ctx.queueName = ""
		ctx.exchangeType = constant.TOPIC_TYPE
	default:
		ctx.routingKey = ""
		ctx.exchangeType = constant.FANOUT_TYPE
	}
	return &RabbitMQClient{
		conn:    conn,
		channel: ch,
		ctx:     ctx,
	}
}
