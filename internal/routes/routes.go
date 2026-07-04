package routes

import (
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/auth/login"
	"github.com/Abhishek48Shah/blog/internal/auth/signup"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRoutes(pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/login", login.HandleLogin(pool))
	mux.HandleFunc("POST /v1/signup", signup.HandleSignup(pool))
	return mux
}
