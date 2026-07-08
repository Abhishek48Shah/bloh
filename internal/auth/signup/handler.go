package signup

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxBodySize = 1024

func HandleSignup(pool *pgxpool.Pool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var headerValue = r.Header.Get("Content-Type")
		if headerValue != "application/json" {
			util.UnsupportedMediaType(w, r)
			return
		}
		var input SignupRequest
		err := util.ReadJSON(w, r, maxBodySize, &input)
		fmt.Println(input)
		if err != nil {
			util.BadRequest(w, r, err)
			return
		}
		validator := util.NewValidator()
		SignupValidator(&validator, input)
		if ok := validator.Ok(); !ok {
			util.FailedValidationResponse(w, r, validator.Error)
			return
		}
		ctx := r.Context()
		userId, err := SignupService(pool, ctx, &input)
		if err != nil {
			switch {

			}
			util.InternalServerError(w, r, err)
			return
		}
		err = util.WriteJSON(w, r, userId, nil)
		if err != nil {
			util.InternalServerError(w, r, err)
		}
	})

}
