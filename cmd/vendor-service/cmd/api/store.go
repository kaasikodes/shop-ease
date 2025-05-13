package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/kaasikodes/shop-ease/services/vendor-service/internal/store"
	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
	"github.com/kaasikodes/shop-ease/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type SaveStorePayload = types.Store

func (app *application) getInventoriesHandler(w http.ResponseWriter, r *http.Request) {

	initialTraceCtx, span := app.trace.Start(r.Context(), "Get Store Inventories")

	defer span.End()
	storeIdStr := chi.URLParam(r, "storeId")
	storeId, err := strconv.Atoi(storeIdStr)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading storeId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("storeId", storeId),
	)

	app.logger.WithContext(initialTraceCtx).Info("getting inventory for vendor/seller")
	productIdStr := r.URL.Query().Get("productId")
	span.SetAttributes(
		attribute.String("filter.productId", productIdStr),
	)
	var productId int
	if productIdStr != "" {
		productId, _ = strconv.Atoi(productIdStr)

	}

	products, total, err := app.store.store.GetInventories(utils.GetPaginationFromQuery(r), &types.InventoryFilter{StoreId: storeId, ProductId: productId})
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error getting inventory for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed getting inventory for vendor/seller")

	var result = make([]any, len(products))
	for i, product := range products {
		result[i] = product

	}

	app.jsonResponse(w, http.StatusOK, "Inventory retrieved successfully!", createPaginatedResponse(result, total))
	return

}
func (app *application) getProductsHandler(w http.ResponseWriter, r *http.Request) {

	initialTraceCtx, span := app.trace.Start(r.Context(), "Get Store products")

	defer span.End()
	storeIdStr := chi.URLParam(r, "storeId")
	storeId, err := strconv.Atoi(storeIdStr)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading storeId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("storeId", storeId),
	)

	app.logger.WithContext(initialTraceCtx).Info("getting store for vendor/seller")
	products, total, err := app.store.store.GetProducts(utils.GetPaginationFromQuery(r), &types.ProductFilter{StoreId: storeId})
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error getting store for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed getting store for vendor/seller")

	var result = make([]any, len(products))
	for i, product := range products {
		result[i] = product

	}

	app.jsonResponse(w, http.StatusOK, "Products retrieved successfully!", createPaginatedResponse(result, total))
	return

}
func (app *application) getStoreHandler(w http.ResponseWriter, r *http.Request) {

	initialTraceCtx, span := app.trace.Start(r.Context(), "Get Store")

	defer span.End()
	storeIdStr := chi.URLParam(r, "storeId")
	storeId, err := strconv.Atoi(storeIdStr)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading storeId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("storeId", storeId),
	)

	app.logger.WithContext(initialTraceCtx).Info("getting store for vendor/seller")
	store, err := app.store.store.GetStoreById(int64(storeId))
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error getting store for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed getting store for vendor/seller")
	span.SetAttributes(attribute.String("name", store.Name),
		attribute.String("description", store.Description),
		attribute.Int("vendorId", store.VendorId),
		attribute.String("account.bank", store.Account.Bank),
		attribute.String("account.swiftCode", store.Account.SwiftCode),
		attribute.String("account.number", store.Account.Number),
		attribute.String("contact.email", store.Contact.Email),
		attribute.String("contact.phone", store.Contact.Phone),
		attribute.String("address.country", store.Address.Country),
		attribute.String("address.state", store.Address.State),
		attribute.String("address.lga", store.Address.Lga),
		attribute.String("address.location", store.Address.Location),
		attribute.String("address.lat", store.Address.Lat),
		attribute.String("address.long", store.Address.Long),
		attribute.String("address.landmark", store.Address.Landmark))
	app.jsonResponse(w, http.StatusOK, "Store retrieved successfully!", &store)
	return

}
func (app *application) updateStoreHandler(w http.ResponseWriter, r *http.Request) {

	initialTraceCtx, span := app.trace.Start(r.Context(), "Update Store")

	defer span.End()
	storeIdStr := chi.URLParam(r, "storeId")
	storeId, err := strconv.Atoi(storeIdStr)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading storeId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	// get the parameters from the request
	var payload SaveStorePayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading save store payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error validating save strore payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("storeId", storeId),
		attribute.String("name", payload.Name),
		attribute.String("description", payload.Description),
		attribute.String("vendorId", fmt.Sprint(payload.VendorId)),
		attribute.String("account.bank", payload.Account.Bank),
		attribute.String("account.swiftCode", payload.Account.SwiftCode),
		attribute.String("account.number", payload.Account.Number),
		attribute.String("contact.email", payload.Contact.Email),
		attribute.String("contact.phone", payload.Contact.Phone),
		attribute.String("address.country", payload.Address.Country),
		attribute.String("address.state", payload.Address.State),
		attribute.String("address.lga", payload.Address.Lga),
		attribute.String("address.location", payload.Address.Location),
		attribute.String("address.lat", payload.Address.Lat),
		attribute.String("address.long", payload.Address.Long),
		attribute.String("address.landmark", payload.Address.Landmark),
	)

	app.logger.WithContext(initialTraceCtx).Info("updating store for vendor/seller")
	store, err := app.store.store.UpdateStore(int64(storeId), payload)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error updating store for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed updating store for vendor/seller")
	app.jsonResponse(w, http.StatusOK, "Store updated successfully!", &store)
	return

}
func (app *application) createStoreHandler(w http.ResponseWriter, r *http.Request) {
	initialTraceCtx, span := app.trace.Start(r.Context(), "Create Store")

	defer span.End()

	// get the parameters from the request
	var payload SaveStorePayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading save store payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error validating save strore payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.String("name", payload.Name),
		attribute.String("description", payload.Description),
		attribute.String("vendorId", fmt.Sprint(payload.VendorId)),
		attribute.String("account.bank", payload.Account.Bank),
		attribute.String("account.swiftCode", payload.Account.SwiftCode),
		attribute.String("account.number", payload.Account.Number),
		attribute.String("contact.email", payload.Contact.Email),
		attribute.String("contact.phone", payload.Contact.Phone),
		attribute.String("address.country", payload.Address.Country),
		attribute.String("address.state", payload.Address.State),
		attribute.String("address.lga", payload.Address.Lga),
		attribute.String("address.location", payload.Address.Location),
		attribute.String("address.lat", payload.Address.Lat),
		attribute.String("address.long", payload.Address.Long),
		attribute.String("address.landmark", payload.Address.Landmark),
	)

	app.logger.WithContext(initialTraceCtx).Info("Creating store for vendor/seller")
	store, err := app.store.store.CreateStore(payload)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error creating store for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed creating store for vendor/seller")
	app.jsonResponse(w, http.StatusOK, "Store created successfully!", &store)
	return

}
func (app *application) bulkAddInventoryHandler(w http.ResponseWriter, r *http.Request) {
	initialTraceCtx, span := app.trace.Start(r.Context(), "Bulk adding inventory")

	defer span.End()
	storeIdStr := chi.URLParam(r, "storeId")
	storeId, err := strconv.Atoi(storeIdStr)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading storeId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	// get the parameters from the request
	var payload []store.Inventory
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading save store payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error validating save store payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("number of inventory records to be added", len(payload)),
		attribute.Int("storeId", storeId),
	)
	for _, item := range payload {
		item.StoreId = storeId //specify the storeId for the inventory

	}

	app.logger.WithContext(initialTraceCtx).Info("Creating store for vendor/seller")
	err = app.store.store.BulkAddInventory(payload)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error creating store for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed adding inventory for vendor/seller")
	app.jsonResponse(w, http.StatusOK, "Inventory added in bulk successfully!", nil)
	return

}
func (app *application) addInventoryHandler(w http.ResponseWriter, r *http.Request) {
	initialTraceCtx, span := app.trace.Start(r.Context(), "Adding inventory")

	defer span.End()
	storeIdStr := chi.URLParam(r, "storeId")
	storeId, err := strconv.Atoi(storeIdStr)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading storeId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	// get the parameters from the request
	var payload store.Inventory
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading save store payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error validating save store payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("productId", payload.ProductId),
		attribute.Int("unitPrice", payload.UnitPrice),
		attribute.Int("quantity", payload.Quantity),
		attribute.Int("storeId", storeId),
	)

	app.logger.WithContext(initialTraceCtx).Info("Creating store for vendor/seller")
	inventories := make([]store.Inventory, 1) //define a payload to match the bulk inventory signature
	payload.StoreId = storeId
	inventories[1] = payload
	err = app.store.store.BulkAddInventory(inventories)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error creating store for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed adding inventory for vendor/seller")
	app.jsonResponse(w, http.StatusOK, "Inventory added successfully!", nil)
	return

}
func (app *application) updateInventoryHandler(w http.ResponseWriter, r *http.Request) {
	initialTraceCtx, span := app.trace.Start(r.Context(), "Updating inventory")

	defer span.End()
	storeIdStr := chi.URLParam(r, "storeId")
	inventoryIdStr := chi.URLParam(r, "inventoryId")
	inventoryId, errId := strconv.Atoi(inventoryIdStr)
	storeId, errStore := strconv.Atoi(storeIdStr)
	err := errors.Join(errId, errStore)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading storeId or inventoryId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	// get the parameters from the request
	var payload store.Inventory
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading save store payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error validating save store payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("productId", payload.ProductId),
		attribute.Int("unitPrice", payload.UnitPrice),
		attribute.Int("quantity", payload.Quantity),
		attribute.Int("storeId", storeId),
	)

	app.logger.WithContext(initialTraceCtx).Info("Creating store for vendor/seller")

	err = app.store.store.UpdateInventory(int64(inventoryId), payload)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error creating store for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed updating inventory for vendor/seller")
	app.jsonResponse(w, http.StatusOK, "Inventory updated successfully!", nil)
	return

}
func (app *application) deleteInventoryHandler(w http.ResponseWriter, r *http.Request) {
	initialTraceCtx, span := app.trace.Start(r.Context(), "Deleting inventory")

	defer span.End()
	storeIdStr := chi.URLParam(r, "storeId")
	inventoryIdStr := chi.URLParam(r, "inventoryId")
	inventoryId, errId := strconv.Atoi(inventoryIdStr)
	storeId, errStore := strconv.Atoi(storeIdStr)
	err := errors.Join(errId, errStore)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading storeId or inventoryId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	// get the parameters from the request
	var payload store.Inventory
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading save store payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error validating save store payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("productId", payload.ProductId),
		attribute.Int("unitPrice", payload.UnitPrice),
		attribute.Int("quantity", payload.Quantity),
		attribute.Int("storeId", storeId),
	)

	app.logger.WithContext(initialTraceCtx).Info("Removing inventory from store for vendor/seller")

	err = app.store.store.UpdateInventory(int64(inventoryId), payload)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error removing inventory from store for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed deleting inventory for vendor/seller")
	app.jsonResponse(w, http.StatusOK, "Inventory deleted successfully!", nil)
	return

}
