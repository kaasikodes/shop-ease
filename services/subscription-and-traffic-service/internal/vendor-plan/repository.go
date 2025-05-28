package vendorplan

import (
	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

type VendorPlanRepo interface {
	CreateVendorPlan(payload VendorPlanPayload) (*VendorPlan, error)
	GetVendorPlans(pagination *types.PaginationPayload, filter *VendorPlanFilter) (result []VendorPlan, total int, err error)
	BulkActivateOrDeactivateVendorPlan(planIds []int, isActive bool) error
	CreateVendorPlanSubscription(planId int, vendorId int) (*VendorSubsription, error)
	UpdateVendorPlanSubscription(subscriptionId int, payload VendorSubsription) (*VendorSubsription, error)
	CreateVendorUserInteractionRecord(userId int, interactionType VendorUserInteractionType) (*VendorUserInteraction, error)
	GetVendorUserInteractionRecords(pagination *types.PaginationPayload, filter *VendorUserInteractionFilter) (result []VendorUserInteraction, total int, err error)
	GetActiveSubscriptionsForVendor(vendorId int64) ([]*VendorSubsription, error)
	GetVendorSubscriptionID(subscriptionId int64) (*VendorSubsription, error)
	MarkSubscriptionPaid(subscriptionId int64) error
}
