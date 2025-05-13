package service

import (
	"context"
	"errors"

	"github.com/kaasikodes/shop-ease/services/notification-service/store"
	"github.com/kaasikodes/shop-ease/shared/logger"
)

type InAppNotificationService struct {
	store  store.NotificationStore
	logger logger.Logger
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

func NewInAppNotificationService(store store.NotificationStore, logger logger.Logger) *InAppNotificationService {
	return &InAppNotificationService{store, logger}

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
	e.logger.Info("IN_APP:", "notification sent ....")
	return nil

}
func (e *InAppNotificationService) SendMultiple(ctx context.Context, notifications []store.Notification) error {
	if notifications == nil {
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
