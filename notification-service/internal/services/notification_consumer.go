package services

import (
	"encoding/json"
	"fmt"
	"log"
	"notification-service/config"
	"notification-service/internal/broker"
	"notification-service/internal/models"
	"notification-service/internal/smtp"
	"notification-service/internal/sse"
	"strconv"

	"gorm.io/gorm"
)

// consumeMessages is a generic function to consume messages from a specified queue and process them.
func consumeMessages(db *gorm.DB, queueName string, handler func(models.NotificationSchema) error) {
	conn := broker.RabbitMQConn
	if conn == nil {
		log.Println("RabbitMQ connection not initialized")
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
		return
	}
	defer func() {
		if err := ch.Close(); err != nil {
			log.Println("Failed to close channel:", err)
		}
	}()

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Println("Failed to declare a queue:", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Println("Failed to register a consumer:", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var message models.NotificationSchema
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Println("Failed to deserialize JSON:", err)
				if nackErr := d.Nack(false, false); nackErr != nil {
					log.Println("Failed to nack message:", nackErr)
				}
				continue
			}

			fmt.Printf("Received a message from %s: %+v\n", queueName, message)

			// Save notification to the database
			notification := models.Notification{
				UserID:           message.UserID,
				MessageType:      message.MessageType,
				NotificationType: message.NotificationType,
				Message:          message.Message,
				Status:           "pending",
			}
			if err := db.Create(&notification).Error; err != nil {
				log.Println("Failed to save notification to database:", err)
				if nackErr := d.Nack(false, false); nackErr != nil {
					log.Println("Failed to nack message:", nackErr)
				}
				continue
			}

			// Process the notification
			if err := handler(message); err != nil {
				log.Println("Failed to process message:", err)
				// Update notification status to failed
				db.Model(&notification).Update("Status", "failed")
				if nackErr := d.Nack(false, false); nackErr != nil {
					log.Println("Failed to nack message:", nackErr)
				}
				continue
			}

			// Update notification status to sent
			db.Model(&notification).Update("Status", "sent")

			// Acknowledge the message manually if processing was successful
			if err := d.Ack(false); err != nil {
				log.Println("Failed to acknowledge message:", err)
			}
		}
	}()

	fmt.Printf("Waiting for messages from %s...\n", queueName)
	<-forever
}

// ConsumeFollowNotificationMessages consumes messages from the "SendNewFollowerNotification" queue.
func ConsumeFollowNotificationMessages(cfg *config.Config, db *gorm.DB, broker *sse.Broker) {
	consumeMessages(db, "SendNewFollowerNotification", func(message models.NotificationSchema) error {
		if message.NotificationType == "email" {
			emailSender := smtp.NewEmailSender(cfg)
			log.Println("sending email...")
			return emailSender.SendEmail(message.Email, "New Follow Notification", message.Message)
		} else if message.NotificationType == "push" {
			return broker.SendMessageToClient(strconv.Itoa(int(message.UserID)), message.Message)
		}
		return nil
	})
}

// ConsumeLikeNotificationMessages consumes messages from the "SendNewLikeNotification" queue.
func ConsumeLikeNotificationMessages(cfg *config.Config, db *gorm.DB, broker *sse.Broker) {
	consumeMessages(db, "SendNewLikeNotification", func(message models.NotificationSchema) error {
		if message.NotificationType == "email" {
			emailSender := smtp.NewEmailSender(cfg)
			return emailSender.SendEmail(message.Email, "New Like Notification", message.Message)
		} else if message.NotificationType == "push" {
			return broker.SendMessageToClient(strconv.Itoa(int(message.UserID)), message.Message)
		}
		return nil
	})
}

// ConsumePostNotificationMessages consumes messages from the "SendNewPostNotification" queue.
func ConsumePostNotificationMessages(cfg *config.Config, db *gorm.DB, broker *sse.Broker) {
	consumeMessages(db, "SendNewPostNotification", func(message models.NotificationSchema) error {
		if message.NotificationType == "email" {
			emailSender := smtp.NewEmailSender(cfg)
			return emailSender.SendEmail(message.Email, "New Post Notification", message.Message)
		} else if message.NotificationType == "push" {
			return broker.SendMessageToClient(strconv.Itoa(int(message.UserID)), message.Message)
		}
		return nil
	})
}

// ConsumeCommentNotificationMessages consumes messages from the "SendNewCommentNotification" queue.
func ConsumeCommentNotificationMessages(cfg *config.Config, db *gorm.DB, broker *sse.Broker) {
	consumeMessages(db, "SendNewCommentNotification", func(message models.NotificationSchema) error {
		if message.NotificationType == "email" {
			emailSender := smtp.NewEmailSender(cfg)
			return emailSender.SendEmail(message.Email, "New Post Notification", message.Message)
		} else if message.NotificationType == "push" {
			return broker.SendMessageToClient(strconv.Itoa(int(message.UserID)), message.Message)
		}
		return nil
	})
}
