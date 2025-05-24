package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Users interface {
	CreateWithVerificationToken(ctx context.Context, user *User, tokenValue string, tokenIsValidFor time.Duration) error
	RemoveMultipleUsers(ctx context.Context, tx *sql.Tx, emails []string) error
	Create(context.Context, *sql.Tx, *User, *UserRole) error
	Verify(context.Context, *sql.Tx, *User) error
	AssignRole(ctx context.Context, userId int, roleId DefaultRoleID) (*UserRole, error)
	ActivateOrDeactivateRole(ctx context.Context, userId int, roleId DefaultRoleID, isActive bool) (*UserRole, error)
	Update(context.Context, *User) (*User, error)
	GetByEmailOrId(context.Context, *User) (*User, error)
	Get(ctx context.Context, pagination PaginationPayload, filter UserFilterQuery) ([]User, int, error)
}
type Tokens interface {
	Create(context.Context, *sql.Tx, *Token) error
	Remove(context.Context, *Token) error
	GetOne(ctx context.Context, value string, entityId int, tokenType TokenType) (*Token, error)
}
type Roles interface {
	CreateDefaultRoles(ctx context.Context) ([]Role, error)
	GetByName(context.Context, DefaultRoleName) (*Role, error)
	Get(ctx context.Context, pagination PaginationPayload, status string) ([]Role, error)
}
type Storage interface {
	Users() Users
	Tokens() Tokens

	Roles() Roles
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

var (
	QueryTimeoutDuration = 4 * time.Second
	ErrNotFound          = errors.New("entity does not exists")
	ErrConflict          = errors.New("entity already exists")
)

// - Identified tables - token, user, role, user_role (all tables have createdAt & updatedAt)
