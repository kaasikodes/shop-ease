package main

import (
	"net/http"
	"strconv"
	"time"

	vendorplan "github.com/kaasikodes/shop-ease/services/subscription-and-traffic-service/internal/vendor-plan"
	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

const (
	ExpiresAtVerificationToken = time.Hour * 24 * 5
)

type paginatedResponse struct {
	Total  int   `json:"total"`
	Result []any `json:"result"`
}

func createPaginatedResponse(result []any, total int) paginatedResponse {
	return paginatedResponse{
		Total:  total,
		Result: result,
	}

}

func (app *application) isProduction() bool {
	return app.config.env == "production"
}

func (app *application) getPaginationFromQuery(r *http.Request) *types.PaginationPayload {
	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if limit <= 0 {
		limit = 10
	}

	return &types.PaginationPayload{
		Offset: offset,
		Limit:  limit,
	}
}

func (app *application) getVendorPlanFilterFromQuery(r *http.Request) *vendorplan.VendorPlanFilter {
	query := r.URL.Query()

	var filter vendorplan.VendorPlanFilter

	if isActiveStr := query.Get("isActive"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			filter.IsActive = &isActive
		}
	}

	if name := query.Get("name"); name != "" {
		filter.Name = name
	}

	return &filter
}

func (app *application) getVendorUserInteractionFilterFromQuery(r *http.Request) *vendorplan.VendorUserInteractionFilter {
	query := r.URL.Query()

	var filter vendorplan.VendorUserInteractionFilter
	userId, _ := strconv.Atoi(query.Get("userId"))
	vendorId, _ := strconv.Atoi(query.Get("vendorId"))
	interactionType := (query.Get("type"))

	filter.UserId = userId
	filter.VendorId = vendorId
	if interactionType != "" {
		switch vendorplan.VendorUserInteractionType(interactionType) {
		case vendorplan.UserOrderedItem, vendorplan.UserIntrestedInItem:
			filter.Type = vendorplan.VendorUserInteractionType(interactionType)

		}
	}

	return &filter
}
