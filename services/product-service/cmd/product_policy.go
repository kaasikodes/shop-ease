package main

import (
	"net/http"

	"github.com/kaasikodes/shop-ease/shared/types"
	"go.opentelemetry.io/otel/codes"
)

func (app *application) saveProductPolicyHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Save Product Policy")
	defer span.End()

	var input struct {
		SharingFormulaId int                     `json:"sharingFormulaId"`
		PriceToUse       types.DominantPriceType `json:"priceToUse"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.SaveAppProductPolicy(ctx, input.SharingFormulaId, input.PriceToUse); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Product policy saved successfully", nil)
}

func (app *application) getProductPolicyHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "Get Product Policy")
	defer span.End()

	policy, err := app.store.GetAppProductPolicy(ctx)
	if err != nil {
		app.logger.WithContext(ctx).Error("Error retrieving product policy", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Product policy retrieved successfully", policy)
}
