package vendorplan

import (
	"database/sql"

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
}

type SqlVendorRepo struct {
	db *sql.DB
}

func (r *SqlVendorRepo) GetVendorUserInteractionRecords(pagination *types.PaginationPayload, filter *VendorUserInteractionFilter) (result []VendorUserInteraction, total int, err error)
func (r *SqlVendorRepo) GetVendorPlans(pagination *types.PaginationPayload, filter *VendorPlanFilter) (result []VendorPlan, total int, err error)
func (r *SqlVendorRepo) CreateVendorPlan(payload VendorPlanPayload) (*VendorPlan, error)
func (r *SqlVendorRepo) BulkActivateOrDeactivateVendorPlan(planIds []int, isAcvtive bool) error
func (r *SqlVendorRepo) CreateVendorPlanSubscription(planId int, vendorId int) (*VendorSubsription, error)
func (r *SqlVendorRepo) UpdateVendorPlanSubscription(subscriptionId int, payload VendorSubsription) (*VendorSubsription, error)
func (r *SqlVendorRepo) CreateVendorUserInteractionRecord(userId int, interactionType VendorUserInteractionType) (*VendorUserInteraction, error)

func NewSqlVendorRepo(db *sql.DB) *SqlVendorRepo {
	return &SqlVendorRepo{db}

}
