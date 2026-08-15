package auth

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserStorage struct {
	pool  *pgxpool.Pool
	cache *redis.Client
}

func NewAuthStorage(pool *pgxpool.Pool, cache *redis.Client) *UserStorage {
	return &UserStorage{pool: pool, cache: cache}
}
func (u *UserStorage) GetByEmail(ctx context.Context, email string, user *User) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `SELECT id,password,role FROM users WHERE email=$1`
	err := u.pool.QueryRow(ctxTimeout, query, email).Scan(&user.Id, &user.HashedPassword, &user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserStorage) CheckUser(ctx context.Context, email string) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`
	var exists bool
	if err := u.pool.QueryRow(ctxTimeout, query, email).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil

}
func (u *UserStorage) Create(ctx context.Context, payload SignupRequest, user *User) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `INSERT INTO users (username,email,password) VALUES ($1,$2,$3) RETURNING id,role`
	err := u.pool.QueryRow(ctxTimeout, query, *payload.Username, *payload.Email, *payload.Password).Scan(&user.Id, &user.Role)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserStorage) StoreSession(ctx context.Context, token string, info []byte, expireAt time.Duration) error {
	cyxTimeout, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := u.cache.Set(cyxTimeout, token, info, expireAt).Err()
	if err != nil {
		return err
	}
	return nil
}
