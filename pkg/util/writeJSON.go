package util

import (
	"encoding/json"
	"maps"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, body any, status int, headers http.Header) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	maps.Insert(w.Header(), maps.All(headers))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
	return nil
}
