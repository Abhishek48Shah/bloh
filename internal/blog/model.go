package blog

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Abhishek48Shah/blog/internal/util"
	"github.com/go-playground/validator/v10"
)

type Blog struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
type CreateBlogInput struct {
	Title   string
	Content string
}
type UpdateBlogInput struct {
	Title   *string
	Content *string
}
type Pagination struct {
	limit  int
	curosr *string
}

var ErrMismatchedHMC = errors.New("crypto/hmac: signature error")
var ErrInvalidTokenFormat = errors.New("invalid token format")

func validateBlog(validate *validator.Validate, input any) error {
	if err := validate.Struct(input); err != nil {
		validationErr := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Title":
				switch err.Tag() {
				case "required":
				case "min":
					validationErr["title"] = "title must be 20 character long"
				case "max":
				}
			case "Content":
				switch err.Tag() {
				case "required":
				}
				validationErr[err.Field()] = fmt.Sprintf("")
			}
		}
		return util.ValidationErr(validationErr)
	}
	return nil
}
func sign(data []byte, secretKey []byte) []byte {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(data)
	return mac.Sum(nil)
}
func Decode(token string, secret []byte) (string, error) {
	if token == "" {
		return "", nil
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", ErrInvalidTokenFormat
	}
	payload, err := base32.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	if !hmac.Equal(sign(payload, secret), []byte(parts[1])) {
		return "", ErrMismatchedHMC
	}
	var cursor string
	err = json.Unmarshal(payload, &cursor)
	if err != nil {
		return "", err
	}
	return cursor, nil
}
func Encode(cursor string, secretKey []byte) (string, error) {
	payload, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	raw := base32.StdEncoding.EncodeToString(payload)
	signature := string(sign(payload, secretKey))
	return fmt.Sprintf("%s.%s", raw, signature), nil
}
