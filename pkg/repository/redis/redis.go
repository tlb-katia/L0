package redis

import (
	"L0/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type Redis struct {
	Client   *redis.Client
	Duration time.Duration
}

func NewRedis(cfg *config.Config) (*Redis, error) {
	const op = "Storage/Redis.NewRedis"

	redisDuration, err := strconv.Atoi(cfg.RedisDuration)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, op)
	}

	db, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, op)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       db,
	})

	return &Redis{
		Client:   client,
		Duration: time.Duration(redisDuration) * time.Minute,
	}, nil
}
