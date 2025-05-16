package main

import (
	"net/http"
	"time"

	"github.com/kaasikodes/shop-ease/services/product-service/internal/model"
	"github.com/kaasikodes/shop-ease/shared/types"
	"github.com/kaasikodes/shop-ease/shared/utils"
	"go.opentelemetry.io/otel/codes"
)

func (app *application) createDiscountHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Create Discount")
	defer span.End()

	var input model.Discount
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.CreateDiscount(ctx, input); err != nil {
		app.logger.WithContext(ctx).Error("Error creating discount", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusCreated, "Discount created successfully", nil)
}

func (app *application) getDiscountsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Get Discounts")
	defer span.End()

	pagination := utils.GetPaginationFromQuery(r)
	discounts, total, err := app.store.GetDiscounts(ctx, &utils.PaginationPayload{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	}, nil)
	if err != nil {
		app.logger.WithContext(ctx).Error("Error getting discounts", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.internalServerError(w, r, err)
		return
	}

	var result = make([]any, len(discounts))
	for i, discount := range discounts {
		result[i] = discount
	}

	app.jsonResponse(w, http.StatusOK, "Discounts retrieved successfully", createPaginatedResponse(result, total))
}

func (app *application) updateDiscountApplicabilityHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Update Discount Applicability")
	defer span.End()

	discountID, err := app.readIntParam(r, "discountId")
	if err != nil {
		app.notFoundResponse(w, r, err)
		return
	}

	var input types.DiscountApplicability
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.UpdateDiscountApplicability(ctx, discountID, input); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Discount applicability updated successfully", nil)
}

func (app *application) updateDiscountExpiryDateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Update Discount Expiry Date")
	defer span.End()

	discountID, err := app.readIntParam(r, "discountId")
	if err != nil {
		app.notFoundResponse(w, r, err)
		return
	}

	var payload struct {
		ExpiryDate time.Time `json:"expiryDate"`
	}
	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.UpdateDiscountExpiryDate(ctx, discountID, payload.ExpiryDate); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Discount expiry date updated", nil)
}
