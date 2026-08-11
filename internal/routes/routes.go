package routes

import (
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/auth/login"
	"github.com/Abhishek48Shah/blog/internal/auth/signup"
	"github.com/Abhishek48Shah/blog/internal/blog"
	"github.com/Abhishek48Shah/blog/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func NewRoutes(pool *pgxpool.Pool, cache *redis.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/login", login.HandleLogin(pool, cache))
	mux.HandleFunc("POST /v1/signup", signup.HandleSignup(pool, cache))
	blogSorage := blog.NewBlogStorage(pool)
	blogSerives := blog.NewBlogServices(blogSorage)
	blogHandler := blog.NewBlogHandlers(blogSerives)
	mux.HandleFunc("GET /v1/blog/{id}", blogHandler.GetSingleBlog)
	mux.Handle("POST /blog", middleware.VerifyAuth(http.HandlerFunc(blogHandler.CreateBlog), cache))
	return middleware.RecoverPanic(middleware.CorsMiddleware(mux))
}
