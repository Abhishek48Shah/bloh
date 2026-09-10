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
		util.WriteError(w, r, util.UnsupportedMediaTypeErr("malformed header", nil), h.logger)
		return
	}
	input := LoginInput{}
	err = util.ReadJSON(w, r, LoginMaxBodySize, &input)
	if err != nil {
		util.WriteError(w, r, util.BadRequestErr(err.Error(), err), h.logger)
		return
	}
	if err := h.validate.Struct(input); err != nil {
		validationErr := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Println("tag: ", err.Tag(), "field: ", err.Field(), "actual-field: ", err.ActualTag())
			fmt.Println("namespace: ", err.Namespace(), "kind: ", err.Kind(), "param: ", err.Param())

			switch err.Field() {
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
		util.WriteError(w, r, util.ValidationErr(validationErr), h.logger)
		return
	}
	session, err := h.service.GetByEmail(r.Context(), &input)
	if err != nil {
		util.WriteError(w, r, err, h.logger)
		return
	}
	cookie := &http.Cookie{Name: "session", Value: session.Key, MaxAge: int(session.ExpireAt.Seconds())}
	http.SetCookie(w, cookie)
	if err := util.WriteJson(w, "login successfull", http.StatusAccepted, nil); err != nil {
		util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		return
	}
}
func (h *AuthHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		util.WriteError(w, r, util.BadRequestErr("malformed header", err), h.logger)
		return
	}
	if mediaType != "application/json" {
		util.WriteError(w, r, util.UnsupportedMediaTypeErr("malformed header", fmt.Errorf("malformd headers: %s", mediaType)), h.logger)
		return
	}
	var input SignupInput
	if err := util.ReadJSON(w, r, SignupMaxBodySize, &input); err != nil {
		util.WriteError(w, r, util.BadRequestErr(err.Error(), err), h.logger)
		return
	}
	if err := h.validate.Struct(input); err != nil {
		validationErr := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				switch err.Tag() {
				case "min":
				case "required":

					validationErr["password"] = fmt.Sprintf("password must be %s characters long ", err.Param())
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
		util.WriteError(w, r, util.ValidationErr(validationErr), h.logger)
		return
	}
	session, err := h.service.Register(r.Context(), &input)
	if err != nil {
		util.WriteError(w, r, err, h.logger)
	}
	cookie := &http.Cookie{Name: "session", Value: session.Key, MaxAge: int(session.ExpireAt.Seconds())}
	http.SetCookie(w, cookie)
	if err := util.WriteJson(w, "account register successfully", http.StatusAccepted, nil); err != nil {
		util.WriteError(w, r, util.InternalServerErr("internal server error", err), h.logger)
		return
	}
}
