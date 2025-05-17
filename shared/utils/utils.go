package utils

import (
	"net/http"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
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

func ToProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

// ParseInt parses a string into an int.
// Returns 0 if parsing fails.
func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
