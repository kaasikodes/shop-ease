package main

import (
	"errors"
	"net/http"

	"github.com/kaasikodes/shop-ease/services/product-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/utils"
)

func (app *application) bulkAddProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Bulk Add Products")
	defer span.End()

	var input []repository.ProductInput
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if len(input) == 0 {
		app.badRequestResponse(w, r, errors.New("product list cannot be empty"))
		return
	}

	if err := app.store.BulkAddProducts(ctx, input); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, "Products added successfully", nil)
}

func (app *application) getProductsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Get Products")
	defer span.End()

	pagination := utils.GetPaginationFromQuery(r)
	products, total, err := app.store.GetProducts(ctx, &utils.PaginationPayload{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	})
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	var result = make([]any, len(products))
	for i, p := range products {
		result[i] = p
	}

	app.jsonResponse(w, http.StatusOK, "Products retrieved successfully", createPaginatedResponse(result, total))
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Delete Product")
	defer span.End()

	productID, err := app.readIntParam(r, "productId")
	if err != nil {
		app.notFoundResponse(w, r, err)
		return
	}

	if err := app.store.DeleteProduct(ctx, productID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Product deleted successfully", nil)
}

func (app *application) updateProductInventoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Update Product Inventory")
	defer span.End()

	inventoryID, err := app.readIntParam(r, "inventoryId")
	if err != nil {
		app.notFoundResponse(w, r, err)
		return
	}

	var payload struct {
		StoreID   int               `json:"storeId"`
		ProductID int               `json:"productId"`
		Quantity  int               `json:"quantity"`
		MetaData  map[string]string `json:"metaData"`
	}

	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.UpdateProductInventory(ctx, inventoryID, payload.StoreID, payload.ProductID, payload.Quantity, &payload.MetaData)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Inventory updated successfully", nil)
}
