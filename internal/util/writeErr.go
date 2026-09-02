package util

import (
	"errors"
	"log/slog"
	"net/http"
)

type ApiError struct {
	Status   int
	Message  string
	Err      error
	metadata map[string]string
}

func (e *ApiError) error() string {
	return e.Message
}

type envlope map[string]string

func WriteError(w http.ResponseWriter, r *http.Request, err error, logger *slog.Logger) {
	var apiErr *ApiError
	logger.Error("")
	if errors.As(err, &apiErr) {
		logger.Error("api", "error", apiErr.Err, "method", r.Method, "uri", r.RequestURI)
		env := envlope{"error": apiErr.Message}
		WriteJson(w, env, apiErr.Status, nil)
		return
	}
	WriteJson(w, envlope{"error": "Internal server error"}, http.StatusInternalServerError, nil)
}
func BadRequestErr(Message string, Err error) *ApiError {
	return &ApiError{Status: http.StatusBadRequest, Message: Message, Err: Err}
}
func ValidationErr(message map[string]string) *ApiError {
	return &ApiError{Status: http.StatusUnprocessableEntity, Message: message, Err: nil}
}
func TimeoutErr(message string, Err error) *ApiError {
	return &ApiError{Status: http.StatusGatewayTimeout, Message: message, Err: nil}
}
func NotFoundErr(message string, Err error) *ApiError {
	return &ApiError{Status: http.StatusNotFound, Message: message, Err: Err}
}
