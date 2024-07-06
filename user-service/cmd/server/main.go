package main

import (
	"github.com/user-service/config"
	"github.com/user-service/db"
	"github.com/user-service/internal/broker"
	"github.com/user-service/internal/cache"
	"github.com/user-service/internal/handlers"
	"github.com/user-service/internal/models"
	"github.com/user-service/internal/routes"
	"github.com/user-service/internal/services"
	"log"
)

func main() {
	// Load application configs
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initializing redis
	cache.InitializeRedis(cfg)

	// Initializing rabbitmq
	rabbitConn, err := broker.InitRabbitMQ(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the rabbit: %v", err)
	}

	// Initializing postgres DB
	dbConn, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	if err := dbConn.AutoMigrate(
		&models.User{}, &models.Profile{}, &models.UserFollower{}, &models.UserSettings{}); err != nil {
		log.Fatalf("Error while trying to run auto migrations: %v", err)
	}

	// Initialize services
	userService := services.NewUserService(dbConn, cfg, rabbitConn)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Setup routes
	router := routes.SetupRouter(userHandler)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server on port: %v, Error: %v", cfg.ServerPort, err)
	}
}
