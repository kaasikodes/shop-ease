package utils

import (
	"net/http"
	"strconv"
)

type PaginationPayload struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func GetPaginationFromQuery(r *http.Request) *PaginationPayload {
	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	if limit <= 0 {
		limit = 10
	}

	return &PaginationPayload{
		Offset: offset,
		Limit:  limit,
	}
}
