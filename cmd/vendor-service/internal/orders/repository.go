package orders

import (
	"database/sql"

	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

type OrderRepo interface {
	GetOrders(pagination *types.PaginationPayload, filter OrderFilter) (result []Order, total int, err error)
	GetOrderById(orderId int) (Order, error)
	UpdateOrder(orderId int, payload Order) (Order, error)
	BulkUpdateStatus(orderIds []int, status OrderStatus) ([]Order, error)
}

type SqlOrderRepo struct {
	db *sql.DB
}

func NewSqlOrderRepo(db *sql.DB) *SqlOrderRepo {
	return &SqlOrderRepo{db}
}

func (o *SqlOrderRepo) GetOrders(pagination *types.PaginationPayload, filter OrderFilter) (result []Order, total int, err error)
func (o *SqlOrderRepo) GetOrderById(orderId int) (Order, error)
func (o *SqlOrderRepo) UpdateOrder(orderId int, payload Order) (Order, error)
func (o *SqlOrderRepo) BulkUpdateStatus(orderIds []int, status OrderStatus) ([]Order, error)
