package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
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

func (app *application) readIntParam(r *http.Request, key string) (int, error) {
	param := chi.URLParam(r, key)
	id, err := strconv.Atoi(param)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return id, nil
}
