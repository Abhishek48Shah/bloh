package api

import (
	"log/slog"
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/auth"
	"github.com/Abhishek48Shah/blog/internal/blog"
	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/Abhishek48Shah/blog/internal/token"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Handler(pool *pgxpool.Pool, cache *redis.Client, logger *slog.Logger, validator *validator.Validate, tokenSvc *token.TokenService, middleware *middleware.Middleware) http.Handler {
	mux := http.NewServeMux()
	/**
	// auth endpoints
	**/
	authStorage := auth.NewAuthStorage(pool, cache)
	authService := auth.NewAuthService(authStorage, tokenSvc)
	authHandler := auth.NewAuthHandler(authService, logger, validator)
	mux.HandleFunc("POST /v1/login", authHandler.LoginHandler)
	mux.HandleFunc("POST /v1/signup", authHandler.SignupHandler)

	/**
	// blog endpoints
	**/
	blogStorage := blog.NewStorage(pool, cache)
	blogService := blog.NewService(blogStorage)
	blogHandler := blog.NewHandler(blogService, logger, validator)
	mux.HandleFunc("GET /v1/blog/{id}", blogHandler.GetSingleBlog)
	mux.Handle("POST /v1/blog", middleware.Auth(http.HandlerFunc(blogHandler.CreateBlog)))
	return middleware.HandlePanic(mux)
}
