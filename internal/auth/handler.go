package auth

import (
	"log/slog"
	"net/http"
)

type AuthHandler struct {
	service *AuthService
	logger  *slog.Logger
}

func NewAuthHandler(service *AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{service: service, logger: logger}
}
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

}
