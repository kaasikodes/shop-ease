package service

import (
	"context"
	"errors"

	"github.com/kaasikodes/shop-ease/notification/store"
)

type SmsNotificationService struct {
	mailer struct {
		apiKey string
	}
}
type SmsNotificationsPayload struct {
	Notifications []SmsNotificationPayload `json:"notifications" validate:"min=1"`
}
type SmsNotificationPayload struct {
	Phone   string `json:"phone" validate:"required,phone,min=10,max=15"`
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content" validate:"required,min=5"`
}

func (e *SmsNotificationService) SendMultiple(ctx context.Context, notifications []store.Notification) error {
	if notifications == nil || len(notifications) == 0 {
		return errors.New("no notifications were passed in")
	}
	payload := SmsNotificationsPayload{}
	for _, n := range notifications {
		payload.Notifications = append(payload.Notifications, SmsNotificationPayload{
			Phone:   n.Phone,
			Title:   n.Title,
			Content: n.Content,
		})

	}
	if err := Validate.Struct(payload); err != nil {
		return err

	}
	return nil
}
func (e *SmsNotificationService) Send(ctx context.Context, notification *store.Notification) error {
	if notification == nil {
		return errors.New("notification can not be nil")
	}
	payload := SmsNotificationPayload{
		Phone:   notification.Phone,
		Title:   notification.Title,
		Content: notification.Content,
	}
	if err := Validate.Struct(payload); err != nil {
		return err

	}
	return nil

}
