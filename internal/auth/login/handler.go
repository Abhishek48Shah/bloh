package login

import (
	"net/http"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleLogin(pool *pgxpool.Pool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value := r.Header.Get("Content-Type"); value != "application/json" || value == "" {
			util.UnsupportedMediaType(w, r)
			return
		}
	})
}
