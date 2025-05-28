package orders

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

type SqlOrderRepo struct {
	db *sql.DB
}

func NewSqlOrderRepo(db *sql.DB) *SqlOrderRepo {
	return &SqlOrderRepo{db}
}

func (o *SqlOrderRepo) GetOrders(pagination *types.PaginationPayload, filter OrderFilter) ([]Order, int, error) {
	var (
		args   []interface{}
		orders []Order
		where  []string
	)

	query := `SELECT id, quantity, unitPrice, status, productId, storeId, fulfillingInventoryId, sharingFormulaId, createdAt, updatedAt
			  FROM orders`
	countQuery := `SELECT COUNT(*) FROM orders`

	if filter.ProductId != 0 {
		where = append(where, "productId = ?")
		args = append(args, filter.ProductId)
	}
	if filter.StoreId != 0 {
		where = append(where, "storeId = ?")
		args = append(args, filter.StoreId)
	}
	if filter.Status != "" {
		where = append(where, "status = ?")
		args = append(args, filter.Status)
	}

	if len(where) > 0 {
		whereClause := " WHERE " + strings.Join(where, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	query += " ORDER BY createdAt DESC LIMIT ? OFFSET ?"
	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := o.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID,
			&order.Quantity,
			&order.UnitPrice,
			&order.Status,
			&order.ProductId,
			&order.StoreId,
			&order.FulfillingInventoryId,
			&order.Sharing.Id,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, order)
	}

	var total int
	err = o.db.QueryRow(countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (o *SqlOrderRepo) GetOrderById(orderId int) (Order, error) {
	query := `SELECT id, quantity, unitPrice, status, productId, storeId, fulfillingInventoryId, sharingFormulaId, createdAt, updatedAt 
			  FROM orders WHERE id = ?`

	var order Order
	err := o.db.QueryRow(query, orderId).Scan(
		&order.ID,
		&order.Quantity,
		&order.UnitPrice,
		&order.Status,
		&order.ProductId,
		&order.StoreId,
		&order.FulfillingInventoryId,
		&order.Sharing.Id,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, fmt.Errorf("order not found")
		}
		return order, err
	}

	return order, nil
}

func (o *SqlOrderRepo) UpdateOrder(orderId int, payload Order) (Order, error) {
	query := `UPDATE orders SET 
		quantity = ?, 
		unitPrice = ?, 
		status = ?, 
		productId = ?, 
		storeId = ?, 
		fulfillingInventoryId = ?, 
		sharingFormulaId = ?, 
		updatedAt = ? 
		WHERE id = ?`

	_, err := o.db.Exec(query,
		payload.Quantity,
		payload.UnitPrice,
		payload.Status,
		payload.ProductId,
		payload.StoreId,
		payload.FulfillingInventoryId,
		payload.Sharing.Id,
		time.Now(),
		orderId,
	)

	if err != nil {
		return Order{}, err
	}

	return o.GetOrderById(orderId)
}

func (o *SqlOrderRepo) BulkUpdateStatus(orderIds []int, status OrderStatus) ([]Order, error) {
	if len(orderIds) == 0 {
		return nil, nil
	}

	placeholders := strings.Repeat("?,", len(orderIds))
	placeholders = strings.TrimSuffix(placeholders, ",")

	query := fmt.Sprintf("UPDATE orders SET status = ?, updatedAt = ? WHERE id IN (%s)", placeholders)

	args := make([]interface{}, 0, len(orderIds)+2)
	args = append(args, status, time.Now())
	for _, id := range orderIds {
		args = append(args, id)
	}

	_, err := o.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	var updatedOrders []Order
	for _, id := range orderIds {
		order, err := o.GetOrderById(id)
		if err != nil {
			return nil, err
		}
		updatedOrders = append(updatedOrders, order)
	}

	return updatedOrders, nil
}

func (o *SqlOrderRepo) CreateOrder(payload Order) (Order, error) {
	query := `INSERT INTO orders (
		quantity, 
		unitPrice, 
		status, 
		productId, 
		storeId, 
		fulfillingInventoryId, 
		sharingFormulaId, 
		createdAt, 
		updatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	res, err := o.db.Exec(
		query,
		payload.Quantity,
		payload.UnitPrice,
		payload.Status,
		payload.ProductId,
		payload.StoreId,
		payload.FulfillingInventoryId,
		payload.Sharing.Id,
		now,
		now,
	)
	if err != nil {
		return Order{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Order{}, err
	}

	return o.GetOrderById(int(id))
}
