package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/redis/go-redis/v9"
)

type contextKey string

const userId contextKey = "userID"

func VerifyAuth(next http.Handler, cache *redis.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()
		if len(cookies) != 1 || cookies[0].Name != "Session" {
			util.UnauthorizedError(w, r, fmt.Errorf("middleware/VerifyAuth: cookie lenght %d", len(cookies)), "authentication cookie missing")
			return
		}
		key := cookies[0].Value
		ctxTimeout, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		isExist, err := cache.Get(ctxTimeout, key).Result()
		if err != nil {
			switch {
			case errors.Is(err, redis.Nil):
				util.UnauthorizedError(w, r, fmt.Errorf("middleware/VerifyAuth: redis GET value: %w", err), "session expired or invalid")
			default:
				util.InternalServerError(w, r, err)
			}
			return
		}
		defer cache.Close()
		userID, err := strconv.Atoi(isExist)
		if err != nil {
			util.InternalServerError(w, r, fmt.Errorf("internal/middleware/auth.go/ userId parse: %s", err))
		}
		ctx := context.WithValue(r.Context(), userId, userID)
		newRequest := r.WithContext(ctx)
		next.ServeHTTP(w, newRequest)
	})
}
func GetUserId(ctx context.Context) (int, bool) {
	val, ok := ctx.Value(userId).(int)
	return val, ok
}
