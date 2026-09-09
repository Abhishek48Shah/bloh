package database

import (
	"context"
	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/redis/go-redis/v9"
)

func Redis(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{Addr: cfg.Redis.Adder, DB: cfg.Redis.DB, Password: cfg.Redis.Password})
	cmd := rdb.Ping(ctx)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return rdb, nil
}
