package signup

import (
	"github.com/Abhishek48Shah/blog/pkg/util"
	"regexp"
)

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id string `json:"id"`
}

var emailRex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func SignupValidator(validator *util.Validator, payload SignupRequest) {
	validator.Check("username", payload.Username != "", "username is required")
	validator.Check("username", len(payload.Username) < 3, "username must be 3 or more character long")
	validator.Check("email", payload.Email != "", "email is required")
	validator.Check("email", emailRex.MatchString(payload.Email), "invalid email format")
	validator.Check("password", payload.Password != "", "password is required")
	validator.Check("password", len(payload.Password) >= 8, "password must be  8 or more characters long")

}
