package database

import (
	"context"
	"time"

	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Postgres(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.Db.DatabaseURL)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(cfg.Db.MaxConn)
	config.MinConns = int32(cfg.Db.MinConn)
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	withTimeoutCtx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	err = pool.Ping(withTimeoutCtx)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
