package vendorplan

type VendorPlanRepo interface {
	CreateVendorPlan() error
	BulkActivateOrDeactivateVendorPlan() error
	CreateVendorPlanSubscription() error
	UpdateVendorPlanSubscription() error
	CreateVendorUserInteractionRecord() error
}
