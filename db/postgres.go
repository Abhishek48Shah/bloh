package db

import (
	"context"

	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Load(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DB.Dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
