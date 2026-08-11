package login

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func LoginService(ctx context.Context, pool *pgxpool.Pool, cache *redis.Client, email, password string) (*Session, error) {
	user, err := GetByEmail(ctx, pool, email)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return nil, err
	}
	buff := make([]byte, 32)
	if _, err := rand.Read(buff); err != nil {

		return nil, errors.New("failed to generate token")
	}
	userInfo := struct {
		Id   int
		Role string
	}{
		Id:   user.Id,
		Role: user.Role,
	}
	data, err := json.Marshal(userInfo)
	if err != nil {
		return nil, err
	}
	token := hex.EncodeToString(buff)
	expireAt := 7 * 24 * time.Hour
	err = StoreSession(ctx, cache, token, data, expireAt)
	if err != nil {
		return nil, err
	}
	return &Session{Key: token, ExpireAt: expireAt}, nil
}
