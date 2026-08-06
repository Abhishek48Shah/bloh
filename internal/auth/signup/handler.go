package signup

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Abhishek48Shah/blog/pkg/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const maxBodySize = 2048

func HandleSignup(pool *pgxpool.Pool, redisClient *redis.Client) http.HandlerFunc {
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
		token, err := SignupService(pool, redisClient, ctx, &input)
		if err != nil {
			util.InternalServerError(w, r, err)
			return
		}
		cookie := &http.Cookie{Name: "Session", Value: token.Key, Domain: "localhost:4000", Expires: time.Now().Add(token.Expire_at)}
		http.SetCookie(w, cookie)
		util.WriteJSON(w, nil, 201, nil)
	})
}
