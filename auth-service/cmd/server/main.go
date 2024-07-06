package main

import (
	"auth-service/config"
	"auth-service/internal/handlers"
	"auth-service/internal/routes"
	"auth-service/internal/services"
	"log"
)

func main() {
	// Load application configs
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// keycloak service
	keycloakService := services.NewKeycloakService(cfg)

	// auth requests handler
	authHandler := handlers.NewAuthHandler(keycloakService)

	// setting up routes
	router := routes.SetupRouter(authHandler)

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(cfg.ServerPort); err != nil {
		log.Fatalf("Failed to run server on port: %v, Error: %v", cfg.ServerPort, err)
	}
}
