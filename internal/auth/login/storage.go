package login

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func GetByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (*User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `SELECT id,email,password FROM users WHERE email=$1`
	var user User
	err := pool.QueryRow(ctxTimeout, query, email).Scan(&user.Id, &user.Email, &user.HashedPassword)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
