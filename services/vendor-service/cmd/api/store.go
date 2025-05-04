package main

import (
	"fmt"
	"net/http"

	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type SaveStorePayload = types.Store

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
	_, err := app.store.store.CreateStore(payload)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error creating store for vendor/seller", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed creating store for vendor/seller")

}
