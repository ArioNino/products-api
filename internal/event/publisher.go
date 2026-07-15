package event

import (
	"context"
	"encoding/json"
	"fmt"

	"product-api/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ProductPublisher interface {
	PublishProductCreated(ctx context.Context, product model.Product) error
}

type productPublisher struct {
	channel *amqp.Channel
}

func NewProductPublisher(channel *amqp.Channel) *productPublisher {
	return &productPublisher{channel: channel}
}

func (p *productPublisher) PublishProductCreated(ctx context.Context, product model.Product) error {
	body, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("gagal marshal produk untuk publish: %w", err)
	}

	return p.channel.PublishWithContext(
		ctx,
		ExchangeName,
		RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}