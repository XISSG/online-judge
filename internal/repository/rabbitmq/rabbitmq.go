package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/xissg/online-judge/internal/constant"
)

func (c *RabbitMQClient) ExchangeDeclare() error {
	err := c.channel.ExchangeDeclare(
		c.ctx.exchangeName,
		c.ctx.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("repository layer: rabbitmq, exchange declare: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (c *RabbitMQClient) Publish(message string) error {
	err := c.channel.Publish(
		c.ctx.exchangeName,
		c.ctx.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return fmt.Errorf("repository layer: rabbitmq, publish: %w %+v", constant.ErrInternal, err)
	}
	return nil
}

func (c *RabbitMQClient) QueueDeclareAndBind() error {
	q, err := c.channel.QueueDeclare(
		c.ctx.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("repository layer: rabbitmq, queue declare: %w %+v", constant.ErrInternal, err)
	}
	err = c.channel.QueueBind(
		q.Name,
		c.ctx.routingKey,
		c.ctx.exchangeName,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("repository layer: rabbitmq, queue bind: %w %+v", constant.ErrInternal, err)
	}
	return nil
}
func (c *RabbitMQClient) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := c.channel.Consume(
		c.ctx.queueName,
		c.ctx.consumerTag,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("repository layer: rabbitmq, consume: %w %+v", constant.ErrInternal, err)
	}
	return msgs, err
}
func (c *RabbitMQClient) Close() error {
	err := c.conn.Close()
	err = c.channel.Close()
	if err != nil {
		return fmt.Errorf("repository layer: rabbitmq,conn and channel close: %w %+v", constant.ErrInternal, err)
	}
	return nil
}
