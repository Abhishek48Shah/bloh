package routes

import (
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/auth/login"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRoutes(pool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", login.HandleLogin(pool))
	return mux
}
