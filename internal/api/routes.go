package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/auth"
	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Handler(pool *pgxpool.Pool, cache *redis.Client, logger *slog.Logger, validator *validator.Validate) http.Handler {
	mux := http.NewServeMux()
	fmt.Printf("pool nil=%v\n", pool == nil)
	authStorage := auth.NewAuthStorage(pool, cache)
	authService := auth.NewAuthService(authStorage)
	authHandler := auth.NewAuthHandler(authService, logger, validator)
	mux.HandleFunc("POST /v1/login", authHandler.LoginHandler)
	return middleware.HandlePanic(mux)
}
