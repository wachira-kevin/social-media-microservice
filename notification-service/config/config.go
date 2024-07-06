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
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL:   viper.GetString("DATABASE_URL"),
		ServerPort:    viper.GetString("SERVER_PORT"),
		RedisAddr:     viper.GetString("REDIS_ADDR"),
		RedisPassword: viper.GetString("REDIS_PASSWORD"),
		RedisDB:       viper.GetInt("REDIS_DB"),
		RabbitMQURL:   viper.GetString("RABBITMQ_URL"),
		SMTPHost:      viper.GetString("SMTP_HOST"),
		SMTPPort:      viper.GetInt("SMTP_PORT"),
		SMTPUsername:  viper.GetString("SMTP_USERNAME"),
		SMTPPassword:  viper.GetString("SMTP_PASSWORD"),
	}, nil
}
