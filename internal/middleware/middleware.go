package middleware

import (
	"log/slog"

	"github.com/Abhishek48Shah/blog/internal/token"
	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	logger   *slog.Logger
	cache    *redis.Client
	tokenSrv *token.TokenService
}

func NewMiddleware(logger *slog.Logger, cache *redis.Client, tokenSrv *token.TokenService) *Middleware {
	return &Middleware{logger: logger, cache: cache, tokenSrv: tokenSrv}
}
