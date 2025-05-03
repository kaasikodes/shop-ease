package orders

import "github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"

type OrderRepo interface {
	GetOrders(pagination *types.PaginationPayload, filter *types.OrderFilter) (result []types.Order, total int, err error)
}
