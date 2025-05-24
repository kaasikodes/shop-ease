package jwttoken

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email,omitempty"` // Optional field
	jwt.RegisteredClaims
}

type JwtMaker struct {
	secretKey string
}

func NewJwtMaker(secret string) *JwtMaker {
	return &JwtMaker{secretKey: secret}
}

var (
	ErrExpiredToken = errors.New("expired token")
	ErrInvalidToken = errors.New("invalid token")
	ErrNoAuthHeader = errors.New("authorization header is missing")
	ErrWrongFormat  = errors.New("authorization header format must be Bearer {token}")
)

// CreateToken generates a JWT signed with HS256
func (j *JwtMaker) CreateToken(userID, userEmail string, duration time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Email:  userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// VerifyToken parses and validates the JWT token
func (j *JwtMaker) VerifyToken(tokenStr string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

// ExtractToken extracts token from Authorization header
func (j *JwtMaker) ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrWrongFormat
	}

	return parts[1], nil
}

// ExtractAndVerifyToken extracts the token from the request and verifies it
func (j *JwtMaker) ExtractAndVerifyToken(r *http.Request) (*CustomClaims, error) {
	tokenStr, err := j.ExtractToken(r)
	if err != nil {
		return nil, err
	}
	return j.VerifyToken(tokenStr)
}
