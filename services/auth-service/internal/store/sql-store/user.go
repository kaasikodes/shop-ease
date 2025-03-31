package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
)

type UserRoleStatus int

type UserRole = store.UserRole
type UserFilterQuery = store.UserFilterQuery
type User = store.User

type SQLUserStore struct {
	db *sql.DB
}

var (
	QueryTimeoutDuration = store.QueryTimeoutDuration
)

func (u *SQLUserStore) CreateWithVerificationToken(ctx context.Context, user *User, tokenValue string, tokenIsValidFor time.Duration) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = u.Create(ctx, tx, user, &store.UserRole{
		IsActive: true,
		ID:       store.CustomerID,
		Name:     store.Customer,
	})
	if err != nil {
		return err
	}
	err = createToken(ctx, tx, &store.Token{
		EntityId:  user.ID,
		TokenType: store.VerificationTokenType,
		ExpiresAt: time.Now().Add(tokenIsValidFor),
		Value:     tokenValue,
	})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func (u *SQLUserStore) RemoveMultipleUsers(ctx context.Context, tx *sql.Tx, emails []string) error {
	if len(emails) == 0 {
		return nil // No emails provided, nothing to remove
	}

	// Dynamically construct the placeholders for the IN clause
	placeholders := make([]string, len(emails))
	args := make([]interface{}, len(emails)) //Dev Comment: interface is used here purely to satisfy the ...any fn interface of ExecContext, if not string[] could have been used, any yes you can use any - but I'm not a fan!
	for i, email := range emails {
		placeholders[i] = "?" // MySQL uses "?" as a placeholder
		args[i] = email
	}

	query := fmt.Sprintf(`DELETE FROM users WHERE email IN (%s)`, strings.Join(placeholders, ","))

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// Handle transaction
	usesDefaultTx := false
	var err error
	if tx == nil {
		tx, err = u.db.BeginTx(ctx, nil)
		usesDefaultTx = true
		if err != nil {
			return err
		}
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Execute deletion query
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if usesDefaultTx {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: Update the user return types to populate createdAt, and updatedAt, as well as verifiedAt and fill where necessary
// Create inserts a new user and role for user, takes in transaction because it will probably have to be done as a single unit of work when creating token to be used to verify
func (u *SQLUserStore) Create(ctx context.Context, tx *sql.Tx, user *User, role *store.UserRole) error {
	// Queries
	queryUser := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	queryRole := `INSERT INTO userRoles (userId, roleId, isActive) VALUES (?, ?, ?)`
	queryGetUser := `SELECT id, createdAt, updatedAt FROM users WHERE id = LAST_INSERT_ID()`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// Handle transaction
	usesDefaultTx := false
	var err error
	if tx == nil {
		tx, err = u.db.BeginTx(ctx, nil)
		usesDefaultTx = true
		if err != nil {
			return err
		}
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert user
	result, err := tx.ExecContext(ctx, queryUser, user.Name, user.Email, user.Password.GetHash())
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // Duplicate entry
				return store.ErrDuplicateEmail
			}
		}
		return err
	}

	// Get the last inserted user ID
	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(userID)

	// Get createdAt and updatedAt
	err = tx.QueryRowContext(ctx, queryGetUser).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert role
	_, err = tx.ExecContext(ctx, queryRole, user.ID, role.ID, role.IsActive)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // Duplicate entry for user-role relation
				return store.ErrDuplicateUserRole
			}
		}
		return err
	}

	if usesDefaultTx {
		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	user.Roles = append(user.Roles, *role)
	return nil
}

// Verifying a user
func (u *SQLUserStore) Verify(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET isVerified = ? WHERE email = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var err error
	usesDefaultTx := false

	if tx == nil {
		tx, err = u.db.BeginTx(ctx, nil)
		usesDefaultTx = true
		if err != nil {
			return err
		}
	}

	// Execute update
	result, err := tx.ExecContext(ctx, query, true, user.Email)
	if err != nil {
		return store.ErrVerifyUser
	}

	// Check if any row was updated
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return store.ErrNotFound
	}

	user.IsVerified = true

	if usesDefaultTx {
		return tx.Commit()
	}

	return nil
}

