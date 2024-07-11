package main

import (
	"log"
	"notification-service/config"
	"notification-service/internal/broker"
	"notification-service/internal/cache"
	"notification-service/internal/db"
	"notification-service/internal/handlers"
	"notification-service/internal/models"
	"notification-service/internal/routes"
	"notification-service/internal/services"
)

func main() {
	// Load application configs
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	dbConn, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	if err := dbConn.AutoMigrate(&models.Notification{}); err != nil {
		log.Fatalf("Error while trying to run auto migrations: %v", err)
	}

	// Initializing rabbitmq
	_, err = broker.InitRabbitMQ(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the rabbit: %v", err)
	}

	// Initializing redis
	cache.InitializeRedis(cfg)

	// Initialize SSE broker and handlers
	handlers.InitSSEHandler()

	// message consumers routines
	go services.ConsumeFollowNotificationMessages(cfg, dbConn, handlers.Broker)
	go services.ConsumePostNotificationMessages(cfg, dbConn, handlers.Broker)
	go services.ConsumeLikeNotificationMessages(cfg, dbConn, handlers.Broker)
	go services.ConsumeCommentNotificationMessages(cfg, dbConn, handlers.Broker)

	router := routes.SetupRouter(dbConn)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server on port: %v, Error: %v", cfg.ServerPort, err)
	}
}
