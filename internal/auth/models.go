package auth

import (
	"encoding/json"
	"fmt"

	"github.com/Abhishek48Shah/blog/internal/util"
	"github.com/go-playground/validator/v10"
)

type LoginInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
type SignupInput struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
type User struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type Session struct {
	Id   int    `redis:"id"`
	Role string `redis:"role"`
}

func validateForm(validate *validator.Validate, inputs any) error {
	validationErr := make(map[string]string)
	if err := validate.Struct(inputs); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				switch err.Tag() {
				case "min":
					validationErr["username"] = fmt.Sprintf("username must be %s characters long ", err.Param())
				case "required":

					validationErr["username"] = fmt.Sprintf("user name is required ", err.Param())
				case "max":
					validationErr["username"] = fmt.Sprintf("username must be at most %s characters long", err.Param())
				}
				validationErr[err.Field()] = fmt.Sprintf("username must be 3 character long or greater")
			case "Email":
				switch err.Tag() {
				case "email":
					validationErr["email"] = "email must be valid"
				case "required":
					validationErr["email"] = "email is required"
				}
			case "Password":
				switch err.Tag() {
				case "required":
					validationErr["password"] = "password is required"
				case "min":
					validationErr["password"] = fmt.Sprintf("password must be %s characters long ", err.Param())
				case "max":
					validationErr["password"] = fmt.Sprintf("password must be at most %s characters long", err.Param())
				}
			}

		}
		return util.ValidationErr(validationErr)
	}
	return nil
}
func marshalSessios(session *Session) ([]byte, error) {
	return json.Marshal(session)
}
