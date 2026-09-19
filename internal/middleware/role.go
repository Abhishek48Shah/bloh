package middleware

import (
	"fmt"
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/util"
)

func (m *Middleware) RequiredRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := GetUserContext(r)
		if session == nil {
			util.WriteError(w, r, util.UnauthorizedErr("unauthorized access", fmt.Errorf("no session in context")), m.logger)
			return
		}
		if session.Role != role {
			util.WriteError(w, r, util.ForbiddenErr("you don't have permission to access this resources", fmt.Errorf("no role in context")), m.logger)
			return
		}
		next.ServeHTTP(w, r)
	})
}
