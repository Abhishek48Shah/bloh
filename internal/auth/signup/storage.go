package signup

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
func Create(pool *pgxpool.Pool, ctx context.Context, payload *SignupRequest) (*User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `INSERT INTO users (username,email,password) VALUES ($1,$2,$3) RETURNING id,role`
	var user User
	err := pool.QueryRow(ctxTimeout, query, payload.Username, payload.Email, payload.Password).Scan(&user.Id, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func StoreSession(redisClient *redis.Client, ctx context.Context, token string, user []byte, expires_At time.Duration) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := redisClient.Set(ctxTimeout, token, user, expires_At).Err()
	if err != nil {
		return err
	}
	return nil
}
