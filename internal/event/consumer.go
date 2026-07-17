package event

import (
	"encoding/json"
	"log/slog"
	"runtime/debug"

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
		processMsg(msg)
	}
}

func processMsg(msg amqp.Delivery) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic saat proses pesan, recovered",
				"panic", r,
				"stack", string(debug.Stack()),
			)
			msg.Nack(false, false)
		}
	}()

	var product model.Product
	if err := json.Unmarshal(msg.Body, &product); err != nil {
		slog.Error("gagal parse pesan produk", "error", err)
		msg.Nack(false, false)
		return
	}

	slog.Info("pesan produk baru diterima", "product_id", product.ID, "name", product.Name)
	if err := notifyAdminNewProduct(product); err != nil {
		slog.Warn("gagal kirim notifikasi email", "product_id", product.ID, "error", err)
	}

	msg.Ack(false)
}