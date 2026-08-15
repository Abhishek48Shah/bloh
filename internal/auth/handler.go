package auth

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	service   *AuthService
	validator *validator.Validate
}

func NewAuthHandler(service *AuthService, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{service: service, validator: validator}
}

const MaxBodySize_L = 1024
const MaxBodySize_S = 2048

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		util.UnsupportedMediaType(w, r)
		return
	}
	var payload LoginRequest
	err = util.ReadJSON(w, r, MaxBodySize_L, &payload)
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}
	err = h.validator.Struct(payload)
	if err != nil {
		var let = make(map[string]string)
		util.FailedValidationResponse(w, r, let)
		return
	}
	ctx := r.Context()
	token, err := h.service.Authenticate(ctx, payload)
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword), errors.Is(err, pgx.ErrNoRows):
			util.UnauthorizedError(w, r, err, "invalid email or password")
			return
		}
		util.InternalServerError(w, r, err)
		return
	}
	session := http.Cookie{Name: "Session", Value: token.Key, Expires: time.Now().Add(token.ExpireAt)}
	http.SetCookie(w, &session)
	headers := make(http.Header)
	headers.Set("Cache-Controll", "no-store")
	err = util.WriteJSON(w, "login successfull", http.StatusOK, headers)
	if err != nil {
		util.InternalServerError(w, r, err)
		return
	}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	memiType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || memiType != "application/json" {
		util.UnsupportedMediaType(w, r)
		return
	}
	var input SignupRequest
	err = util.ReadJSON(w, r, MaxBodySize_S, &input)
	if err != nil {
		util.BadRequest(w, r, err)
		return
	}
	ctx := r.Context()
	token, err := h.service.Register(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserExit):
			util.ConfilctError(w, r, err)
		default:
			util.InternalServerError(w, r, err)
		}
		return
	}
	cookie := &http.Cookie{Name: "Session", Value: token.Key, Domain: "localhost:4000", Expires: time.Now().Add(token.ExpireAt)}
	http.SetCookie(w, cookie)
	headers := make(http.Header)
	headers.Set("Cache-Control", "no-store")
	err = util.WriteJSON(w, "Signup of successfull", 201, headers)
	if err != nil {
		util.InternalServerError(w, r, fmt.Errorf("/signup/handler.go: writeJSON err: %w", err))
		return
	}
}
