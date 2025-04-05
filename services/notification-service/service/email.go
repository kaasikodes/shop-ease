package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/kaasikodes/shop-ease/services/notification-service/config"
	"github.com/kaasikodes/shop-ease/services/notification-service/store"
	"github.com/kaasikodes/shop-ease/shared/logger"
	gomail "gopkg.in/mail.v2"
)

type MailConfig = config.MailConfig
type EmailNotificationService struct {
	config MailConfig
	logger logger.Logger
}
type EmailNotificationsPayload struct {
	Notifications []EmailNotificationPayload `json:"notifications" validate:"min=1"`
}
type EmailNotificationPayload struct {
	Email   string `json:"email" validate:"required,email,min=5,max=255"`
	Title   string `json:"title" validate:"required,min=3"`
	Content string `json:"content" validate:"required,min=5"`
}

func NewEmailNotificationService(cfg MailConfig, logger logger.Logger) *EmailNotificationService {
	log.Println("Email Service Address ...", cfg.FromEmail)
	return &EmailNotificationService{config: cfg, logger: logger}

}
func (e *EmailNotificationService) sendMail(toEmails []string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", e.config.FromEmail)
	message.SetHeader("To", toEmails...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", body)

	dialer := gomail.NewDialer(e.config.Host, e.config.Port, e.config.Username, e.config.Password)
	e.logger.Info("WHAT ARE THOSE", e.config.Host, e.config.Port, e.config.Username, e.config.Password)

	err := dialer.DialAndSend(message)
	if err != nil {
		e.logger.Error("notification: error sending mail %v", err)
		return err
	}
	return nil

}

func (e *EmailNotificationService) SendMultiple(ctx context.Context, notifications []store.Notification) error {
	if notifications == nil {
		return errors.New("no notifications were passed in")
	}
	payload := EmailNotificationsPayload{}
	payload.Notifications = make([]EmailNotificationPayload, len(notifications))

	for i, n := range notifications {
		payload.Notifications[i] = EmailNotificationPayload{
			Email:   n.Email,
			Title:   n.Title,
			Content: n.Content,
		}

	}
	if err := Validate.Struct(payload); err != nil {
		return err

	}
	var errCh = make(chan error, len(payload.Notifications))
	var wg sync.WaitGroup
	wg.Add(len(payload.Notifications))
	for _, n := range payload.Notifications {
		go func(notification EmailNotificationPayload) {
			defer wg.Done()
			toEmails := []string{notification.Email}

			if err := e.sendMail(toEmails, notification.Title, notification.Content); err != nil {
				errCh <- err
			}
			e.logger.Info("EMAIL SENT to: %s", notification.Email)

		}(n)

	}
	wg.Wait()
	close(errCh)
	var errs []error
	for err := range errCh {
		errs = append(errs, err)

	}
	if len(errs) > 0 {
		e.logger.Error(fmt.Sprintf("Notiticatio::Errors errors while sending mail to %v", notifications), errors.Join(errs...))
		return errors.Join(errs...)

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
	toEmails := []string{notification.Email}

	if err := e.sendMail(toEmails, notification.Title, notification.Content); err != nil {
		return err
	}
	e.logger.Info("EMAIL SENT to: %s", notification.Email)

	return nil

}
