package routes

import (
	"github.com/Abhishek48Shah/blog/internal/auth"
	"github.com/Abhishek48Shah/blog/internal/blog"
	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"net/http"
)

func NewRoutes(pool *pgxpool.Pool, cache *redis.Client, validator *validator.Validate) http.Handler {
	mux := http.NewServeMux()
	// Login and Signup endpoints
	authStorage := auth.NewAuthStorage(pool, cache)
	authService := auth.NewAuthService(authStorage)
	authHandler := auth.NewAuthHandler(authService, validator)
	mux.HandleFunc("POST /v1/login", authHandler.Login)
	mux.HandleFunc("POST /v1/register", authHandler.Signup)
	//Blog endpoints
	blogSorage := blog.NewBlogStorage(pool)
	blogService := blog.NewBlogServices(blogSorage)
	blogHandler := blog.NewBlogHandlers(blogService, validator)
	mux.HandleFunc("GET /v1/blog/{id}", blogHandler.GetSingleBlog)
	mux.Handle("POST /blog", middleware.VerifyAuth(http.HandlerFunc(blogHandler.CreateBlog), cache))
	return middleware.RecoverPanic(middleware.CorsMiddleware(mux))
}
