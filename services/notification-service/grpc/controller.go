package main

import (
	"context"

	"github.com/kaasikodes/shop-ease/services/notification-service/service"
	"github.com/kaasikodes/shop-ease/services/notification-service/store"
	"github.com/kaasikodes/shop-ease/shared/proto/notification"
	"google.golang.org/grpc"
)

type NotificationGrpcHandler struct {
	service []service.NotificationService
	notification.UnimplementedNotificationServiceServer
}

func NewNotificiationGRPCHandler(s *grpc.Server, service []service.NotificationService) {
	handler := &NotificationGrpcHandler{service: service}

	// register the NotificationServiceServer
	notification.RegisterNotificationServiceServer(s, handler)

}

func (n *NotificationGrpcHandler) Send(ctx context.Context, payload *notification.NotificationRequest) (*notification.Notification, error) {
	phone := ""
	if payload.Phone != nil {
		phone = *payload.Phone

	}
	res := &store.Notification{
		Email:   payload.Email,
		Phone:   phone,
		Title:   payload.Title,
		Content: payload.Content,
	}
	var err error
	for _, v := range n.service {
		err = v.Send(ctx, res)
		if err != nil {
			return nil, err
		}

	}

	not := &(notification.Notification{
		Id:        int32(res.ID),
		Email:     res.Email,
		Phone:     res.Phone,
		Title:     res.Title,
		Content:   res.Content,
		IsRead:    res.IsRead,
		ReadAt:    res.ReadAt.String(),
		CreatedAt: res.CreatedAt.String(),
		UpdatedAt: res.UpdatedAt.String(),
	})
	return not, err
}
func (n *NotificationGrpcHandler) SendMultiple(ctx context.Context, payload *notification.SendMultipleRequest) (*notification.SendMultipleResponse, error) {
	nots := []store.Notification{}
	for _, v := range payload.Notifications {
		phone := ""
		if v.Phone != nil {
			phone = *v.Phone

		}
		nots = append(nots, store.Notification{
			Email:   v.Email,
			Phone:   phone,
			Title:   v.Title,
			Content: v.Content,
		})

	}

	var err error
	for _, v := range n.service {
		err = v.SendMultiple(ctx, nots)
		if err != nil {
			return nil, err
		}

	}
	rNots := []*notification.Notification{}
	for _, v := range nots {
		rNots = append(rNots, &notification.Notification{
			Id:        int32(v.ID),
			Email:     v.Email,
			Phone:     v.Phone,
			Title:     v.Title,
			Content:   v.Content,
			IsRead:    v.IsRead,
			ReadAt:    v.ReadAt.String(),
			CreatedAt: v.CreatedAt.String(),
			UpdatedAt: v.UpdatedAt.String(),
		})

	}
	return &notification.SendMultipleResponse{
		Notifications: rNots,
	}, err
}
