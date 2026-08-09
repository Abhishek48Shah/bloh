package blog

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type BlogStorage struct {
	pool *pgxpool.Pool
	rdb  *redis.Client
}

func LoadStorage(pool *pgxpool.Pool, rdb *redis.Client) *BlogStorage {
	return &BlogStorage{pool: pool, rdb: rdb}
}
