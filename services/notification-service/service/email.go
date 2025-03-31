package service

import (
	"context"
	"errors"
	"log"

	"github.com/kaasikodes/shop-ease/services/notification-service/store"
)

type MailConfig struct {
	Addr string
}
type EmailNotificationService struct {
	config MailConfig
}
type EmailNotificationsPayload struct {
	Notifications []EmailNotificationPayload `json:"notifications" validate:"min=1"`
}
type EmailNotificationPayload struct {
	Email   string `json:"email" validate:"required,email,min=5,max=255"`
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content" validate:"required,min=5"`
}

func NewEmailNotificationService(cfg MailConfig) *EmailNotificationService {
	log.Println("Address ...", cfg.Addr)
	return &EmailNotificationService{config: cfg}

}

func (e *EmailNotificationService) SendMultiple(ctx context.Context, notifications []store.Notification) error {
	log.Println("Email Notifications sent ...", len(notifications))
	if notifications == nil {
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
	log.Println("Email Notification sent ...", notification)

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
