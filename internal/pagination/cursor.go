package pagination

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Codec struct {
	key []byte
}
type Cursor struct {
	BlogID    *int       `json:"blog_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

var ErrMismatchedHMC = errors.New("crypto/hmac: signature error")
var ErrInvalidTokenFormat = errors.New("invalid token format")

func NewCodec(key []byte) *Codec {
	return &Codec{key: key}
}
func sign(data []byte, secretKey []byte) []byte {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(data)
	return mac.Sum(nil)
}

func (c *Codec) Encode(cur *Cursor) (string, error) {
	if cur == nil {
		return "", nil
	}
	payload, err := json.Marshal(cur)
	if err != nil {
		return "", err
	}
	raw := base32.StdEncoding.EncodeToString(payload)
	signature := string(sign(payload, c.key))
	return fmt.Sprintf("%s.%s", raw, signature), nil

}
func (c *Codec) Decode(token string) (*Cursor, error) {
	if token == "" {
		return nil, nil
	}
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, ErrInvalidTokenFormat
	}
	payload, err := base32.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	if !hmac.Equal(sign(payload, c.key), []byte(parts[1])) {
		return nil, ErrMismatchedHMC
	}
	var cursor Cursor
	err = json.Unmarshal(payload, &cursor)
	if err != nil {
		return nil, err
	}
	return &cursor, nil

}
