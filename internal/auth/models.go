package auth

import "time"

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type SignupInput struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type User struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type Session struct {
	Id   int
	Role string
}
type Token struct {
	Key      string
	ExpireAt time.Duration
}
