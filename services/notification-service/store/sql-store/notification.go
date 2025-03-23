package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kaasikodes/shop-ease/notification/store"
)

type SQLNotificationStore struct {
	db *sql.DB
}

// Get fetches notifications with optional filters and pagination
func (s *SQLNotificationStore) Get(ctx context.Context, pagination *PaginationPayload, filter *NotificationFilter) ([]Notification, int, error) {
	query := "SELECT id, email, phone, isRead, readAt, title, content, createdAt, updatedAt FROM notifications WHERE 1=1"
	countQuery := "SELECT COUNT(*) FROM notifications WHERE 1=1"

	var args []interface{}
	if pagination == nil {
		pagination = &store.DefaultPagination

	}

	if filter.Email != nil {
		query += " AND user_id = ?"
		countQuery += " AND user_id = ?"
		args = append(args, *filter.Email)
	}
	if filter.Phone != nil {
		query += " AND phone = ?"
		countQuery += " AND phone = ?"
		args = append(args, *filter.Phone)
	}
	if filter.IsRead != nil {
		query += " AND isRead = ?"
		countQuery += " AND isRead = ?"
		args = append(args, *filter.IsRead)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pagination.Limit, pagination.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.Email, &n.Phone, &n.IsRead, &n.ReadAt, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, n)
	}

	var count int
	err = s.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return notifications, count, nil
}

// GetOne fetches a single notification by ID
func (s *SQLNotificationStore) GetOne(ctx context.Context, notificationId int) (*Notification, error) {
	query := "SELECT id, email, phone, isRead, readAt, title, content, createdAt, updatedAt FROM notifications WHERE id = ?"
	var n Notification
	err := s.db.QueryRowContext(ctx, query, notificationId).Scan(&n.ID, &n.Email, &n.Phone, &n.IsRead, &n.ReadAt, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, err
	}
	return &n, nil
}

// Create inserts a new notification into the database
func (s *SQLNotificationStore) Create(ctx context.Context, notification *Notification) (*Notification, error) {
	query := "INSERT INTO notifications (email, phone, title, content) VALUES (?, ?, ?, ?)"
	result, err := s.db.ExecContext(ctx, query, notification.Email, notification.Phone, notification.Title, notification.Content)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	notification.ID = int(lastID)
	return notification, nil
}

// CreateMultiple inserts multiple notifications at once
func (s *SQLNotificationStore) CreateMultiple(ctx context.Context, notifications []Notification) ([]Notification, error) {
	query := "INSERT INTO notifications (email, phone, title, content) VALUES "
	var args []interface{}
	var values []string

	for _, n := range notifications {
		values = append(values, "(?, ?, ?, ?)")
		args = append(args, n.Email, n.Phone, n.Title, n.Content)
	}

	query += fmt.Sprintf("%s", values)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
