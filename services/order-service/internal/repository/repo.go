package repository

import (
	"context"

	"github.com/kaasikodes/shop-ease/services/order-service/internal/model"
	"github.com/kaasikodes/shop-ease/shared/utils"
)

type CreateOrderInputItem struct {
	ProductId      int
	StoreId        int
	Quantity       int
	Price          float64 //TODO: Come up with a better discount logic/model, when you have a bit of spare time
	Discount       float64
	AmountToBePaid float64
}
type OrderFilter struct {
	Status    model.OrderStatus
	UserId    int
	StoreId   int
	ProductId int
}
type OrderRepo interface {
	CreateOrder(ctx context.Context, userId int, items []CreateOrderInputItem) (*int, error)
	UpdateOrderStatus(ctx context.Context, orderId int, status model.OrderStatus) error
	UpdateOrderItemStatus(ctx context.Context, orderItemId int, status model.OrderStatus) error
	GetOrderById(ctx context.Context, orderId int) (model.Order, error)
	GetOrders(ctx context.Context, pagination *utils.PaginationPayload, filter *OrderFilter) (result []model.OrderListItem, total int, err error)
}
