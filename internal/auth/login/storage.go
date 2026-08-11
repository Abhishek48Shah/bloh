package login

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func GetByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (*User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `SELECT id,email,password,role FROM users WHERE email=$1`
	var user User
	err := pool.QueryRow(ctxTimeout, query, email).Scan(&user.Id, &user.Email, &user.HashedPassword, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func StoreSession(ctx context.Context, cache *redis.Client, token string, info []byte, expireAt time.Duration) error {
	cyxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := cache.Set(cyxTimeout, token, info, expireAt).Err()
	if err != nil {
		return err
	}
	return nil
}
