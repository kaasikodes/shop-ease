package service

import (
	"context"
	"errors"

	"github.com/kaasikodes/shop-ease/notification/store"
)

type EmailNotificationService struct {
	mailer struct {
		apiKey string
	}
}
type EmailNotificationsPayload struct {
	Notifications []EmailNotificationPayload `json:"notifications" validate:"min=1"`
}
type EmailNotificationPayload struct {
	Email   string `json:"email" validate:"required,email,min=5,max=255"`
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content" validate:"required,min=5"`
}

func (e *EmailNotificationService) SendMultiple(ctx context.Context, notifications []store.Notification) error {
	if notifications == nil || len(notifications) == 0 {
		return errors.New("no notifications were passed in")
	}
	payload := EmailNotificationsPayload{}
	for _, n := range notifications {
		payload.Notifications = append(payload.Notifications, EmailNotificationPayload{
			Email:   n.Email,
			Title:   n.Title,
			Content: n.Content,
		})

	}
	if err := Validate.Struct(payload); err != nil {
		return err

	}
	return nil
}
func (e *EmailNotificationService) Send(ctx context.Context, notification *store.Notification) error {
	if notification == nil {
		return errors.New("notification can not be nil")
	}
	payload := EmailNotificationPayload{
		Email:   notification.Email,
		Title:   notification.Title,
		Content: notification.Content,
	}
	if err := Validate.Struct(payload); err != nil {
		return err

	}
	return nil

}
