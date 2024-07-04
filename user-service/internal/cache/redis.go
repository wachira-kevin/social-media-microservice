package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/user-service/config"
	"sync"
)

var (
	client *redis.Client
	once   sync.Once
)

func InitializeRedis(cfg *config.Config) {
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
	})
}

func GetRedisClient() *redis.Client {
	return client
}

var ctx = context.Background()

func Set(key string, value interface{}) error {
	return client.Set(ctx, key, value, 0).Err()
}

func Get(key string) (string, error) {
	return client.Get(ctx, key).Result()
}
