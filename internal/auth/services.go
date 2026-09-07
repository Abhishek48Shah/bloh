package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Abhishek48Shah/blog/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage *Storage
}

func NewAuthService(storage *Storage) *AuthService {
	return &AuthService{storage: storage}
}
func (s *AuthService) GetByEmail(ctx context.Context, data *LoginInput) (*Token, error) {
	var user User
	err := s.storage.FindByEmail(ctx, data.Email, &user)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return nil, util.UnauthorizedErr("invalid email or password", err)
	}
	randByt := make([]byte, 32)
	if _, err = rand.Read(randByt); err != nil {
		return nil, util.InternalServerErr("internal server error", err)
	}
	key := hex.EncodeToString(randByt)
	expiredAt := 7 * 24 * time.Hour
	session := Session{Id: user.ID, Role: user.Role}
	if err := s.storage.StoreSession(ctx, &key, expiredAt, &session); err != nil {
		return nil, err
	}
	return &Token{Key: key, ExpireAt: expiredAt}, nil
}
func (s *AuthService) Register(ctx context.Context, data *SignupInput) (*Token, error) {
	isExist, err := s.storage.Check(ctx, data.Email)
	if err != nil {
		return nil, err
	}
	if isExist {
		return nil, util.ConflictErr("email already registered", fmt.Errorf("registration attempted with an existing email"))
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), 12)
	if err != nil {
		return nil, util.InternalServerErr("internal server error", err)
	}
	data.Password = string(hashedPassword)
	var user User
	err = s.storage.Create(ctx, data, user)
	if err != nil {
		return nil, err
	}
	randByt := make([]byte, 32)
	if _, err = rand.Read(randByt); err != nil {
		return nil, util.InternalServerErr("internal server error", err)
	}
	key := hex.EncodeToString(randByt)
	expiredAt := 7 * 24 * time.Hour
	session := Session{Id: user.ID, Role: user.Role}
	if err := s.storage.StoreSession(ctx, &key, expiredAt, &session); err != nil {
		return nil, err
	}
	return &Token{Key: key, ExpireAt: expiredAt}, nil
}
