package auth

import (
	"context"
	"fmt"

	"github.com/Abhishek48Shah/blog/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage *Storage
}

func NewAuthService(storage *Storage) *AuthService {
	return &AuthService{storage: storage}
}
func (s *AuthService) GetByEmail(ctx context.Context, data *LoginInput) (*User, error) {
	var user User
	err := s.storage.FindByEmail(ctx, data.Email, &user)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		return nil, util.UnauthorizedErr("invalid email or password", err)
	}
	return nil, nil
}
