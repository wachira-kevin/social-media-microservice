package publishers

import (
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/user-service/internal/models"
	"log"
)

func PublishFollowingMessage(schema *models.NotificationSchema, conn *amqp.Connection) error {
	// check connection
	if conn == nil {
		return fmt.Errorf("RabbitMQ connection not initialized")
	}
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()
	// queue declaration
	q, err := ch.QueueDeclare(
		"new_follower", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Serialize message to JSON
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

	log.Printf("publisher follow message for user %d successfully", schema.UserID)

	return nil
}
