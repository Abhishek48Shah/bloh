package login

import "github.com/Abhishek48Shah/blog/pkg/util"

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func payloadValidation(validator *util.Validator) {
	validator.Check()
}
