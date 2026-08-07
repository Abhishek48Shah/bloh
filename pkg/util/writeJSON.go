package util

import (
	"encoding/json"
	"log/slog"
	"maps"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, body any, status int, headers http.Header) error {
	data, err := json.Marshal(body)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	maps.Insert(w.Header(), maps.All(headers))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
	return nil
}
func SucessResponse(w http.ResponseWriter, message string, data any, status int, headers http.Header) error {
	payload := envlope{"message": message}
	if data != nil {
		payload["body"] = data
	}
	err := WriteJSON(w, payload, status, headers)
	if err != nil {
		return err
	}
	return nil
}
