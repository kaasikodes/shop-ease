package store

import (
	"context"
	"database/sql"
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
 

type TokenStore struct {
	db *sql.DB
}

// Create inserts a new token into the database
func (t *TokenStore) Create(ctx context.Context, tx *sql.Tx, token *Token) error {
	query := `
		INSERT INTO tokens (entityId, tokenType, value, expiresAt)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := tx.QueryRowContext(ctx, query, token.EntityId, token.TokenType, token.Value, token.ExpiresAt).
		Scan(&token.Id)
	if err != nil {
		return err
	}
	return nil
}

// Remove deletes a token from the database
func (t *TokenStore) Remove(ctx context.Context, token *Token) error {
	query := `
		DELETE FROM tokens
		WHERE id = $1 OR (entityId = $2 AND tokenType = $3)
	`
	_, err := t.db.ExecContext(ctx, query, token.Id, token.EntityId, token.TokenType)
	if err != nil {
		return err
	}
	return nil
}

// GetOne retrieves a token based on ID or entityId and tokenType
func (t *TokenStore) GetOne(ctx context.Context, token *Token) (*Token, error) {
	query := `
		SELECT id, entityId, tokenType, value, expiresAt
		FROM tokens
		WHERE id = $1 OR (entityId = $2 AND tokenType = $3)
		LIMIT 1
	`
	err := t.db.QueryRowContext(ctx, query, token.Id, token.EntityId, token.TokenType).Scan(&token.Id, &token.EntityId, &token.TokenType, &token.Value, &token.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoTokenFound // No token found
		}
		return nil, err
	}
	return token, nil
}
