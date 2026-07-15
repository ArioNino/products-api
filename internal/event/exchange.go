package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "product_exchange"
	QueueName    = "product_created_queue"
	RoutingKey   = "product.created"
)

func SetupProductBroker(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return ch.QueueBind(
		q.Name,
		RoutingKey,
		ExchangeName,
		false,
		nil,
	)
}