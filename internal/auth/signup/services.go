package signup

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SignupService(pool *pgxpool.Pool, ctx context.Context, payload *SignupRequest) (int, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}
	payload.Password = string(hashPassword)
	fmt.Println(payload.Password)
	return 0, nil
}
