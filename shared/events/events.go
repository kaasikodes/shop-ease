package events

var (
	ProductCreatedEvent       = "product.created"
	ProductUpdatedEvent       = "product.updated"
	ProductLowStockEvent      = "product.low_stock"
	UserCreatedEvent          = "user.created"
	UserUpdatedEvent          = "user.updated"
	UserOrderedItemEvent      = "user.ordered_item"
	UserInterestedInItemEvent = "user.interested_in_item"
)

const (
	ProductTopic      = "product"
	AuthTopic         = "auth"
	SubscriptionTopic = "subscription"
	VendorTopic       = "vendor"
	PaymentTopic      = "payment"
	OrderTopic        = "order"
)
