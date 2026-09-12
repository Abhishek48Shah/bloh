package api

import (
	"github.com/Abhishek48Shah/blog/internal/auth"
	"github.com/Abhishek48Shah/blog/internal/blog"
	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"net/http"
)

func Handler(pool *pgxpool.Pool, cache *redis.Client, logger *slog.Logger, validator *validator.Validate) http.Handler {
	mux := http.NewServeMux()
	middleware := middleware.NewMiddleware(logger, cache)
	/**
	// auth endpoints
	**/
	authStorage := auth.NewAuthStorage(pool, cache)
	authService := auth.NewAuthService(authStorage)
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
