package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Abhishek48Shah/blog/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	pool  *pgxpool.Pool
	cache *redis.Client
}

func NewAuthStorage(pool *pgxpool.Pool, cache *redis.Client) *Storage {
	return &Storage{pool: pool, cache: cache}
}
func (s *Storage) FindByEmail(ctx context.Context, email string, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `SELECT id,password,role FROM users WHERE email=$1`
	err := s.pool.Ping(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = s.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Password, &user.Role)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return util.UnauthorizedErr("invalid email or password", err)
		case errors.Is(err, context.DeadlineExceeded):
			return util.TimeoutErr("request time out.", err)
		default:
			return util.InternalServerErr("internal server error", err)
		}
	}
	return nil
}
func (s *Storage) Create(ctx context.Context, data *SignupInput, user *User) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `INSERT INTO USERS username,email,password VALUES ($1,$2,$3) RETURNING id,role`
	err := s.pool.QueryRow(ctx, query, data.Username, data.Email, data.Password).Scan(&user.ID, &user.Role)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return util.TimeoutErr("request time out.", err)
		default:
			return util.InternalServerErr("internal server error", err)
		}

	}
	return nil
}
func (s *Storage) Check(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE email=$1)`
	var isExist bool
	err := s.pool.QueryRow(ctx, query, email).Scan(&isExist)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return isExist, util.TimeoutErr("request time out.", err)
		default:
			return isExist, util.InternalServerErr("internal server error", err)
		}
	}
	return isExist, nil
}
func (s *Storage) StoreSession(ctx context.Context, key *string, expireAt time.Duration, session *Session) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := s.cache.Set(ctx, *key, session, expireAt).Err(); err != nil {
		switch {
		case errors.Is(err, redis.ErrPoolTimeout), errors.Is(err, context.DeadlineExceeded):
			return util.TimeoutErr("request time out.", err)
		default:
			return util.InternalServerErr("internal server error", err)
		}
	}
	return nil
}
