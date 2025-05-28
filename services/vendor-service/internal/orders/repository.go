package orders

import (
	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

type OrderRepo interface {
	GetOrders(pagination *types.PaginationPayload, filter OrderFilter) (result []Order, total int, err error)
	GetOrderById(orderId int) (Order, error)
	UpdateOrder(orderId int, payload Order) (Order, error)
	BulkUpdateStatus(orderIds []int, status OrderStatus) ([]Order, error)
	CreateOrder(payload Order) (Order, error)
}
