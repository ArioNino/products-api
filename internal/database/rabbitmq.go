package database

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ(dsn string) (*amqp.Connection, error){
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat koneksi RabbitMQ: %w", err)
	}

	return conn, nil
}