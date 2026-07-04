package login

import (
	"regexp"

	"github.com/Abhishek48Shah/blog/pkg/util"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type User struct {
	Id             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"`
}

var emailRex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func payloadValidation(validator *util.Validator, payload *LoginRequest) {
	validator.Check("email", payload.Email != "", "email is required")
	validator.Check("email", emailRex.MatchString(payload.Email), "invalid email format")
	validator.Check("password", payload.Password != "", "password is required")
	validator.Check("password", len(payload.Password) >= 8, "password must be  8 or more characters long")
}
