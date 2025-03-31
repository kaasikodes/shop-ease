package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
)

// Define DefaultRoleName & DefaultRoleID enum
type DefaultRoleName = store.DefaultRoleName
type DefaultRoleID = store.DefaultRoleID
type PaginationPayload = store.PaginationPayload

const (
	Admin    = store.Admin
	Vendor   = store.Vendor
	Customer = store.Customer
)
const (
	AdminID    = store.AdminID
	VendorID   = store.VendorID
	CustomerID = store.CustomerID
)

var (
	ErrDefaultRolesAlreadyExists = store.ErrDefaultRolesAlreadyExists
	ErrNotFound                  = store.ErrNotFound
)

type Role = store.Role

type SQLRoleStore struct {
	db *sql.DB
}

var DefaultRoles = store.DefaultRoles

func (r *SQLRoleStore) CreateDefaultRoles(ctx context.Context) ([]Role, error) {

	// Insert roles if they don't exist
	insertQuery := `
		INSERT INTO roles (id, name, isDefault)
        VALUES (?, ?, ?)
        ON DUPLICATE KEY UPDATE id = id;
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, role := range DefaultRoles {
		_, err := tx.ExecContext(ctx, insertQuery, role.ID, role.Name, role.IsDefault)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Retrieve all default roles from the database
	query := `SELECT id, name, isDefault FROM roles WHERE isDefault = TRUE ORDER BY id;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name, &role.IsDefault); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
func (r *SQLRoleStore) GetByName(ctx context.Context, name DefaultRoleName) (*Role, error) {
	query := `
		SELECT id, name, isDefault
		FROM roles
		WHERE name = $1
		LIMIT 1;
	`
	role := &Role{}

	err := r.db.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name, &role.IsDefault)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound

		}
		return nil, err
	}

	return role, err

}
func (r *SQLRoleStore) Get(ctx context.Context, pagination PaginationPayload, status string) ([]Role, error) {
	query := `
		SELECT id, name, isDefault
		FROM roles
		WHERE status = $1
		LIMIT $2 OFFSET $3;
	`
	var roles []Role
	rows, err := r.db.QueryContext(ctx, query, status, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var role Role
		err := rows.Scan(&role.ID, &role.Name, &role.IsDefault)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil

}
