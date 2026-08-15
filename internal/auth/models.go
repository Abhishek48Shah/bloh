package auth

import (
	"time"
)

type LoginRequest struct {
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required,min=8"`
}
type SignupRequest struct {
	Username *string `json:"username" validate:"required,min=3"`
	Email    *string `json:"email" validate:"required,email"`
	Password *string `json:"password" validate:"required,min=8"`
}

type User struct {
	Id             int
	HashedPassword string
	Role           string
}

type Session struct {
	Key      string
	ExpireAt time.Duration
}
