package main

import "time"

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
	return app.config.Env == "production"
}
