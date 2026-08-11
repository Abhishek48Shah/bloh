package middleware

import (
	"github.com/Abhishek48Shah/blog/pkg/util"
	"net/http"
)

func CorsMiddleware(next http.Handler) http.Handler {
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
