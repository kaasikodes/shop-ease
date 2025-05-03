package events

var (
	ProductCreatedEvent  = "product.created"
	ProductUpdatedEvent  = "product.updated"
	ProductLowStockEvent = "product.low_stock"
	UserCreatedEvent     = "user.created"
	UserUpdatedEvent     = "user.updated"
)

const (
	ProductTopic      = "product"
	AuthTopic         = "auth"
	SubscriptionTopic = "subscription"
	VendorTopic       = "vendor"
	PaymentTopic      = "payment"
	OrderTopic        = "order"
)
