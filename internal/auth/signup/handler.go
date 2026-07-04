package signup

import (
	"net/http"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HandleSignup(pool *pgxpool.Pool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if headerValue := r.Header.Get("Content-Type"); headerValue == "" || headerValue != "application/json" {
			util.UnsupportedMediaType(w, r)
			return
		}
		var input SignupRequest
		err := util.ReadJSON(w, r, 1024, &input)
		if err != nil {
			util.BadRequest(w, r, err)
			return
		}
		validator := util.NewValidator()
		SignupValidator(validator, input)
		if ok := validator.Ok(); !ok {
			util.BadRequest(w, r, err)
		}
		ctx := r.Context()
		user_id, err := SignupService(pool, ctx, &input)
		if err != nil {
			switch {

			}
			util.InternalServerError(w, r, err)
			return
		}
		util.WriteJSON(w, r, user_id, nil)
	})

}
