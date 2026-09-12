package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Abhishek48Shah/blog/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Postgres(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	//postgres:://user/password@host:port/dbname?mode
	URL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.DbName, cfg.Db.SSLMode)
	config, err := pgxpool.ParseConfig(URL)
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
