package service

import (
	"context"
	"errors"

	"github.com/kaasikodes/shop-ease/notification/store"
)

type InAppNotificationService struct {
	store store.NotificationStore
}
type InAppNotificationsPayload struct {
	Notifications []InAppNotificationPayload `json:"notifications" validate:"min=1"`
}
type InAppNotificationPayload struct {
	Email   string `json:"email" validate:"required,email,min=5,max=255"`
	Phone   string `json:"phone" validate:"-"`
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content" validate:"required,min=5"`
}

func (e *InAppNotificationService) Send(ctx context.Context, notification *store.Notification) error {
	if notification == nil {
		return errors.New("notification can not be nil")
	}
	payload := InAppNotificationPayload{
		Email:   notification.Email,
		Phone:   notification.Phone,
		Title:   notification.Title,
		Content: notification.Content,
	}
	if err := Validate.Struct(payload); err != nil {
		return err

	}
	if _, err := e.store.Create(ctx, notification); err != nil {
		return err
	}
	return nil

}
func (e *InAppNotificationService) SendMultiple(ctx context.Context, notifications []store.Notification) error {
	if notifications == nil || len(notifications) == 0 {
		return errors.New("no notifications were passed in")
	}
	payload := InAppNotificationsPayload{}
	for _, n := range notifications {
		payload.Notifications = append(payload.Notifications, InAppNotificationPayload{
			Email:   n.Email,
			Phone:   n.Phone,
			Title:   n.Title,
			Content: n.Content,
		})

	}

	if err := Validate.Struct(payload); err != nil {
		return err

	}
	if _, err := e.store.CreateMultiple(ctx, notifications); err != nil {
		return err
	}
	return nil

}
