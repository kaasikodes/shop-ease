package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (app *application) getSellerHandler(w http.ResponseWriter, r *http.Request) {

	initialTraceCtx, span := app.trace.Start(r.Context(), "Get Seller/Vendor")

	defer span.End()
	sellerIdStr := chi.URLParam(r, "sellerId")
	sellerId, err := strconv.Atoi(sellerIdStr)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading sellerId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("sellerId", sellerId),
	)

	app.logger.WithContext(initialTraceCtx).Info("getting vendor/seller from store ")
	seller, err := app.store.seller.GetVendor(int64(sellerId))
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error getting  vendor/seller from store ", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed getting store for vendor/seller")
	span.SetAttributes(
		attribute.String("name", seller.Name),
		attribute.String("email", seller.Email),
		attribute.Int("userId", seller.UserId),
		attribute.String("phone", seller.Phone),
	)
	app.jsonResponse(w, http.StatusOK, "Seller/vendor retrieved successfully!", &seller)
	return

}
