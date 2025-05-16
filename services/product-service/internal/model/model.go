package model

import (
	"github.com/kaasikodes/shop-ease/shared/types"
)

type Category struct {
	types.CommonDescriptiveModel
}

type Product struct {
	Inventory     []Inventory `json:"inventory"`
	Category      Category    `json:"category"`
	SubCategories []Category  `json:"subCategories"`
	Price         types.Price `json:"price"`
	Tags          []string    `json:"tags"`

	types.CommonDescriptiveModel
}
type AppProductPolicy struct {
	Id                      int                     `json:"id"`
	CurrentSharingFormulaId int                     `json:"sharingFormulaId"`
	CurrentSharingFormula   types.SharingFormula    `json:"sharingFormula"`
	ProductPriceToUse       types.DominantPriceType `json:"productPriceToUse"` //defaults to store
	// TODO: What happens when the app wants to have a universal discount, who is responsible for compensating the buyer
	// How is money remiited to the vendor and the app - sharing formula and what is the standard
	// e.g lets say payment is made for a product via the order, and the app gets 10% while the vendor gets 90%. Its all recorded as paid to company's account, and on the sales record of the vendor its calculated and shown dynamically, but this formula can change so the sharing formula will be historical with each order made having a sharing formula attached to it.

}

type Inventory struct { //just keeps a record of the total current inventory of stores that have this product
	Id        int `json:"id"`
	Quantity  int `json:"quantity" validate:"required"`
	ProductId int `json:"productId" validate:"required"`
	StoreId   int `json:"storeId" validate:"required"`
	MetaData  *map[string]string
}

type Discount = types.Discount
