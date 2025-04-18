package sql_store

import (
	"database/sql"

	"github.com/kaasikodes/shop-ease/services/notification-service/store"
)

type PaginationPayload = store.PaginationPayload
type NotificationFilter = store.NotificationFilter
type Notification = store.Notification

type SqlStorage struct {
	db *sql.DB
}

func NewSQLStorage(db *sql.DB) store.Storage {
	return &SqlStorage{
		db,
	}

}

func (sq *SqlStorage) Notification() store.NotificationStore {
	return &SQLNotificationStore{
		db: sq.db,
	}

}
