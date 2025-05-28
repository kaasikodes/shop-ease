package model

import (
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/model"
)

type PaymentProvider = model.PaymentProvider
type EntityPaymentType = model.EntityPaymentType
type PaymentStatus = model.PaymentStatus

var (
	EntityPaymentTypeVendorSubscriptionPayment EntityPaymentType = model.EntityPaymentTypeVendorSubscriptionPayment
	EntityPaymentTypeOrderPayment              EntityPaymentType = model.EntityPaymentTypeOrderPayment
)
var (
	PaymentProviderPaystack PaymentProvider = model.PaymentProviderPaystack
	PaymentProviderFlutter  PaymentProvider = model.PaymentProviderFlutter
)

type TransactionFilter = model.TransactionFilter
type Transaction = model.Transaction
