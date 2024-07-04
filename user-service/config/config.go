package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL: viper.GetString("DATABASE_URL"),
		ServerPort:  viper.GetString("SERVER_PORT"),
	}, nil
}
