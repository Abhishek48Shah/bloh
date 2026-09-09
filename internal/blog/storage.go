package blog

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	pool  *pgxpool.Pool
	cache *redis.Client
}

func NewStorage(pool *pgxpool.Pool, cache *redis.Client) *Storage {
	return &Storage{pool: pool, cache: cache}
}
