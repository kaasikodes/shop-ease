package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kaasikodes/shop-ease/services/order-service/internal/model"
	"github.com/kaasikodes/shop-ease/shared/utils"
)

type PostgresOrderRepo struct {
	db *sql.DB
}

func NewPostgresOrderRepo(db *sql.DB) *PostgresOrderRepo {
	return &PostgresOrderRepo{db}
}

func (r *PostgresOrderRepo) CreateOrder(ctx context.Context, userId int, items []CreateOrderInputItem) (*int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert Order
	var orderId int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, status, is_paid, is_canceled, created_at, updated_at)
		VALUES ($1, $2, false, false, NOW(), NOW())
		RETURNING id
	`, userId, string(model.UnpaidOrPendingOrderStatus)).Scan(&orderId)
	if err != nil {
		return nil, err
	}

	// Insert Order Items
	for _, item := range items {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO order_items (order_id, product_id, store_id, price, quantity,amount_to_be_paid, created_at, updated_at, status)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), $7)
		`, orderId, item.ProductId, item.StoreId, item.Price, item.Quantity, item.AmountToBePaid, model.UnpaidOrPendingOrderStatus)
		if err != nil {
			return nil, err
		}
	}

	return &orderId, tx.Commit()
}

func (r *PostgresOrderRepo) UpdateOrderStatus(ctx context.Context, orderId int, status model.OrderStatus) error {
	var args []interface{}
	args = append(args, string(status), time.Now(), orderId)

	query := `
		UPDATE orders
		SET status = $1, updated_at = $2`

	switch status {
	case model.CancelledOrderStatus:
		query += `, is_canceled = true, canceled_at = NOW()`
	case model.PaidOrderStatus:
		query += `, is_paid = true, paid_at = NOW()`
	}

	query += ` WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *PostgresOrderRepo) GetOrderById(ctx context.Context, orderId int) (model.Order, error) {
	var order model.Order

	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, status, is_paid, is_canceled, paid_at, canceled_at, created_at, updated_at
		FROM orders
		WHERE id = $1
	`, orderId).Scan(
		&order.Id,
		&order.UserId,
		&order.Status,
		&order.IsPaid,
		&order.IsCanceled,
		&order.PaidAt,
		&order.CanceledAt,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return order, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, product_id, store_id, quantity, created_at, updated_at
		FROM order_items
		WHERE order_id = $1
	`, orderId)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(
			&item.Id,
			&item.ProductId,
			&item.StoreId,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return order, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *PostgresOrderRepo) GetOrders(ctx context.Context, pagination *utils.PaginationPayload, filter *OrderFilter) (result []model.OrderListItem, total int, err error) {
	var (
		args       []interface{}
		conditions []string
	)

	query := `
		SELECT 
			o.id, o.user_id, o.status, COUNT(oi.id) as item_count, 
			o.is_paid, o.is_canceled, o.paid_at, o.canceled_at, 
			o.created_at, o.updated_at
		FROM orders o
		LEFT JOIN order_items oi ON oi.order_id = o.id
	`

	// Build filters dynamically
	if filter != nil {
		if filter.Status != "" {
			conditions = append(conditions, fmt.Sprintf("o.status = $%d", len(args)+1))
			args = append(args, filter.Status)
		}
		if filter.UserId != 0 {
			conditions = append(conditions, fmt.Sprintf("o.user_id = $%d", len(args)+1))
			args = append(args, filter.UserId)
		}
		if filter.StoreId != 0 {
			conditions = append(conditions, fmt.Sprintf("oi.store_id = $%d", len(args)+1))
			args = append(args, filter.StoreId)
		}
		if filter.ProductId != 0 {
			conditions = append(conditions, fmt.Sprintf("oi.product_id = $%d", len(args)+1))
			args = append(args, filter.ProductId)
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += `
		GROUP BY o.id
		ORDER BY o.created_at DESC
		LIMIT $` + fmt.Sprint(len(args)+1) + `
		OFFSET $` + fmt.Sprint(len(args)+2)

	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []model.OrderListItem
	for rows.Next() {
		var o model.OrderListItem
		err := rows.Scan(
			&o.Id,
			&o.UserId,
			&o.Status,
			&o.ItemCount,
			&o.IsPaid,
			&o.IsCanceled,
			&o.PaidAt,
			&o.CanceledAt,
			&o.CreatedAt,
			&o.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	// Count total (without limit/offset)
	countQuery := `
		SELECT COUNT(DISTINCT o.id)
		FROM orders o
		LEFT JOIN order_items oi ON oi.order_id = o.id
	`
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	err = r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *PostgresOrderRepo) UpdateOrderItemStatus(ctx context.Context, orderItemId int, status model.OrderStatus) error {
	query := `
		UPDATE order_items
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, orderItemId)
	return err
}
