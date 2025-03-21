package store

import (
	"context"
	"database/sql"
	"errors"
)

// Define DefaultRoleName & DefaultRoleID enum
type DefaultRoleName string
type DefaultRoleID int

const (
	Admin DefaultRoleName = "admin"
	Vendor DefaultRoleName  = "vendor"
	Customer DefaultRoleName = "customer"
	


)
const (
	AdminID DefaultRoleID = 1
	VendorID DefaultRoleID = 2
	CustomerID DefaultRoleID = 3
)

var (
	ErrDefaultRolesAlreadyExists = errors.New("default role already exists")
)



type Role struct {
	ID        DefaultRoleID `json:"id"`
	Name      DefaultRoleName `json:"name"`
	IsDefault bool `json:"isDefault"`
	Common
}

type RoleStore struct {
	db *sql.DB
}

var DefaultRoles = []Role{
	{ID: AdminID, Name: Admin, IsDefault: true},
	{ID: VendorID, Name: Vendor, IsDefault: true},
	{ID: CustomerID, Name: Customer, IsDefault: true},
}
func (r *RoleStore) CreateDefaultRoles(ctx context.Context) ([]Role, error) {
	

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
func (r *RoleStore) GetByName(ctx context.Context,  name DefaultRoleName) (*Role, error) {
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
func (r *RoleStore)  Get(ctx context.Context, pagination PaginationPayload, status string) ( []Role, error,) {
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
	if err := rows.Err(); err !=nil {
		return nil, err
	}

	return roles, nil
	
}