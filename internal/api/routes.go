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
	mux.Handle("POST /v1/logout", middleware.Auth(http.HandlerFunc(authHandler.LogoutHandler)))

	/**
	// blog endpoints
	**/
	blogStorage := blog.NewStorage(pool, cache)
	blogService := blog.NewService(blogStorage)
	blogHandler := blog.NewHandler(blogService, logger, validator)
	mux.HandleFunc("GET /v1/blog/{id}", blogHandler.GetSingleBlogHandler)
	// mux.HandleFunc("GET /v1/blog", middleware.Auth(middleware.RequiredRole("user", http.HandlerFunc(blogHandler.GetUserBlogsHandler))))
	mux.Handle("POST /v1/blog", middleware.Auth(middleware.RequiredRole("user", http.HandlerFunc(blogHandler.CreateBlogHandler))))
	mux.Handle("DELETE /v1/blog/{id}", middleware.Auth(middleware.RequiredRole("user", http.HandlerFunc(blogHandler.DeleteBlogHandler))))
	mux.Handle("PATCH /v1/blog/{id}", middleware.Auth(middleware.RequiredRole("user", http.HandlerFunc(blogHandler.EditBlogHandler))))
	/**
	//
	**/
	return middleware.HandlePanic(middleware.HandleCors(mux))
}
