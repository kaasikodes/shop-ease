package model

import (
	"time"

	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

type PaymentProvider string
type EntityPaymentType string
type PaymentStatus string

var (
	VendorSubscriptionPayment EntityPaymentType = "vendor"
	OrderPayment              EntityPaymentType = "order"
)
var (
	PaystackPaymentProvider PaymentProvider = "paystack"
	FlutterPaymentProvider  PaymentProvider = "flutter"
)

type TransactionFilter struct {
	Provider          PaymentProvider   `json:"provider"`
	Amount            float64           `json:"amount"`
	EntityPaymentType EntityPaymentType `json:"entityPaymentType"`
	Status            PaymentStatus     `json:"status"`
	PaidAt            *time.Time        `json:"paidAt"`
}
type Transaction struct {
	ID                int               `json:"id"`
	Provider          PaymentProvider   `json:"provider"`
	TransactionId     string            `json:"transactionId"`
	MetaData          map[string]string `json:"metaData"`
	EntityId          int               `json:"entityId"`
	Amount            float64           `json:"amount"`
	EntityPaymentType EntityPaymentType `json:"entityPaymentType"`
	Status            PaymentStatus     `json:"status"`
	PaidAt            *time.Time        `json:"paidAt"`
	types.Common
}
