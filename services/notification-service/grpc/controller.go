package grpc_server

import (
	"context"
	"fmt"
	"sync"

	"github.com/kaasikodes/shop-ease/services/notification-service/service"
	"github.com/kaasikodes/shop-ease/services/notification-service/store"
	"github.com/kaasikodes/shop-ease/shared/proto/notification"
	"google.golang.org/grpc"
)

// Current Flow
// Work on code to use air & Make file - Done
// Optimize Code here - Done
// Refactor to use go routines and take note of the time
// Also Flesh out the other service - sms, email (mail trap - and ensure you use templates to send the mail content), and in-app(use web sockets, and push notifications, as well as background jobs)
// Work on grpc security
// Remember to share project ...

// Next Project - should be an MCP.

type NotificationGrpcHandler struct {
	services []service.NotificationService
	notification.UnimplementedNotificationServiceServer
}

func NewNotificiationGRPCHandler(s *grpc.Server, services []service.NotificationService) {
	handler := &NotificationGrpcHandler{services: services}

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
	var errs []error
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(n.services))
	for _, v := range n.services {
		go func(service service.NotificationService) {
			defer wg.Done()
			err := service.Send(ctx, res)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}

		}(v)

	}
	wg.Wait()

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
	// If there were any errors, return them
	if len(errs) > 0 {
		// Return a single error that encapsulates all encountered errors
		return not, fmt.Errorf("errors encountered: %v", errs)
	}

	return not, nil
}
func (n *NotificationGrpcHandler) SendMultiple(ctx context.Context, payload *notification.SendMultipleRequest) (*notification.SendMultipleResponse, error) {
	nots := make([]store.Notification, len(payload.Notifications))
	for i, v := range payload.Notifications {
		phone := ""
		if v.Phone != nil {
			phone = *v.Phone

		}
		nots[i] = store.Notification{
			Email:   v.Email,
			Phone:   phone,
			Title:   v.Title,
			Content: v.Content,
		}

	}

	var errs []error
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(n.services))
	for _, v := range n.services {
		go func(service service.NotificationService) {
			defer wg.Done()
			err := v.SendMultiple(ctx, nots)
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}

		}(v)

	}
	wg.Wait()
	rNots := make([]*notification.Notification, len(nots))
	for i, v := range nots {
		rNots[i] = &notification.Notification{
			Id:        int32(v.ID),
			Email:     v.Email,
			Phone:     v.Phone,
			Title:     v.Title,
			Content:   v.Content,
			IsRead:    v.IsRead,
			ReadAt:    v.ReadAt.String(),
			CreatedAt: v.CreatedAt.String(),
			UpdatedAt: v.UpdatedAt.String(),
		}

	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("error encoutered: %v", errs)
	}
	return &notification.SendMultipleResponse{
		Notifications: rNots,
	}, nil
}
