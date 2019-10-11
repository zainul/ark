package xtoken

import (
	"errors"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	secret string
	once   sync.Once
)

// MyClaims ...
type MyClaims struct {
	UserID    int64     `json:"user_id"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
	jwt.StandardClaims
}

// NewSecret ...
func NewSecret(secret string) {
	once.Do(func() {
		secret = secret
	})
}

// Claim ...
func Claim(token string) (*MyClaims, error) {
	if secret == "" {
		return nil, errors.New("secret not provided, please initiate")
	}

	cl := MyClaims{}

	tokenClaim, err := jwt.ParseWithClaims(token, &cl, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("Token not valid")
	}

	if claims, ok := tokenClaim.Claims.(*MyClaims); ok && tokenClaim.Valid {
		return claims, nil
	}

	return nil, errors.New("Token not valid")
}
