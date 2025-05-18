package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
	"github.com/kaasikodes/shop-ease/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// type SaveStorePayload = types.Store

func (app *application) getOrdersHandler(w http.ResponseWriter, r *http.Request) {

	initialTraceCtx, span := app.trace.Start(r.Context(), "Get Store Orders")

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
