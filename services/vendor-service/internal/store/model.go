package store

import "github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"

type Inventory struct {
	Quantity  int            `json:"quantity" validate:"required"`
	UnitPrice int            `json:"unitPrice" validate:"required"`
	ProductId int            `json:"productId" validate:"required"`
	StoreId   int            `json:"storeId" validate:"required"`
	Product   *types.Product `json:"product"`

	types.Common
}

type Store = types.Store
