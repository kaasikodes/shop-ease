package model

import (
	"time"

	"github.com/kaasikodes/shop-ease/shared/types"
)

type OrderStatus string

var (
	UnpaidOrPendingOrderStatus OrderStatus = "pending/unpaid"
	PaidOrderStatus            OrderStatus = "paid"
	ProcessingOrderStatus      OrderStatus = "processing"
	DeliveredOrderStatus       OrderStatus = "delivered"
	FulfilledOrderStatus       OrderStatus = "fulfilled"
	CancelledOrderStatus       OrderStatus = "canceled"
)

type OrderListItem struct {
	Id         int
	UserId     int
	Status     OrderStatus
	ItemCount  int
	IsPaid     bool
	IsCanceled bool
	PaidAt     *time.Time
	CanceledAt *time.Time

	types.Common
}
type Order struct {
	Id         int
	UserId     int
	Status     OrderStatus
	Items      []OrderItem
	IsPaid     bool
	IsCanceled bool
	PaidAt     *time.Time
	CanceledAt *time.Time

	types.Common
}

type OrderItem struct {
	Id             int
	Status         OrderStatus
	ProductId      int
	StoreId        int
	Quantity       int
	Price          float64
	Discount       float64
	AmountToBePaid float64
	types.Common
}
