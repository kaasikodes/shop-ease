package main

import (
	"errors"
	"fmt"
	"net/http"

	vendorplan "github.com/kaasikodes/shop-ease/services/subscription-and-traffic-service/internal/vendor-plan"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// create vendor plans
func (app *application) createVendorPlan(w http.ResponseWriter, r *http.Request) {
	parentCtx, span := app.trace.Start(r.Context(), "Create Vendor Subscription Plan")

	defer span.End()

	// get the parameters from
	var payload vendorplan.VendorPlanPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(parentCtx).Error("Error reading vendor subscription payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(parentCtx).Error("Error validating vendor subscription payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.String("name", payload.Name),
		attribute.String("content", payload.Content),
		attribute.Int("userInteractions", payload.UserInteractionsAllowed),
		attribute.Int("durationInSecs", int(payload.DurationInSecs)),
		attribute.Float64("price", float64(payload.Price)),
	)

	// saving the vendpor plan in storage repo
	storeVendorPlanCtx, span := app.trace.Start(parentCtx, "Saving Vendor Subscription Plan")
	defer span.End()
	plan, err := app.store.plan.CreateVendorPlan(payload)
	if err != nil {
		errMsg := "error encountered while saving the vendor plan"
		app.logger.WithContext(storeVendorPlanCtx).Error(errMsg, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.internalServerError(w, r, errors.New(errMsg))
		return
	}

	app.jsonResponse(w, http.StatusCreated, "Vendor plan subscription created successfully!", &plan)
	return

}

// activate or deactivate vendor plans in bulk
func (app *application) toggleVendorPlanActivation(w http.ResponseWriter, r *http.Request) {
	parentCtx, span := app.trace.Start(r.Context(), "Toggle Vendor Subscription Plan")

	defer span.End()

	// get the parameters from
	var payload vendorplan.VendorPlanActivationPayload
	if err := readJson(w, r, &payload); err != nil {
		app.logger.WithContext(parentCtx).Error("Error reading vendor subscription payload as json", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.logger.WithContext(parentCtx).Error("Error validating vendor subscription payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Bool("isActive", payload.IsActive),
		attribute.IntSlice("planIds", payload.PlanIds),
	)

	// updating the vendpor plan in storage repo
	storeVendorPlanCtx, span := app.trace.Start(parentCtx, "Updating Vendor Subscription Plan")
	defer span.End()
	err := app.store.plan.BulkActivateOrDeactivateVendorPlan(payload.PlanIds, payload.IsActive)
	if err != nil {
		errMsg := "error encountered while updating the vendor plan"
		app.logger.WithContext(storeVendorPlanCtx).Error(errMsg, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.internalServerError(w, r, errors.New(errMsg))
		return
	}

	app.jsonResponse(w, http.StatusCreated, fmt.Sprintf("%v Vendor subscription plans toggled successfully to %v!", len(payload.PlanIds), payload.IsActive), nil)
	return

}

// retrieve vendor subscription plans
func (app *application) getVendorPlans(w http.ResponseWriter, r *http.Request) {
	parentCtx, span := app.trace.Start(r.Context(), "Retrieve Vendor Subscription Plans")

	defer span.End()

	// Parse pagination query params (optional enhancement)
	pagination := app.getPaginationFromQuery(r)
	// Parse filter query params (optional enhancement)
	filter := app.getVendorPlanFilterFromQuery(r)
	span.SetAttributes(
		attribute.Int("pagination.limit", pagination.Limit),
		attribute.Int("pagination.offset", pagination.Offset),
		attribute.Bool("filter.isActive", *filter.IsActive),
	)

	// retrieving the vendpor plans in storage repo
	storeCtx, innerSpan := app.trace.Start(parentCtx, "Retrieve Vendor Subscription Plans from storage")
	defer innerSpan.End()
	result, total, err := app.store.plan.GetVendorPlans(pagination, filter)
	if err != nil {
		errMsg := "error encountered while retrieving the vendor plans"
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

	app.jsonResponse(w, http.StatusCreated, ("Vendor subscription plans retrieved successfully!"), createPaginatedResponse(plans, total))
	return

}

// handler to listen to the following events - user interactions: order made for product(order service); item added to wishlist(search & recommend - can change ), payment made for subscription
// emits the following events
