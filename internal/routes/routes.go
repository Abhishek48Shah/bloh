package routes

import (
	"fmt"
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/auth/login"
	"github.com/Abhishek48Shah/blog/internal/auth/signup"
	"github.com/Abhishek48Shah/blog/internal/blog"
	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/mod/sumdb/storage"
)

func NewRoutes(pool *pgxpool.Pool, rdb *redis.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/login", login.HandleLogin(pool, rdb))
	mux.HandleFunc("POST /v1/signup", signup.HandleSignup(pool, rdb))
	blogSorage := blog.LoadStorage(pool, rdb)
	blogSerives := blog.LoadServices(blogSorage)
	blogHandler := blog.LoadHandlers(blogSerives)
	mux.HandleFunc("GET /v1/blog/{id}", blogHandler.GetBlogById)
	return recoverPanic(corsMiddleware(mux))
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := make(http.Header)
		headers.Add("Access-Control-Allow-Orign", "*")
		headers.Add("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTION")
		headers.Add("Access-Controle-Allow-Headers", "Content-Type")
		headers.Add("Access-Controle-Allow-Credentails", "true")
		if r.Method == "OPTIONS" {
			util.WriteJSON(w, nil, http.StatusNoContent, headers)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				util.InternalServerError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
