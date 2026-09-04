package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"mime"
	"net/http"

	"github.com/Abhishek48Shah/blog/internal/util"
	"github.com/go-playground/validator/v10"
)

var (
	SignupMaxBodySize int64 = 1024
	LoginMaxBodySize  int64 = 1024
)

type AuthHandler struct {
	service  *AuthService
	logger   *slog.Logger
	validate *validator.Validate
}

func NewAuthHandler(service *AuthService, logger *slog.Logger, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{service: service, logger: logger, validate: validator}
}

var BadRequestErr = errors.New("invalid input formate")

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		util.WriteError(w, r, util.BadRequestErr("malformed header", err), h.logger)
		return
	}
	if mediaType != "application/json" {
		fmt.Println("hello error")
		util.WriteError(w, r, util.UnsupportedMediaTypeErr("malformed header", nil), h.logger)
		return
	}
	input := LoginInput{}
	err = util.ReadJSON(w, r, LoginMaxBodySize, &input)
	if err != nil {
		util.WriteError(w, r, util.BadRequestErr(err.Error(), err), h.logger)
		return
	}
	err = h.validate.Struct(input)
	if err != nil {
		validationErr := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Tag())
			switch err.Field() {
			case "required":
				validationErr[err.Field()] = fmt.Sprintf("%s is required,err.Field()")
			case "Email":
				validationErr[err.Field()] = fmt.Sprint("email must be valid")
			case "Password":
				validationErr[err.Field()] = fmt.Sprint("password must be eight character long or greater")
			}
		}
		errMessage := util.ValidationErr(validationErr)
		util.WriteError(w, r, errMessage, h.logger)
		return
	}
	_, err = h.service.GetByEmail(r.Context(), &input)
	if err != nil {
		util.WriteError(w, r, err, h.logger)
		return
	}

}
