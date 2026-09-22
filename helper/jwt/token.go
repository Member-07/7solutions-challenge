package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager interface {
	Generate(userID string) (string, error)
	Parse(tokenString string) (string, error)
}
type jwtToken struct {
	secret []byte
	ttl    time.Duration
}

type claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJwtManager(secret string, ttl time.Duration) (JWTManager, error) {
	if secret == "" {
		return nil, fmt.Errorf("jwt secret is required")
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &jwtToken{secret: []byte(secret), ttl: ttl}, nil
}

func (s *jwtToken) Generate(userID string) (string, error) {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	})
	return token.SignedString(s.secret)
}

func (s *jwtToken) Parse(tokenString string) (string, error) {
	parsed, err := jwt.ParseWithClaims(tokenString, &claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid token")
		}
		return s.secret, nil
	})
	if err != nil || !parsed.Valid {
		return "", fmt.Errorf("invalid token")
	}

	c, ok := parsed.Claims.(*claims)
	if !ok || c.UserID == "" {
		return "", fmt.Errorf("invalid token")
	}
	return c.UserID, nil
}
