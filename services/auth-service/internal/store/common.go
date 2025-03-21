package store

import "time"

type Common struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type PaginationPayload struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}