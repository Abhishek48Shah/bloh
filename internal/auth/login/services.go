package login

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func LoginService(ctx context.Context, pool *pgxpool.Pool, email, password string) (int, error) {
	user, err := GetByEmail(ctx, pool, email)
	if err != nil {
		return 0, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return 0, err
	}
	return user.Id, nil
}
