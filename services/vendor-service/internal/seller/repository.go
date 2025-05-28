package seller

import (
	"database/sql"
	"fmt"
)

type SellerRepo interface {
	// create vendor
	CreateVendor(payload Seller) (*Seller, error)
	GetVendor(sellerId int64) (*Seller, error)
}

type SqlSellerRepo struct {
	db *sql.DB
}

func NewSqlSellerRepo(db *sql.DB) *SqlSellerRepo {
	return &SqlSellerRepo{db}

}

// CreateVendor inserts a new seller into the sellers table and returns the inserted record
func (r *SqlSellerRepo) CreateVendor(payload Seller) (*Seller, error) {
	query := `
		INSERT INTO sellers (userId, name, email, phone)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, payload.UserId, payload.Name, payload.Email, payload.Phone)
	if err != nil {
		return nil, fmt.Errorf("error creating vendor: %w", err)
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last inserted ID: %w", err)
	}

	return r.GetVendor(insertedID)
}

// GetVendor fetches a seller by their ID from the sellers table
func (r *SqlSellerRepo) GetVendor(sellerId int64) (*Seller, error) {
	query := `
		SELECT id, userId, name, email, phone, createdAt, updatedAt
		FROM sellers
		WHERE id = ?
	`

	var seller Seller
	err := r.db.QueryRow(query, sellerId).
		Scan(&seller.ID, &seller.UserId, &seller.Name, &seller.Email, &seller.Phone, &seller.CreatedAt, &seller.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No vendor found
		}
		return nil, fmt.Errorf("error fetching vendor: %w", err)
	}

	return &seller, nil
}
