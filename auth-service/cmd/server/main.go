package main

import (
	"auth-service/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Load application configs
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	router := gin.Default()

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server on port: %v, Error: %v", cfg.ServerPort, err)
	}
}
