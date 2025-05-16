package events

// TODO: Refactor to be  a map say -> map[EventTopic][Events] => map[AuthTopic] ... think more
var (
	ProductCreatedEvent       = "product.created"
	ProductUpdatedEvent       = "product.updated"
	ProductLowStockEvent      = "product.low_stock"
	UserCreatedEvent          = "user.created"
	UserUpdatedEvent          = "user.updated"
	UserOrderedItemEvent      = "user.ordered_item"
	UserInterestedInItemEvent = "user.interested_in_item"
	// payment listens
	VendorSubscriptionCreated = "subscription.vendor_subcription_created"
	OrderCreated              = "order.order_placed"
	// payment sends
	VendorSubscriptionPaymnentMade = "payment.vendor_subcription_paid_for"
	OrderPaymnentMade              = "payment.order_paid_for"
)

const (
	ProductTopic      = "product"
	AuthTopic         = "auth"
	SubscriptionTopic = "subscription"
	VendorTopic       = "vendor"
	PaymentTopic      = "payment"
	OrderTopic        = "order"
)
