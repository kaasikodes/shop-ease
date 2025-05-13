package types

import (
	"time"

	"github.com/kaasikodes/shop-ease/shared/utils"
)

type Common struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type PaginationPayload = utils.PaginationPayload

type ProductFilter struct {
	StoreId int `json:"storeId"`
}
type InventoryFilter struct {
	ProductId int `json:"productId"`
	StoreId   int `json:"storeId"`
}
type OrderFilter struct {
	ProductId int `json:"productId"`
	StoreId   int `json:"storeId"`
}

type Order struct {
	ID        int `json:"id"`
	ProductId int `json:"productId"`
	StoreId   int `json:"storeId"`
	Quantity  int `json:"quantity"`
	Common
}
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Common
}
type Vendor struct {
	ID     int    `json:"id"`
	UserId int    `json:"userId"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Common
}

type Store struct {
	ID          int    `json:"id"`
	VendorId    int    `json:"vendorId"`
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"required,max=100"`
	Address     Address
	Contact     Contact
	Account     Account
	Common
}
type Account struct {
	Bank      string `json:"bank" validate:"required"`
	Number    string `json:"number" validate:"required"`
	SwiftCode string `json:"swiftCode" validate:"required"`
}
type Contact struct {
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email" validate:"required,email,max=255"`
}
type Address struct {
	Location   string `json:"location" validate:"required"`
	Lat        string `json:"lat" validate:"required"`
	Long       string `json:"long" validate:"required"`
	Country    string `json:"country" validate:"required"`
	State      string `json:"state" validate:"required"`
	Lga        string `json:"lga" validate:"-"`
	Landmark   string `json:"landmark" validate:"-"`
	Timezone   string `json:"timezone" validate:"-"`
	PostalCode string `json:"postalCode" validate:"-"`
}
