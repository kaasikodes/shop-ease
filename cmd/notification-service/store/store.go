package store

import (
	"context"
	"time"
)

type Notification struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`  // can be null
	IsRead    bool       `json:"isRead"` //defaults to false
	ReadAt    *time.Time `json:"readAt"` //can be null
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
type NotificationFilter struct {
	Email  *string `json:"email"`
	Phone  *string `json:"phone"`
	IsRead *bool   `json:"isRead"`
}
type PaginationPayload struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

var (
	DefaultPagination PaginationPayload = PaginationPayload{20, 0}
)

type NotificationStore interface {
	Get(ctx context.Context, pagination *PaginationPayload, filter *NotificationFilter) ([]Notification, int, error)
	GetOne(ctx context.Context, notificationId int) (*Notification, error)
	Create(ctx context.Context, notification *Notification) (*Notification, error)
	CreateMultiple(ctx context.Context, notification []Notification) ([]Notification, error)
}

type Storage interface {
	Notification() NotificationStore
}
