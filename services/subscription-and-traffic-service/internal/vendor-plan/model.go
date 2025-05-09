package vendorplan

import (
	"time"

	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

type VendorPlanActivationPayload struct {
	PlanIds  []int `json:"planIds" validate:"required"`
	IsActive bool  `json:"isActive" validate:"required"`
}
type VendorPlanPayload struct {
	Name                    string        `json:"name" validate:"required,min=5,max=17"`
	Content                 string        `json:"content" validate:"required,min=5,max=350"`
	Price                   float32       `json:"price" validate:"required"`
	UserInteractionsAllowed int           `json:"userInteractionsAllowed" validate:"required"`
	DurationInSecs          time.Duration `json:"duration" validate:"required"` // this in seconds
}
type VendorUserInteractionType string

const (
	UserOrderedItem     VendorUserInteractionType = "user.ordered_item"
	UserIntrestedInItem VendorUserInteractionType = "user.interested_in_item"
)

type VendorUserInteractionFilter struct {
	VendorId int                       `json:"vendorId"`
	UserId   int                       `json:"userId"`
	Type     VendorUserInteractionType `json:"type"`
}
type VendorPlanFilter struct {
	IsActive *bool  `json:"isActive"`
	Name     string `json:"name"`
}
type VendorPlan struct {
	ID                      int           `json:"id"`
	Name                    string        `json:"name"`
	Content                 string        `json:"content"`
	Price                   int           `json:"price"`
	UserInteractionsAllowed int           `json:"userInteractionsAllowed"`
	DurationInSecs          time.Duration `json:"duration"` // this in seconds
	IsActive                bool          `json:"isActive"`

	types.Common
}

type VendorSubsription struct {
	ID              int       `json:"id"`
	PlanId          int       `json:"planId"`
	VendorId        int       `json:"vendorId"`
	HasPaid         bool      `json:"hasPaid"`
	LimitExceededAt time.Time `json:"limitExceededAt"`
	PaidAt          time.Time `json:"paidAt"`
	BeganAt         time.Time `json:"beganAt"`
	ExpiresAt       time.Time `json:"expiresAt"`
	types.Common
}

type VendorUserInteraction struct {
	ID       int                       `json:"id"`
	VendorId int                       `json:"vendorId"`
	UserId   int                       `json:"userId"`
	Type     VendorUserInteractionType `json:"type"`
	types.Common
}
