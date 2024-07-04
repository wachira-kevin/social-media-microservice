package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort           string
	KeycloakURL          string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakClientSecret string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// converting redisDB to uint

	return &Config{
		ServerPort:           viper.GetString("SERVER_PORT"),
		KeycloakURL:          viper.GetString("KEYCLOAK_URL"),
		KeycloakRealm:        viper.GetString("KEYCLOAK_REALM"),
		KeycloakClientID:     viper.GetString("KEYCLOAK_CLIENT_ID"),
		KeycloakClientSecret: viper.GetString("KEYCLOAK_CLIENT_SECRET"),
	}, nil
}
