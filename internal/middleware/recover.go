package middleware

import (
	"fmt"
	"github.com/Abhishek48Shah/blog/internal/util"
	"net/http"
)

func (m *Middleware) HandlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				err := fmt.Errorf("/panic recover: %w", err)
				util.WriteError(w, r, util.InternalServerErr("something went wrong", err), m.logger)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
