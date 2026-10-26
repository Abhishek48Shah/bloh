package api

import (
	"log/slog"
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/auth"
	"github.com/Abhishek48Shah/blog/internal/blog"
	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/Abhishek48Shah/blog/internal/pagination"
	"github.com/Abhishek48Shah/blog/internal/token"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Routes struct {
	pool       *pgxpool.Pool
	cache      *redis.Client
	logger     *slog.Logger
	validation *validator.Validate
	tokenSrv   *token.TokenService
	middleware *middleware.Middleware
	codec      *pagination.Codec
}

func NewRoutes(pool *pgxpool.Pool,
	cache *redis.Client,
	logger *slog.Logger,
	validation *validator.Validate,
	tokenSrv *token.TokenService,
	middleware *middleware.Middleware,
	codec *pagination.Codec,
) *Routes {
	return &Routes{pool: pool, cache: cache, logger: logger, validation: validation, tokenSrv: tokenSrv, middleware: middleware, codec: codec}
}

func (r *Routes) Handler() http.Handler {
	mux := http.NewServeMux()
	/**
	// auth endpoints
	**/
	authStorage := auth.NewAuthStorage(r.pool, r.cache)
	authService := auth.NewAuthService(authStorage, r.tokenSrv)
	authHandler := auth.NewAuthHandler(authService, r.logger, r.validation)
	mux.HandleFunc("POST /v1/login", authHandler.LoginHandler)
	mux.HandleFunc("POST /v1/signup", authHandler.SignupHandler)
	mux.Handle("POST /v1/logout", r.middleware.Auth(http.HandlerFunc(authHandler.LogoutHandler)))

	/**
	// blog endpoints
	**/
	blogStorage := blog.NewStorage(r.pool, r.cache)
	blogService := blog.NewService(blogStorage)
	blogHandler := blog.NewHandler(blogService, r.logger, r.validation, r.codec)
	/**
	//public endpoints
	**/
	mux.HandleFunc("GET /v1/blog/{id}", blogHandler.GetSingleBlogHandler)
	/**
	//private endpoints
	**/
	mux.Handle("POST /v1/blog", r.middleware.Auth(r.middleware.RequiredRole("user", http.HandlerFunc(blogHandler.CreateBlogHandler))))
	mux.Handle("DELETE /v1/blog/{id}", r.middleware.Auth(r.middleware.RequiredRole("user", http.HandlerFunc(blogHandler.DeleteBlogHandler))))
	mux.Handle("PATCH /v1/blog/{id}", r.middleware.Auth(r.middleware.RequiredRole("user", http.HandlerFunc(blogHandler.EditBlogHandler))))
	mux.Handle("GET /v1/blog", http.HandlerFunc(blogHandler.GetUserBlogsHandler))
	/**
	//
	**/
	return r.middleware.HandlePanic(r.middleware.HandleCors(mux))
}
