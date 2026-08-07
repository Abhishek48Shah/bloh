package login

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func LoginService(ctx context.Context, pool *pgxpool.Pool, redisClient *redis.Client, email, password string) (*Token, error) {
	user, err := GetByEmail(ctx, pool, email)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, err
	}
	buff := make([]byte, 32)
	n, err := rand.Read(buff)
	if n < 0 {
		return nil, errors.New("failed to generate token")
	}
	err := StoreSession()
	return nil, nil
}
