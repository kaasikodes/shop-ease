package types

import "time"

type PaginationPayload struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Common struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CommonDescriptiveModel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Common
}
type DominantPriceType string

var (
	StoreProductPrice     DominantPriceType = "store"
	ProductPrice          DominantPriceType = "product"
	InventoryProductPrice DominantPriceType = "inventory"
)

type Price struct {
	Amount   int       `json:"amount"`
	Discount *Discount `json:"discount"`
}
type DiscountValueType string
type PaidBy string

var (
	PaidByApp    PaidBy = "app"
	PaidByVendor PaidBy = "vendor"
)

var (
	PercentageDiscount DiscountValueType = "percentage"
	AmountDiscount     DiscountValueType = "amount"
)

type DiscountApplicability struct {
	ProductIds               []int64
	StoreProductIds          []int64
	StoreProductInventoryIds []int64
}
type Discount struct {
	Id           int                   `json:"id"`
	Value        int16                 `json:"percentage" validate:"required"`
	ValueType    DiscountValueType     `json:"type" validate:"required"`
	EffectiveAt  time.Time             `json:"effectiveAt" validate:"required"`
	ExpiresAt    *time.Time            `json:"expiresAt" validate:"-"`
	PaidBy       PaidBy                `json:"paidBy" validate:"required"` //defaults to app
	ApplicableTo DiscountApplicability `json:"applicableTo" validate:"required"`
	CommonDescriptiveModel
}

type SharingFormulaBasedOn string

var (
	SharingFormulaOnVendorBasis SharingFormulaBasedOn = "sale"
	SharingFormulaOnProfitBasis SharingFormulaBasedOn = "profit"
)

type SharingFormula struct {
	Id          int                   `json:"id"`
	App         int                   `json:"app"`
	Vendor      int                   `json:"vendor"`
	BasedOn     SharingFormulaBasedOn `json:"basedOn"`     // defaults to sale - selling price of item, more reliable
	Description string                `json:"description"` // optional notes
	// must sum up to 100
	// concern is it on the sale or the profit because certain products might not have a selling price defined- so how does prfit sharing work or the app make back its money
}
