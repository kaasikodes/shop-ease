package store

import (
	"database/sql"
	"errors"
)

// Define DefaultRoleName & DefaultRoleID enum
type DefaultRoleName string
type DefaultRoleID int

const (
	Admin    DefaultRoleName = "admin"
	Vendor   DefaultRoleName = "vendor"
	Customer DefaultRoleName = "customer"
)
const (
	AdminID    DefaultRoleID = 1
	VendorID   DefaultRoleID = 2
	CustomerID DefaultRoleID = 3
)

var (
	ErrDefaultRolesAlreadyExists = errors.New("default role already exists")
)

type Role struct {
	ID        DefaultRoleID   `json:"id"`
	Name      DefaultRoleName `json:"name"`
	IsDefault bool            `json:"isDefault"`
	Common
}

type SQLRoleStore struct {
	db *sql.DB
}

var DefaultRoles = []Role{
	{ID: AdminID, Name: Admin, IsDefault: true},
	{ID: VendorID, Name: Vendor, IsDefault: true},
	{ID: CustomerID, Name: Customer, IsDefault: true},
}
