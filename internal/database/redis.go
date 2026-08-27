package database

import (
	"context"
	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/redis/go-redis/v9"
)

func Redis(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: cfg.Redis.Adder, DB: cfg.Redis.DB, Password: cfg.Redis.Password})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
