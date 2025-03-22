package store

import (
	"context"
	"database/sql"
	"time"
)

// TODO: Refactor to match sql-store, also do more research as to whether should be exported from here or the nested export is ok

type TestStorage struct {
	file string
}

func (s *TestStorage) Users() Users {
	return &UserStoreR{}

}
func (s *TestStorage) Tokens() Tokens {

	return &TokenStoreR{}

}
func (s *TestStorage) Roles() Roles {
	return &RoleStoreR{}

}

func NewTestStorage(logRecord string) Storage {
	return &TestStorage{
		file: logRecord,
	}

}

type UserStoreR struct {
}
type RoleStoreR struct {
}
type TokenStoreR struct {
}

func (r *UserStoreR) CreateWithVerificationToken(ctx context.Context, user *User, tokenValue string, tokenIsValidFor time.Duration) error {
	return nil
}

func (r *UserStoreR) RemoveMultipleUsers(ctx context.Context, tx *sql.Tx, emails []string) error {
	return nil
}
func (r *UserStoreR) Create(context.Context, *sql.Tx, *User, *UserRole) error {
	return nil
}
func (r *UserStoreR) Verify(context.Context, *sql.Tx, *User) error {
	return nil
}
func (r *UserStoreR) AssignRole(ctx context.Context, userId int, roleId DefaultRoleID) (*UserRole, error) {
	return nil, nil
}
func (r *UserStoreR) ActivateOrDeactivateRole(ctx context.Context, userId int, roleId DefaultRoleID, isActive bool) (*UserRole, error) {
	return nil, nil
}
func (r *UserStoreR) Update(context.Context, *User) (*User, error) {
	return nil, nil
}
func (r *UserStoreR) GetByEmailOrId(context.Context, *User) (*User, error) {
	return nil, nil
}
func (r *UserStoreR) Get(ctx context.Context, pagination PaginationPayload, filter UserFilterQuery) ([]User, int, error) {
	return nil, 0, nil
}

func (t *TokenStoreR) Create(context.Context, *sql.Tx, *Token) error {
	return nil
}
func (t *TokenStoreR) Remove(context.Context, *Token) error {
	return nil
}
func (t *TokenStoreR) GetOne(ctx context.Context, token *Token) (*Token, error) {
	return nil, nil
}

func (t *RoleStoreR) CreateDefaultRoles(ctx context.Context) ([]Role, error) {
	return nil, nil
}
func (t *RoleStoreR) GetByName(context.Context, DefaultRoleName) (*Role, error) {
	return nil, nil
}
func (t *RoleStoreR) Get(ctx context.Context, pagination PaginationPayload, status string) ([]Role, error) {
	return nil, nil
}
