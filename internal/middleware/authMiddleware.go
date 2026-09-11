package middleware

import (
	"context"
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/util"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			util.WriteError(w, r, util.UnauthorizedErr("unauthorized access", err), nil)
		}
		ctx := context.WithValue(r.Context(), "", cookie)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
func GetUserId(r http.Request) int {
	value := r.Context().Value()
}