func (u *SQLUserStore) AssignRole(ctx context.Context, userId int, roleId DefaultRoleID) (*UserRole, error) {
	query := `
		INSERT INTO userRoles (userId, roleId)
		VALUES ($1, $2)
		RETURNING roleId, isActive, 
		(SELECT name FROM roles WHERE id = $2) AS roleName`
	role := &UserRole{}
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := u.db.QueryRowContext(ctx, query, userId, roleId).Scan(&role.ID, &role.IsActive, &role.Name)

	if err != nil { //TODO: Revisit all error to handle with clean error messages, also consider creating an errUtil that accepts errors and just matches them and returns clean err messages or return the err
		return nil, err
	}

	return role, nil

}
func (u *SQLUserStore) ActivateOrDeactivateRole(ctx context.Context, userId int, roleId DefaultRoleID, isActive bool) (*UserRole, error) {
	query := `
		UPDATE userRoles
		SET isActive = $3
		WHERE userId = $1 AND roleId = $2
		RETURNING roleId, isActive, 
		(SELECT name FROM roles WHERE id = $2) AS roleName`
	role := &UserRole{}
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := u.db.QueryRowContext(ctx, query, userId, roleId, isActive).Scan(&role.ID, &role.IsActive, &role.Name)

	if err != nil {
		return nil, err
	}

	return role, nil

}
func (u *SQLUserStore) Update(ctx context.Context, user *User) (*User, error) {
	query := `
		UPDATE users
		SET name = $3
		WHERE id = $1 OR email = $2
		RETURNING id, isVerified
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := u.db.QueryRowContext(ctx, query, user.ID, user.Email, user.Name).Scan(&user.ID, &user.IsVerified)

	if err != nil {
		return nil, err
	}

	return user, nil

}
func (u *SQLUserStore) GetByEmailOrId(ctx context.Context, user *User) (*User, error) {
	queryU := `SELECT (id, email, name, isVerified) FROM users WHERE id = $1 OR email = $2`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := u.db.QueryRowContext(ctx, queryU, user.ID, user.Email).Scan(&user.ID, &user.Email, &user.Name, &user.IsVerified)
	if err != nil {
		return nil, err
	}
	queryR := `
		SELECT ur.roleId, r.name, ur.isActive 
		FROM userRoles ur
		JOIN roles r ON ur.roleId = r.id
		WHERE ur.userId = $1
	`
	rows, err := u.db.QueryContext(ctx, queryR, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	user.Roles = []UserRole{}
	for rows.Next() {
		var role UserRole

		if err := rows.Scan(&role.ID, &role.Name, &role.IsActive); err != nil {
			return nil, err
		}
		user.Roles = append(user.Roles, role)

	}
	return user, nil
}
func (u *SQLUserStore) Get(ctx context.Context, pagination PaginationPayload, filter UserFilterQuery) ([]User, int, error) {
	var total int
	var users []User

	// Base query to get users, 1=1 is set as placeholder to allow for adding dynamic additions to the where clause down the line as it will always evaluate to true
	query := `SELECT u.id, u.name, u.email, u.isVerified, u.verifiedAt 
			  FROM users u 
			  WHERE 1=1 `

	// Conditions for filtering
	args := []interface{}{}
	argIndex := 1

	if filter.IsVerified {
		query += fmt.Sprintf(" AND u.isVerified = $%d", argIndex)
		args = append(args, filter.IsVerified)
		argIndex++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.email ILIKE $%d)", argIndex, argIndex+1)
		args = append(args, "%"+filter.Search+"%", "%"+filter.Search+"%")
		argIndex += 2
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY u.id DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pagination.Limit, pagination.Offset)

	// Execute query
	rows, err := u.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.IsVerified, &user.VerifiedAt)
		if err != nil {
			return nil, 0, err
		}

		// Fetch user roles
		roleQuery := `SELECT r.id, r.name, ur.isActive 
					  FROM userRoles ur 
					  JOIN roles r ON ur.roleId = r.id 
					  WHERE ur.userId = $1`
		roleRows, err := u.db.QueryContext(ctx, roleQuery, user.ID)
		if err != nil {
			return nil, 0, err
		}
		defer roleRows.Close()

		for roleRows.Next() {
			var role UserRole
			if err := roleRows.Scan(&role.ID, &role.Name, &role.IsActive); err != nil {
				return nil, 0, err
			}

			user.Roles = append(user.Roles, role)
		}

		users = append(users, user)
	}

	// Get total count for pagination
	countQuery := `SELECT COUNT(*) FROM users u WHERE 1=1`
	if filter.IsVerified {
		countQuery += " AND u.isVerified = true"
	}
	if filter.Search != "" {
		countQuery += " AND (u.name ILIKE $1 OR u.email ILIKE $1)"
	}

	err = u.db.QueryRowContext(ctx, countQuery, "%"+filter.Search+"%").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
