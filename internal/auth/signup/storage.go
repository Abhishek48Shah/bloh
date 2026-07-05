package signup

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func CheckUser(pool *pgxpool.Pool, ctx context.Context, email string) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	var exists bool
	if err := pool.QueryRow(ctxTimeout, query, email).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil

}
func Create(pool *pgxpool.Pool, ctx context.Context, payload *SignupRequest) (int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `INSERT INTO users (username,email,password) VALUES ($1,$2,$3) RETURNING id`
	var userId int
	err := pool.QueryRow(ctxTimeout, query, payload.Password, payload.Email, payload.Password).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
