package login

import (
	"errors"
	"mime"
	"net/http"
	"time"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const MaxBodySize = 1024

func HandleLogin(pool *pgxpool.Pool, redisClient *redis.Client) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			util.UnsupportedMediaType(w, r)
			return
		}
		var payload LoginRequest
		err = util.ReadJSON(w, r, MaxBodySize, &payload)
		if err != nil {
			util.BadRequest(w, r, err)
			return
		}
		validator := util.NewValidator()
		payloadValidation(&validator, &payload)
		if ok := validator.Ok(); !ok {
			util.FailedValidationResponse(w, r, validator.Error)
			return
		}
		ctx := r.Context()
		token, err := LoginService(ctx, pool, redisClient, payload.Email, payload.Password)
		if err != nil {
			switch {
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword), errors.Is(err, pgx.ErrNoRows):
				util.UnauthorizedError(w, r, "invalid email or password")
				return
			}
			util.InternalServerError(w, r, err)
			return
		}
		session := http.Cookie{Name: "Session", Value: token.Key, Expires: time.Now().Add(token.ExpireAt)}
		http.SetCookie(w, &session)
		err = util.WriteJSON(w, "login successfull", http.StatusOK, nil)
		if err != nil {
			util.InternalServerError(w, r, err)
			return
		}
	})
}
