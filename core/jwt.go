package core

import (
	"crypto/rand"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents JWT claims with generic data payload
type JWTClaims[T any] struct {
	jwt.RegisteredClaims
	Data *T `json:"data,omitempty"`
}

var jwtSecret = func() []byte {
	if secret := os.Getenv("STREMTHRU_JWT_SECRET"); secret != "" {
		return []byte(secret)
	}
	randomSecret := make([]byte, 32)
	if _, err := rand.Read(randomSecret); err != nil {
		panic("failed to generate JWT secret: " + err.Error())
	}
	return randomSecret
}()

// GenerateJWT creates signed JWT token with provided claims
func GenerateJWT[T any](claims JWTClaims[T]) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseJWT verifies and parses JWT token string
func ParseJWT[T any](tokenString string) (*JWTClaims[T], error) {
	claims := &JWTClaims[T]{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}