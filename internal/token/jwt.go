package token

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Abhishek48Shah/blog/internal/util"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
}
type IssuedToken struct {
	JWT      string
	Key      string
	ExpireAt time.Time
}

func NewTokenService(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey, issuer string) *TokenService {
	return &TokenService{privateKey: privKey, publicKey: pubKey, issuer: issuer}
}
func (t *TokenService) GenerateNewToken() (*IssuedToken, error) {
	randByt := make([]byte, 32)
	if _, err := rand.Read(randByt); err != nil {
		return nil, util.InternalServerErr("somthing went wrong", fmt.Errorf("reading random bytes: %w", err))
	}
	key := hex.EncodeToString(randByt)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    t.issuer,
		ID:        key,
	})
	token, err := claims.SignedString(t.privateKey)
	if err != nil {
		return nil, util.InternalServerErr("somthings went wrong", fmt.Errorf("signing claims with private key: %w", err))
	}
	return &IssuedToken{JWT: token, ExpireAt: expiresAt, Key: key}, nil
}

func (t *TokenService) VerifyToken(token string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(tn *jwt.Token) (any, error) {
		if _, ok := tn.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %s", tn.Method)
		}
		return t.publicKey, nil
	}, jwt.WithIssuer(t.issuer), jwt.WithExpirationRequired())
	if err != nil {
		return nil, fmt.Errorf("parse claims: %w", err)
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
