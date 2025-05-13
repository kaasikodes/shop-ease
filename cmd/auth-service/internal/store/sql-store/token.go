package store

import (
	"context"
	"database/sql"
	"log"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
)

type TokenType = store.TokenType

const (
	VerificationTokenType  = store.VerificationTokenType
	PasswordResetTokenType = store.PasswordResetTokenType
	AccessTokenType        = store.AccessTokenType
	RefreshTokenType       = store.RefreshTokenType
)

var (
	ErrNoTokenFound = store.ErrNoTokenFound
)

type Token = store.Token

type SQLTokenStore struct {
	db *sql.DB
}

// Create inserts a new token into the database
func (t *SQLTokenStore) Create(ctx context.Context, tx *sql.Tx, token *Token) error {
	return createToken(ctx, tx, token)
}
func createToken(ctx context.Context, tx *sql.Tx, token *Token) error {
	query := `
		INSERT INTO tokens (entityId, tokenType, value, expiresAt)
		VALUES (?, ?, ?, ?)
	`
	log.Println("expires at ...", token.ExpiresAt)
	result, err := tx.ExecContext(ctx, query, token.EntityId, token.TokenType, token.Value, token.ExpiresAt)
	if err != nil {
		return err
	}

	// Get last inserted ID (since MySQL does not support RETURNING)
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	token.Id = int(lastID)

	return nil
}

// Remove deletes a token from the database
func (t *SQLTokenStore) Remove(ctx context.Context, token *Token) error {
	query := `
		DELETE FROM tokens
		WHERE id = ? OR (entityId = ? AND tokenType = ?)
	`
	_, err := t.db.ExecContext(ctx, query, token.Id, token.EntityId, token.TokenType)
	if err != nil {
		return err
	}
	return nil
}

// GetOne retrieves a token based on ID or entityId and tokenType
func (t *SQLTokenStore) GetOne(ctx context.Context, token *Token) (*Token, error) {
	query := `
		SELECT id, entityId, tokenType, value, expiresAt
		FROM tokens
		WHERE id = ? OR (entityId = ? AND tokenType = ?)
		LIMIT 1
	`
	err := t.db.QueryRowContext(ctx, query, token.Id, token.EntityId, token.TokenType).
		Scan(&token.Id, &token.EntityId, &token.TokenType, &token.Value, &token.ExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoTokenFound // No token found
		}
		return nil, err
	}
	return token, nil
}
