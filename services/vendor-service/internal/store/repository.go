package store

import (
	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
	"github.com/kaasikodes/shop-ease/shared/utils"
)

type StoreRepo interface {
	// create store
	CreateStore(payload Store) (*Store, error)
	// update store
	UpdateStore(id int64, payload Store) (*Store, error)
	// Get store by Id
	GetStoreById(id int64) (*Store, error)
	// Get products (allow for filter by storeId)
	GetProducts(pagination *utils.PaginationPayload, filter *types.ProductFilter) (result []types.Product, total int, err error)
	// Add inventory in bulk
	BulkAddInventory(payload []Inventory) error
	// Update inventory
	UpdateInventory(id int64, payload Inventory) error
	// delete inventory
	DeleteInventory(id int64) (*int64, error)
	// Get Inventories
	GetInventories(pagination *utils.PaginationPayload, filter *types.InventoryFilter) (result []Inventory, total int, err error)
}

// TODO: Create a SqlStoreRepo that implements the interface above

// vendor(grpc) - > subscription (grpc)| -> payment (grpc)
