package auth

import (
	"encoding/json"
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
		fmt.Println(err)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "unsupported media type", http.StatusUnsupportedMediaType)
		return
	}
	input := LoginInput{}
	err = util.ReadJSON(w, r, LoginMaxBodySize, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.validate.Struct(input)
	if err != nil {
		validationErr := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println(err.Field() + " " + err.Tag())
			switch err.Field() {
			case "required":
				fmt.Println(err.Tag())
				fmt.Println(err.Field())
				validationErr[err.Field()] = fmt.Sprintf("%s is required", err.Field())
			case "Email":
				validationErr[err.Field()] = fmt.Sprintf("%s must be valid email", err.Field())
			case "Password":
				validationErr[err.Field()] = fmt.Sprintf("password must be eight character long or greater")
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErr)
		return
	}
	ctx := r.Context()
	_, err = h.service.GetByEmail(ctx, &input)
	if err != nil {
		fmt.Println(err)
	}

}
