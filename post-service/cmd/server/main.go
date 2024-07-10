package main

import (
	"log"
	"post-service/config"
	"post-service/internal/broker"
	"post-service/internal/cache"
	"post-service/internal/db"
	"post-service/internal/handlers"
	"post-service/internal/models"
	"post-service/internal/repositories"
	"post-service/internal/routers"
	"post-service/internal/services"
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
	brokerConn, err := broker.InitRabbitMQ(cfg)
	if err != nil {
		log.Fatalf("Could not connect to the rabbit: %v", err)
	}

	// Initializing postgres DB
	dbConn, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	if err := dbConn.AutoMigrate(
		&models.Post{}, &models.Like{}, &models.Comment{}); err != nil {
		log.Fatalf("Error while trying to run auto migrations: %v", err)
	}

	postRepository := repositories.NewPostRepository(dbConn)
	postService := services.NewPostService(postRepository)
	postHandler := handlers.NewPostHandler(postService, brokerConn)
	router := routers.SetupRouter(postHandler)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server on port: %v, Error: %v", cfg.ServerPort, err)
	}
}
