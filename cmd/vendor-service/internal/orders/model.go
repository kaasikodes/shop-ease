package orders

import (
	"github.com/kaasikodes/shop-ease/services/vendor-service/internal/products"
	"github.com/kaasikodes/shop-ease/services/vendor-service/internal/store"
	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

type OrderStatus string
type OrderFilter struct {
	ProductId int         `json:"productId"`
	StoreId   int         `json:"storeId"`
	Status    OrderStatus `json:"status"`
}

var (
	PendingOrderStatus                OrderStatus = "pending"
	AcceptedOrderStatus               OrderStatus = "accepted"
	RejectedOrderStatus               OrderStatus = "rejected"
	ShippedOrderStatus                OrderStatus = "shipped"
	DeliveredOrderStatus              OrderStatus = "delivered" //courier signs confirms that he has delivered
	FulfilledOrderStatus              OrderStatus = "fulfilled" //customer confirms the delivery
	ReturnedBackByCustomerOrderStatus OrderStatus = "returned-back-by-customer"
)

type Order struct {
	ID                    int               `json:"id"`
	Quantity              int               `json:"quantity"`
	UnitPrice             int               `json:"unitPrice"`
	Status                OrderStatus       `json:"status"`
	ProductId             int               `json:"productId"`
	StoreId               int               `json:"storeId"`
	FulfillingInventoryId int               `json:"fulfillingInventoryId"`
	FulfillingInventory   *store.Inventory  `json:"fulfillingInventory"`
	Product               *products.Product `json:"product"`
	types.Common
}
