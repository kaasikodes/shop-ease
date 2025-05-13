package main

import (
	"errors"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// retrieve user interactions and filter based on vendor
func (app *application) getUserInteractions(w http.ResponseWriter, r *http.Request) {
	parentCtx, span := app.trace.Start(r.Context(), "Retrieve User Interactions")

	defer span.End()

	// Parse pagination query params (optional enhancement)
	pagination := app.getPaginationFromQuery(r)
	// Parse filter query params (optional enhancement)
	filter := app.getVendorUserInteractionFilterFromQuery(r)
	var interactionType interface{} = filter.Type
	span.SetAttributes(
		attribute.Int("pagination.limit", pagination.Limit),
		attribute.Int("pagination.offset", pagination.Offset),
		attribute.String("filter.type", interactionType.(string)),
		attribute.Int("filter.userId", filter.UserId),
		attribute.Int("filter.vendorId", filter.VendorId),
	)

	// retrieving the user interactions in storage repo
	storeCtx, innerSpan := app.trace.Start(parentCtx, "Retrieve user interactions from storage")
	defer innerSpan.End()
	result, total, err := app.store.plan.GetVendorUserInteractionRecords(pagination, filter)
	if err != nil {
		errMsg := "error encountered while retrieving user interactions"
		app.logger.WithContext(storeCtx).Error(errMsg, err)
		innerSpan.RecordError(err)
		innerSpan.SetStatus(codes.Error, err.Error())
		app.internalServerError(w, r, errors.New(errMsg))
		return
	}
	// Convert []VendorPlan to []any for paginatedResponse
	plans := make([]any, len(result))
	for i, plan := range result {
		plans[i] = plan

	}

	app.jsonResponse(w, http.StatusCreated, ("User interactions retrieved successfully!"), createPaginatedResponse(plans, total))
	return

}
