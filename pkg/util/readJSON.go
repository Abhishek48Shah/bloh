package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, maxBodySize int, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBodySize))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dst)
	if err != nil {
		var unMarshallTypeErr json.UnmarshalTypeError
		switch {
		case errors.As(err, unMarshallTypeErr):
			panic(err)
		case errors.Is(err, io.EOF):
		}
		return err
	}
	return nil
}
