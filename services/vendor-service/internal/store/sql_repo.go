package store

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
	"github.com/kaasikodes/shop-ease/shared/utils"
)

type SqlStoreRepo struct {
	db *sql.DB
}

func NewSqlStoreRepo(db *sql.DB) *SqlStoreRepo {
	return &SqlStoreRepo{db}

}

func (r *SqlStoreRepo) CreateStore(payload Store) (*Store, error) {
	query := `
		INSERT INTO stores (vendorId, name, description, location, lat, long, country, state, lga, landmark, timezone, postalCode, phone, email, bank, number, swiftCode)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		payload.VendorId, payload.Name, payload.Description,
		payload.Address.Location, payload.Address.Lat, payload.Address.Long, payload.Address.Country,
		payload.Address.State, payload.Address.Lga, payload.Address.Landmark, payload.Address.Timezone, payload.Address.PostalCode,
		payload.Contact.Phone, payload.Contact.Email,
		payload.Account.Bank, payload.Account.Number, payload.Account.SwiftCode,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating store: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	return r.GetStoreById(id)
}

func (r *SqlStoreRepo) UpdateStore(id int64, payload Store) (*Store, error) {
	query := `
		UPDATE stores
		SET vendorId = ?, name = ?, description = ?, location = ?, lat = ?, long = ?, country = ?, state = ?, lga = ?, landmark = ?, timezone = ?, postalCode = ?, phone = ?, email = ?, bank = ?, number = ?, swiftCode = ?, updatedAt = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(query,
		payload.VendorId, payload.Name, payload.Description,
		payload.Address.Location, payload.Address.Lat, payload.Address.Long, payload.Address.Country,
		payload.Address.State, payload.Address.Lga, payload.Address.Landmark, payload.Address.Timezone, payload.Address.PostalCode,
		payload.Contact.Phone, payload.Contact.Email,
		payload.Account.Bank, payload.Account.Number, payload.Account.SwiftCode,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating store: %w", err)
	}

	return r.GetStoreById(id)
}
func (r *SqlStoreRepo) GetStoreById(id int64) (*Store, error) {
	query := `
		SELECT id, vendorId, name, description, location, lat, long, country, state, lga, landmark, timezone, postalCode, phone, email, bank, number, swiftCode, createdAt, updatedAt
		FROM stores
		WHERE id = ?
	`
	var store Store
	err := r.db.QueryRow(query, id).Scan(
		&store.ID, &store.VendorId, &store.Name, &store.Description,
		&store.Address.Location, &store.Address.Lat, &store.Address.Long, &store.Address.Country,
		&store.Address.State, &store.Address.Lga, &store.Address.Landmark, &store.Address.Timezone, &store.Address.PostalCode,
		&store.Contact.Phone, &store.Contact.Email,
		&store.Account.Bank, &store.Account.Number, &store.Account.SwiftCode,
		&store.CreatedAt, &store.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Store not found
		}
		return nil, fmt.Errorf("error fetching store: %w", err)
	}
	return &store, nil
}

func (r *SqlStoreRepo) GetProducts(pagination *utils.PaginationPayload, filter *types.ProductFilter) ([]types.Product, int, error) {
	var (
		products []types.Product
		args     []interface{}
	)

	query := `
		SELECT p.id, p.name, p.description, p.amount, p.createdAt, p.updatedAt
		FROM products p
		JOIN inventories i ON p.id = i.productId
		WHERE 1=1
	`
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		JOIN inventories i ON p.id = i.productId
		WHERE 1=1
	`

	if filter != nil && filter.StoreId != 0 {
		query += " AND i.storeId = ?"
		countQuery += " AND i.storeId = ?"
		args = append(args, filter.StoreId)
	}

	query += " GROUP BY p.id LIMIT ? OFFSET ?"
	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product types.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price.Amount, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, product)
	}

	var total int
	err = r.db.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting products: %w", err)
	}

	return products, total, nil
}

func (r *SqlStoreRepo) BulkAddInventory(payload []Inventory) error {
	if len(payload) == 0 {
		return nil
	}

	query := `
		INSERT INTO inventories (quantity, unitCostPrice, productId, storeId, arrivalOrProduceDate, createdAt, updatedAt)
		VALUES
	`
	args := []interface{}{}
	valueStrings := []string{}

	for _, inv := range payload {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, NOW(), NOW())")
		args = append(args, inv.Quantity, inv.UnitCostPrice, inv.ProductId, inv.StoreId, inv.ArrivalorProduceDate)
	}

	query += strings.Join(valueStrings, ",")
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error bulk inserting inventory: %w", err)
	}

	return nil
}

func (r *SqlStoreRepo) UpdateInventory(id int64, payload Inventory) error {
	query := `
		UPDATE inventories
		SET quantity = ?, unitCostPrice = ?, productId = ?, storeId = ?, arrivalOrProduceDate = ?, updatedAt = NOW()
		WHERE id = ?
	`
	_, err := r.db.Exec(query, payload.Quantity, payload.UnitCostPrice, payload.ProductId, payload.StoreId, payload.ArrivalorProduceDate, id)
	if err != nil {
		return fmt.Errorf("error updating inventory: %w", err)
	}
	return nil
}

func (r *SqlStoreRepo) DeleteInventory(id int64) (*int64, error) {

	query := `DELETE FROM inventories WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return nil, fmt.Errorf("error deleting inventory: %w", err)
	}

	return &id, nil
}
func (r *SqlStoreRepo) GetInventories(pagination *utils.PaginationPayload, filter *types.InventoryFilter) ([]Inventory, int, error) {
	var (
		inventories  []Inventory
		args         []interface{}
		whereClauses []string
	)

	whereClauses = append(whereClauses, "1=1") // always true for building WHERE clause dynamically

	if filter != nil {
		if filter.ProductId != 0 {
			whereClauses = append(whereClauses, "productId = ?")
			args = append(args, filter.ProductId)
		}
		if filter.StoreId != 0 {
			whereClauses = append(whereClauses, "storeId = ?")
			args = append(args, filter.StoreId)
		}
	}

	whereSQL := strings.Join(whereClauses, " AND ")
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM inventories WHERE %s", whereSQL)
	dataQuery := fmt.Sprintf(`
		SELECT id, quantity, unitCostPrice, productId, storeId, arrivalOrProduceDate, createdAt, updatedAt
		FROM inventories
		WHERE %s
		LIMIT ? OFFSET ?`, whereSQL)

	// Count total
	var total int
	countArgs := args
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting inventories: %w", err)
	}

	// Add pagination to args
	args = append(args, pagination.Limit, pagination.Offset)

	// Fetch rows
	rows, err := r.db.Query(dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error fetching inventories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var inv Inventory
		err := rows.Scan(
			&inv.Id, &inv.Quantity, &inv.UnitCostPrice, &inv.ProductId,
			&inv.StoreId, &inv.ArrivalorProduceDate, &inv.CreatedAt, &inv.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning inventory: %w", err)
		}
		inventories = append(inventories, inv)
	}

	return inventories, total, nil
}
