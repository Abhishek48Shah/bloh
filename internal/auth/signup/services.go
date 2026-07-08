package signup

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SignupService(pool *pgxpool.Pool, ctx context.Context, payload *SignupRequest) (int, error) {
	exists, err := CheckUser(pool, ctx, payload.Email)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("user already exist with that email address")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	payload.Password = string(hashPassword)
	userId, err := Create(pool, ctx, payload)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
