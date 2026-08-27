package middleware

import "net/http"

func HandlePanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {

			}

		}()
		next.ServeHTTP(w, r)
	})
}
