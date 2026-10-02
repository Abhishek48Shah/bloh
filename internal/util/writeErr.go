package util

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type ApiError struct {
	status  int
	message any
	err     error
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%s", e.message)
}

type envelope map[string]any

func WriteError(w http.ResponseWriter, r *http.Request, err error, logger *slog.Logger) {
	var apiErr *ApiError
	if errors.As(err, &apiErr) {
		logger.Error("api", "error", apiErr.err, "method", r.Method, "uri", r.RequestURI)
		env := envelope{"error": apiErr.message}
		if err := WriteJson(w, env, apiErr.status, nil); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
}
func BadRequestErr(message string, err error) *ApiError {
	return &ApiError{status: http.StatusBadRequest, message: message, err: err}
}
func TimeoutErr(message string, err error) *ApiError {
	return &ApiError{status: http.StatusGatewayTimeout, message: message, err: err}

}
func InternalServerErr(message string, err error) *ApiError {
	return &ApiError{status: http.StatusInternalServerError, message: message, err: err}
}
func NotFounErr(message string, err error) *ApiError {
	return &ApiError{status: http.StatusNotFound, message: message, err: err}
}
func ValidationErr(message map[string]string) *ApiError {
	return &ApiError{status: http.StatusUnprocessableEntity, message: message, err: nil}
}
func UnsupportedMediaTypeErr(message string, err error) *ApiError {

	return &ApiError{status: http.StatusUnsupportedMediaType, message: message, err: err}
}
func UnauthorizedErr(message string, err error) *ApiError {
	return &ApiError{status: http.StatusUnauthorized, message: message, err: err}
}
func ConflictErr(message string, err error) *ApiError {
	return &ApiError{status: http.StatusConflict, message: message, err: err}
}
func ForbiddenErr(message string, err error) *ApiError {
	return &ApiError{status: http.StatusForbidden, message: message, err: err}
}
