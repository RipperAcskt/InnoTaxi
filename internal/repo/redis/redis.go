package redis

import (
	"fmt"

	"github.com/RipperAcskt/innotaxi/config"
	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
	cfg    *config.Config
}

func New(cfg *config.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS_DB_HOST,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("client ping failed: %w", err)
	}
	return &Redis{client, cfg}, nil
}
