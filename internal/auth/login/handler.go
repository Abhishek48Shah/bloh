package login

import (
	"net/http"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

const MaxBodySize = 1024

func HandleLogin(pool *pgxpool.Pool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value := r.Header.Get("Content-Type"); value != "application/json" || value == "" {
			util.UnsupportedMediaType(w, r)
			return
		}
		payload := &LoginPayload{}
		err := util.ReadJSON(w, r, MaxBodySize, payload)
		if err != nil {
			util.BadRequest(w, r, err)
		}
		// validator := util.NewValidator()
		// payloadValidation(validator)
		//
		err = util.WriteJSON(w, payload, http.StatusAccepted, nil)
		if err != nil {
			util.InternalServerError(w, r, err)
		}
	})
}
