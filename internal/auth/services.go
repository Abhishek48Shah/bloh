package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage *UserStorage
}

var ErrUserExit = errors.New("user already exist")

func NewAuthService(storage *UserStorage) *AuthService {
	return &AuthService{storage: storage}
}
func (s *AuthService) Authenticate(ctx context.Context, payload LoginRequest) (*Session, error) {
	var user User
	err := s.storage.GetByEmail(ctx, *payload.Email, &user)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(*payload.Password)); err != nil {
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
	err = s.storage.StoreSession(ctx, token, data, expireAt)
	if err != nil {
		return nil, err
	}
	return &Session{Key: token, ExpireAt: expireAt}, nil
}

func (s *AuthService) Register(ctx context.Context, payload SignupRequest) (*Session, error) {
	exists, err := s.storage.CheckUser(ctx, *payload.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExit
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(*payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	*payload.Password = string(hashPassword)
	var user User
	err = s.storage.Create(ctx, payload, &user)
	if err != nil {
		return nil, err
	}
	buff := make([]byte, 32)
	if _, err := rand.Read(buff); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(buff)
	expireAt := 7 * 24 * time.Hour
	userInfo, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = s.storage.StoreSession(ctx, token, userInfo, expireAt)
	if err != nil {
		return nil, err
	}
	return &Session{Key: token, ExpireAt: expireAt}, nil
}
