package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *RabbitMQClient) ExchangeDeclare() error {
	return c.channel.ExchangeDeclare(
		c.ctx.exchangeName,
		c.ctx.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *RabbitMQClient) Publish(message string) error {
	return c.channel.Publish(
		c.ctx.exchangeName,
		c.ctx.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
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
		return err
	}
	err = c.channel.QueueBind(
		q.Name,
		c.ctx.routingKey,
		c.ctx.exchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
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
	return msgs, err
}
func (c *RabbitMQClient) Close() error {
	err := c.conn.Close()
	err = c.channel.Close()
	return err
}
