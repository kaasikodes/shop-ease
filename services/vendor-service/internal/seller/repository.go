package seller

import "database/sql"

type SellerRepo interface {
	// create vendor
	CreateVendor(payload Seller) (*Seller, error)
}

type SqlSellerRepo struct {
	db *sql.DB
}

func (r *SqlSellerRepo) CreateVendor(payload Seller) (*Seller, error)

func NewSqlSellerRepo(db *sql.DB) *SqlSellerRepo {
	return &SqlSellerRepo{db}

}
