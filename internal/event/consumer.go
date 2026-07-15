package event

import (
	"encoding/json"
	"log/slog"

	"product-api/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartProductConsumer(ch *amqp.Channel) {
	msgs, err := ch.Consume(
		QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error("gagal mulai consume", "error", err)
		return
	}

	slog.Info("consumer siap, menunggu pesan...", "queue", QueueName)

	for msg := range msgs {
		var product model.Product
		if err := json.Unmarshal(msg.Body, &product); err != nil {
			slog.Error("gagal parse pesan produk", "error", err)
			msg.Nack(false, false)
			continue
		}

		slog.Info("pesan produk baru diterima", "product_id", product.ID, "name", product.Name)

		msg.Ack(false)
	}
}