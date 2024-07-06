package broker

import (
	"fmt"
	"github.com/user-service/config"

	amqp "github.com/rabbitmq/amqp091-go"
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
