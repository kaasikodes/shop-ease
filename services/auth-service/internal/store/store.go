package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Storage struct {
	Users interface {
	    RemoveMultipleUsers(ctx context.Context, tx *sql.Tx, emails []string) error
		Create(context.Context, *sql.Tx, *User, *UserRole) error
		Verify(context.Context, *sql.Tx, *User,) error
		AssignRole(ctx context.Context, userId int, roleId DefaultRoleID) (*UserRole, error)
		ActivateOrDeactivateRole(ctx context.Context, userId int, roleId DefaultRoleID, isActive bool) (*UserRole, error)
		Update(context.Context,  *User) ( *User, error )
		GetByEmailOrId(context.Context, *User) ( *User, error )
		Get(ctx context.Context, pagination PaginationPayload, filter UserFilterQuery) ( []User, int, error)
	}
	Tokens interface {
		Create(context.Context, *sql.Tx, *Token) error
		Remove(context.Context,  *Token) error
		GetOne(ctx context.Context, token *Token) (*Token, error)
	}
	Roles interface {
		CreateDefaultRoles(ctx context.Context) ([]Role, error)
		GetByName(context.Context,  DefaultRoleName) (*Role, error)
		Get(ctx context.Context, pagination PaginationPayload, status string) ( []Role, error,)
	}
}
var (
	QueryTimeoutDuration = 4 * time.Second
	ErrNotFound = errors.New("entity does not exits")
	ErrConflict = errors.New("entity already exists")
)

func NewStorage(db *sql.DB) Storage  {
	return Storage{
		Users: &UserStore{db},
		Roles: &RoleStore{db},
		Tokens: &TokenStore{db},

	}
	
}
// - Identified tables - token, user, role, user_role (all tables have createdAt & updatedAt)
