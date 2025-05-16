package store

import (
	"time"

	vendor_types "github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
	"github.com/kaasikodes/shop-ease/shared/types"
)

type Inventory struct {
	Quantity             int                  `json:"quantity" validate:"required"`
	UnitCostPrice        *int                 `json:"unitCostPrice" validate:"-"` //not required but if put in will help with sales report
	ProductId            int                  `json:"productId" validate:"required"`
	StoreId              int                  `json:"storeId" validate:"required"`
	Product              vendor_types.Product `json:"product"`
	ArrivalorProduceDate *time.Time           `json:"arrivalOrProduceDate"` //considerations for a agro-ecommerce site
	Price                types.Price
	vendor_types.Common
}
type StoreProductPolicy struct {
	ProductPriceToUse types.DominantPriceType `json:"productPriceToUse"`
}
type Store = vendor_types.Store
