package signup

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

func SignupService(pool *pgxpool.Pool, redisClient *redis.Client, ctx context.Context, payload *SignupRequest) (*Token, error) {
	exists, err := CheckUser(pool, ctx, payload.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("user already exist with that email address")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	payload.Password = string(hashPassword)
	user, err := Create(pool, ctx, payload)
	if err != nil {
		return nil, err
	}
	buff := make([]byte, 32)
	n, err := rand.Read(buff)
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return nil, fmt.Errorf("Token is null")
	}
	token := hex.EncodeToString(buff[:n])
	expire_at := 7 * 24 * time.Hour
	err = StoreSession(redisClient, ctx, token, user, expire_at)
	if err != nil {
		return nil, err
	}
	return &Token{Key: token, Expire_at: expire_at}, nil
}
