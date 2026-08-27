package api

import (
	"log/slog"
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Handler(pool *pgxpool.Pool, cache *redis.Client, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	return middleware.HandlePanic(mux)
}
