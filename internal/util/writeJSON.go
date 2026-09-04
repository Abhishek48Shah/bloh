package util

import (
	"encoding/json"
	"maps"
	"net/http"
)

func WriteJson(w http.ResponseWriter, body any, status int, header http.Header) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	maps.Insert(w.Header(), maps.All(header))
	w.WriteHeader(status)
	w.Write(data)
	return nil
}
