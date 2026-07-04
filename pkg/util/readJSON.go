package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, maxBodySize int, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBodySize))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&dst)
	if err != nil {
		var unMarshalTypeErr *json.UnmarshalTypeError
		var syntaxErr *json.SyntaxError
		var invalidUnmarsahlError *json.InvalidUnmarshalError
		var maxByteErr *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxErr):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxErr.Offset)
		case errors.As(err, &invalidUnmarsahlError):
			panic(err)
		case errors.As(err, &maxByteErr):
			return fmt.Errorf("body must not be large than %d bytes", maxByteErr.Limit)
		case errors.As(err, &unMarshalTypeErr):
			if unMarshalTypeErr.Field != "" {
				return fmt.Errorf("body contains badly-formed JSON type (at character %q)", unMarshalTypeErr.Field)
			}
			return fmt.Errorf("body contains badly-formed JSON type (at character %d)", unMarshalTypeErr.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("json contains badly-formed JSON")
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown field: %s", fieldName)
		}
		return err
	}
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single json value")
	}
	return nil
}
