package util

import (
	"log/slog"
	"net/http"
)

type envlope map[string]any

func logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	slog.Error(err.Error(), "method", method, "uri", uri)
}
func errorResponse(w http.ResponseWriter, status int, message any) {
	env := envlope{"error": message}
	err := WriteJSON(w, env, status, nil)
	if err != nil {

	}

}
func UnsupportedMediaType(w http.ResponseWriter, r *http.Request) {
	message := "Invalid or missing request headers"
	errorResponse(w, http.StatusUnsupportedMediaType, message)
}
func InternalServerError(w http.ResponseWriter, r *http.Request, status int, err error) {
	logError(r, err)
	message := "Interval server error, Please try again later"
	errorResponse(w, http.StatusInternalServerError, message)
}
