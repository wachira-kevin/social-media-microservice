package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"post-service/config"
)

var rabbitConn *amqp.Connection

func InitRabbitMQ(cfg *config.Config) (*amqp.Connection, error) {

	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to RabbitMQ: %w", err)
	}

	rabbitConn = conn
	return rabbitConn, nil
}
