package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/kaasikodes/shop-ease/services/notification-service/store"
)

// go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

// protoc --go_out=. --go-grpc_out=. proto/notification.proto

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

}

// notification can go on through 3 identified channels which are email, sms, in-app(database), others might be addded

type NotificationService interface {
	Send(ctx context.Context, notification *store.Notification) error
	SendMultiple(ctx context.Context, notifications []store.Notification) error
}

type NotificationType string

var (
	EmailNotificationType NotificationType = "email"
	InAppNotificationType NotificationType = "in-app"
	SmsNotificationType   NotificationType = "sms"
)

type NotificationGenerator interface {
	Email() NotificationService
	Sms() NotificationService
	InApp() NotificationService
}

type Notifier struct {
	emailConfig struct {
	}
}

func (n *Notifier) Email() NotificationService {
	return &EmailNotificationService{}

}
func (n *Notifier) Sms() NotificationService {
	return &SmsNotificationService{}

}
func (n *Notifier) InApp() NotificationService {
	return &InAppNotificationService{}

}
func NewNotificationGenerator() NotificationGenerator {
	return &Notifier{}

}
