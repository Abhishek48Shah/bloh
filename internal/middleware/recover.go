package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Abhishek48Shah/blog/internal/util"
)

func HandlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				err := fmt.Errorf("/panic recover: %w", err)
				util.WriteError(w, r, util.InternalServerErr("something went wrong", err), slog.New(slog.NewTextHandler(os.Stdout, nil)))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
