package publishers

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/user-service/internal/models"
	"log"
)

func getChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	if conn == nil {
		return nil, fmt.Errorf("RabbitMQ connection not initialized")
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	return ch, nil
}

func declareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return q, fmt.Errorf("failed to declare a queue: %w", err)
	}
	return q, nil
}

func publishMessage(schema *models.NotificationSchema, queueName string, conn *amqp.Connection) error {
	ch, err := getChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareQueue(ch, queueName)
	if err != nil {
		return err
	}

	jsonBody, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to serialize message to JSON: %w", err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	log.Printf("Published %s message for user %d successfully", queueName, schema.UserID)
	return nil
}

func PublishFollowingMessage(schema *models.NotificationSchema, conn *amqp.Connection) error {
	return publishMessage(schema, "SendNewFollowerNotification", conn)
}

func PublishPostMessage(schema *models.NotificationSchema, conn *amqp.Connection) error {
	return publishMessage(schema, "SendNewPostNotification", conn)
}

func PublishLikeMessage(schema *models.NotificationSchema, conn *amqp.Connection) error {
	return publishMessage(schema, "SendNewLikeNotification", conn)
}

func PublishCommentMessage(schema *models.NotificationSchema, conn *amqp.Connection) error {
	return publishMessage(schema, "SendNewCommentNotification", conn)
}
