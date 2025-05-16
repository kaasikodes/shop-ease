package seller

import "database/sql"

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
func (r *SqlSellerRepo) CreateVendor(payload Seller) (*Seller, error)
func (r *SqlSellerRepo) GetVendor(sellerId int64) (*Seller, error)
