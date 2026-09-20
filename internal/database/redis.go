package database

import (
	"context"
	"fmt"

	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/redis/go-redis/v9"
)

func Redis(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	adder := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	rdb := redis.NewClient(&redis.Options{Addr: adder, DB: cfg.Redis.DbName, Password: cfg.Redis.Password})
	cmd := rdb.Ping(ctx)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return rdb, nil
}
