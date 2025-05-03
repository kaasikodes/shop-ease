package store

import "github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"

type Inventory struct {
	Quantity  int            `json:"quantity"`
	UnitPrice int            `json:"unitPrice"`
	ProductId int            `json:"productId"`
	StoreId   int            `json:"storeId"`
	Product   *types.Product `json:"product"`

	types.Common
}

type Store = types.Store
