package consumers

import (
	"encoding/json"
	"fmt"
	"github.com/user-service/internal/broker"
	"github.com/user-service/internal/models"
	"github.com/user-service/internal/publishers"
	"github.com/user-service/internal/services"
	"log"
)

func consumeMessages(queueName string, handler func(interface{}) error) {
	conn := broker.RabbitConn
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
			var message interface{}
			switch queueName {
			case "CreatePostCreationNotification":
				var event models.CreatePostNotificationEvent
				if err := json.Unmarshal(d.Body, &event); err != nil {
					log.Println("Failed to deserialize CreatePostNotificationEvent:", err)
					if nackErr := d.Nack(false, false); nackErr != nil {
						log.Println("Failed to nack message:", nackErr)
					}
					continue
				}
				message = event
			case "CreateLikeNotification":
				var event models.CreateLikeNotificationEvent
				if err := json.Unmarshal(d.Body, &event); err != nil {
					log.Println("Failed to deserialize CreateLikeNotificationEvent:", err)
					if nackErr := d.Nack(false, false); nackErr != nil {
						log.Println("Failed to nack message:", nackErr)
					}
					continue
				}
				message = event
			case "CreateCommentNotification":
				var event models.CreateCommentNotificationEvent
				if err := json.Unmarshal(d.Body, &event); err != nil {
					log.Println("Failed to deserialize CreateCommentNotificationEvent:", err)
					if nackErr := d.Nack(false, false); nackErr != nil {
						log.Println("Failed to nack message:", nackErr)
					}
					continue
				}
				message = event
			default:
				log.Printf("Unknown queue name: %s\n", queueName)
				if nackErr := d.Nack(false, false); nackErr != nil {
					log.Println("Failed to nack message:", nackErr)
				}
				continue
			}

			fmt.Printf("Received a message from %s: %+v\n", queueName, message)

			if err := handler(message); err != nil {
				log.Printf("Handler failed for message from %s: %+v\n", queueName, err)
				if nackErr := d.Nack(false, false); nackErr != nil {
					log.Println("Failed to nack message:", nackErr)
				}
				continue
			}

			// Acknowledge the message manually if processing was successful
			if err := d.Ack(false); err != nil {
				log.Println("Failed to acknowledge message:", err)
			}
		}
	}()

	fmt.Printf("Waiting for messages from %s...\n", queueName)
	<-forever
}

func HandleCreatePostNotification(userService services.UserService) {
	consumeMessages("CreatePostCreationNotification", func(message interface{}) error {
		event, ok := message.(models.CreatePostNotificationEvent)
		if !ok {
			return fmt.Errorf("incorrect message type: expected CreatePostNotificationEvent")
		}
		user, err := userService.GetUserByID(event.UserID)
		if err != nil {
			return err
		}
		// terminating notification event to all followers of the user
		followers, err := userService.GetFollowers(event.UserID)
		if err != nil {
			return err
		}
		for _, follower := range followers {
			notificationMessage := fmt.Sprintf("%s has posted a new post", user.Username)
			notificationEvent := getNotificationMessage(&follower, notificationMessage, "post")
			err = publishers.PublishLikeMessage(notificationEvent, broker.RabbitConn)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func HandleCreateLikeNotification(userService services.UserService) {
	consumeMessages("CreateLikeNotification", func(message interface{}) error {
		event, ok := message.(models.CreateLikeNotificationEvent)
		if !ok {
			return fmt.Errorf("incorrect message type: expected CreateLikeNotificationEvent")
		}
		// Implement your logic to handle CreateLikeNotificationEvent
		user, err := userService.GetUserByID(event.UserID)
		if err != nil {
			return err
		}
		liker, err := userService.GetUserByID(event.LikerID)
		if err != nil {
			return err
		}
		notificationMessage := fmt.Sprintf("%s liked your post %v", liker.Username, event.PostID)
		notificationEvent := getNotificationMessage(user, notificationMessage, "like")
		err = publishers.PublishPostMessage(notificationEvent, broker.RabbitConn)
		if err != nil {
			return err
		}
		return nil
	})
}

func HandleCreateCommentNotification(userService services.UserService) {
	consumeMessages("CreateCommentNotification", func(message interface{}) error {
		event, ok := message.(models.CreateCommentNotificationEvent)
		if !ok {
			return fmt.Errorf("incorrect message type: expected CreateCommentNotificationEvent")
		}
		// Implement your logic to handle CreateLikeNotificationEvent
		user, err := userService.GetUserByID(event.UserID)
		if err != nil {
			return err
		}
		commentor, err := userService.GetUserByID(event.CommenterID)
		if err != nil {
			return err
		}
		notificationMessage := fmt.Sprintf("%s commented on your post %v. Comment: %s",
			commentor.Username, event.PostID, event.Comment)
		notificationEvent := getNotificationMessage(user, notificationMessage, "comment")
		err = publishers.PublishCommentMessage(notificationEvent, broker.RabbitConn)
		if err != nil {
			return err
		}
		// Implement your logic to handle CreateCommentNotificationEvent
		fmt.Printf("Handling CreateCommentNotificationEvent: %+v\n", event)
		return nil
	})
}

func getNotificationMessage(user *models.User, message string, messageType string) *models.NotificationSchema {
	log.Println(user)
	// adding follow notification to rabbit to be consumed by notification
	if user.Settings.PushNotifications {
		return &models.NotificationSchema{
			UserID:           user.ID,
			NotificationType: "push",
			MessageType:      messageType,
			Message:          message,
			Email:            user.Email,
		}
	} else {
		return &models.NotificationSchema{
			UserID:           user.ID,
			NotificationType: "email",
			MessageType:      messageType,
			Message:          message,
			Email:            user.Email,
		}
	}
}
