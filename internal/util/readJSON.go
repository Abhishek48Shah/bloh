package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ReadJSON(w http.ResponseWriter, r *http.Request, size int64, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, size)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(dst)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var invalidUnmarshalErr *json.InvalidUnmarshalError
		var unMarshalTypeErr *json.UnmarshalTypeError
		var maxByteError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxErr):
			fmt.Println(err)
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxErr.Offset)
		case errors.As(err, &invalidUnmarshalErr):
			fmt.Println(err)
			return err
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &unMarshalTypeErr):
			if unMarshalTypeErr.Field != "" {
				return fmt.Errorf("body contains badly-formed JSON type (at character %q)", unMarshalTypeErr.Field)
			}
			return fmt.Errorf("body contains badly-formed JSON type (at character %d)", unMarshalTypeErr.Offset)
		case strings.HasPrefix(err.Error(), "json: unkown field"):
			field := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", field)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &maxByteError):
			return fmt.Errorf("body must not be larger than %d bytes", maxByteError.Limit)
		default:
			return err
		}
	}
	err = decoder.Decode(struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contains single JSON value")
	}
	return nil
}
