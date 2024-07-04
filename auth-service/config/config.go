package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL   string
	ServerPort    string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	RabbitMQURL   string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// converting redisDB to uint

	return &Config{
		ServerPort: viper.GetString("SERVER_PORT"),
	}, nil
}
