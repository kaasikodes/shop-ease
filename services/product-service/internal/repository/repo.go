package repository

import (
	"context"
	"time"

	"github.com/kaasikodes/shop-ease/services/product-service/internal/model"
	"github.com/kaasikodes/shop-ease/shared/types"
	"github.com/kaasikodes/shop-ease/shared/utils"
)

type CategoryInput struct {
	Name        string
	Description string
}
type ProductInput struct {
	Name             string
	Description      string
	Price            int
	CategoryLabel    string
	SubCategoryLabel []string
	Tags             []string
}

type ProductRepo interface {
	// products: bulkAdd, delete ,getAll, update
	BulkAddProducts(ctx context.Context, payload []ProductInput) error
	DeleteProduct(ctx context.Context, id int) error
	UpdateProduct(ctx context.Context, id int, payload ProductInput) error
	GetProducts(ctx context.Context, id int, pagination *utils.PaginationPayload) (result []model.Product, total int, err error)
	UpdateProductInventory(ctx context.Context, id int, storeId int, productId int, quantity int, metaData *map[string]string) error
	// category: bulkAdd, update, delete, get
	BulkAddCategories(ctx context.Context, payload []CategoryInput) error
	DeleteCategory(ctx context.Context, id int) error
	UpdateCategory(ctx context.Context, id int, payload CategoryInput) error
	GetCategories(ctx context.Context, id int, pagination *utils.PaginationPayload) (result []model.Category, total int, err error)
	// app product policy: save(should be singleton, probably saved as a file, and cached ...), create sharing formula
	CreateSharingFormula(ctx context.Context, id int, basedOn types.SharingFormulaBasedOn, appPercent int, vendorPercent int, description string) error
	SaveAppProductPolicy(ctx context.Context, sharingFormulaId int, priceToUse types.DominantPriceType) error

	// discounts: create, updateApplicability, updateExpiryDate, get
	CreateDiscount(ctx context.Context, payload model.Discount) error
	UpdateDiscountApplicability(ctx context.Context, id int, payload types.DiscountApplicability) error
	UpdateDiscountExpiryDate(ctx context.Context, id int, expiryDate time.Time) error
	GetDiscounts(ctx context.Context, pagination *utils.PaginationPayload) (result []model.Discount, total int, err error)
}
