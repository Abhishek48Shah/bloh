package util

import (
	"log/slog"
	"net/http"

	"golang.org/x/text/message"
)

type envlope map[string]any

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, error map[string]string) {
	errorResponse(w, r, http.StatusUnprocessableEntity, error)
}
func logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	slog.Error(err.Error(), "method", method, "uri", uri)
}
func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envlope{"error": message}
	err := WriteJSON(w, env, status, nil)
	if err != nil {
		InternalServerError(w, r, err)
	}
}
func UnsupportedMediaType(w http.ResponseWriter, r *http.Request) {
	message := "Invalid or missing request headers"
	errorResponse(w, r, http.StatusUnsupportedMediaType, message)
}
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	message := "Bad request, Please try again later"
	errorResponse(w, r, http.StatusBadRequest, message)
}
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	message := "Interval server error, Please try again later"
	errorResponse(w, r, http.StatusInternalServerError, message)
}
func UnauthorizedError(w http.ResponseWriter, r *http.Request, message string) {
	errorResponse(w, r, http.StatusUnauthorized, message)
}
func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	message := "Not found"
	errorResponse(w, r, http.StatusNotFound, message)
}
