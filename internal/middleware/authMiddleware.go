package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Abhishek48Shah/blog/internal/util"
	"github.com/redis/go-redis/v9"
)

type contextKey string
type Payload struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	Key  string `json:"-"`
}

var sessionKey contextKey = "session"

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			util.WriteError(w, r, util.UnauthorizedErr("no cookie provided", err), m.logger)
			return
		}
		value := cookie.Value
		token, err := m.tokenSrv.VerifyToken(value)
		if err != nil {
			util.WriteError(w, r, util.UnauthorizedErr("unauthorized access", fmt.Errorf("token verification failed: %w", err)), m.logger)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		var payload Payload
		data, err := m.cache.Get(ctx, token.ID).Bytes()
		if err != nil {
			fmt.Println(err)
			switch {
			case errors.Is(err, redis.Nil):
				util.WriteError(w, r, util.UnauthorizedErr("unauthorized access", fmt.Errorf("key doesn't exist", err)), m.logger)
			case errors.Is(err, redis.ErrPoolTimeout), errors.Is(err, context.DeadlineExceeded):
				util.WriteError(w, r, util.TimeoutErr("request time out.", err), m.logger)
			default:
				util.WriteError(w, r, util.InternalServerErr("internal server error", err), m.logger)
			}
			return
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			util.WriteError(w, r, util.InternalServerErr("internal server error", fmt.Errorf("unmarshal session: %w", err)), m.logger)
			return
		}
		if err != nil {
			switch {
			case errors.Is(err, redis.ErrPoolTimeout):
				util.WriteError(w, r, util.TimeoutErr("request time out", fmt.Errorf("/authmiddleware redis get: %w", err)), m.logger)
			case errors.Is(err, redis.Nil):
				util.WriteError(w, r, util.UnauthorizedErr("unauthorized access", fmt.Errorf("/authMiddleware redis get : %w", err)), m.logger)
			default:
				util.WriteError(w, r, util.InternalServerErr("internal server error", fmt.Errorf("/authMiddleware redis get: %w", err)), m.logger)

			}
			return
		}
		payload.Key = token.ID
		ctx = context.WithValue(r.Context(), sessionKey, payload)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
func GetUserContext(r *http.Request) *Payload {
	payload := r.Context().Value(sessionKey).(Payload)
	return &payload
}
