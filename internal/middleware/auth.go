package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/redis/go-redis/v9"
)

func VerifyAuth(next http.Handler, cache *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.CookiesNamed("Session")
		if len(cookies) != 1 || cookies[0].Name != "Session" {
			util.UnauthorizedError(w, r, "authentication cookie missing")
			return
		}
		key := cookies[0].Value
		ctxTimeout, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		isExist := cache.Get(ctxTimeout, key)
		if isExist.Err() != nil {

		}
		isExist.Args()
		next.ServeHTTP(w, r)
	})
}
