package broker

import (
	"fmt"
	"github.com/user-service/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection

func InitRabbitMQ(cfg *config.Config) (*amqp.Connection, error) {

	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("could not connect to RabbitMQ: %w", err)
	}

	RabbitConn = conn
	return RabbitConn, nil
}
