package login

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const MaxBodySize = 1024

func HandleLogin(pool *pgxpool.Pool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value := r.Header.Get("Content-Type"); value != "application/json" || value == "" {
			util.UnsupportedMediaType(w, r)
			return
		}
		var payload LoginRequest
		err := util.ReadJSON(w, r, MaxBodySize, &payload)
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
		user_id, err := LoginService(ctx, pool, payload.Email, payload.Password)
		fmt.Println(user_id)
		if err != nil {

			switch {
			case errors.Is(err, pgx.ErrNoRows):
				util.NotFoundError(w, r, err)
				return
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
				util.UnauthorizedError(w, r, "invalid email or password")
				return
			}
			util.InternalServerError(w, r, err)
			return
		}
		err = util.WriteJSON(w, payload, http.StatusAccepted, nil)
		if err != nil {
			util.InternalServerError(w, r, err)
			return
		}
	})
}
