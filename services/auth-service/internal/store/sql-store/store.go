package store

import (
	"context"
	"database/sql"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
)

type SqlStorage struct {
	db *sql.DB
}

func NewSQLStorage(db *sql.DB) store.Storage {
	return &SqlStorage{
		db,
	}

}
func (s *SqlStorage) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)

}
func (s *SqlStorage) Users() store.Users {
	return &SQLUserStore{
		db: s.db,
	}

}
func (s *SqlStorage) Roles() store.Roles {
	return &SQLRoleStore{
		db: s.db,
	}

}
func (s *SqlStorage) Tokens() store.Tokens {
	return &SQLTokenStore{
		db: s.db,
	}

}

// - Identified tables - token, user, role, user_role (all tables have createdAt & updatedAt)
