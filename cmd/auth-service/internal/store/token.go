package store

import (
	"errors"
	"time"
)

type TokenType string

const (
	VerificationTokenType  TokenType = "VERIFICATION"
	PasswordResetTokenType TokenType = "PASSWORD_RESET"
	AccessTokenType        TokenType = "ACCESS_TOKEN"
	RefreshTokenType       TokenType = "REFRESH_TOKEN"
)

var (
	ErrNoTokenFound error = errors.New("no token found")
)

type Token struct {
	Id        int       `json:"id"`
	EntityId  int       `json:"entityId"`
	TokenType TokenType `json:"type"`
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expiresAt"`
	Common
}
