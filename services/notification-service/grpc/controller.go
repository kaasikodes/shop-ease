package grpc_server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/kaasikodes/shop-ease/services/notification-service/service"
	"github.com/kaasikodes/shop-ease/services/notification-service/store"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/proto/notification"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Plan of Action
// - Monitoring/Observability/Alerting - Prometheus, Loki, Grafana, OpenTelemetry, Moving to loki in grafana integration with, also see how to create a trace that has spans across services - Done
// Need to configure all log implemetations to have same format, zap to store logs in same file, and also a time of expiry to clear out the content of the file - Done(but zap, and default have different formats-probably can ensure same format by updating format to match zap)
// Once done ensure that there is a standard of observability - metrics, logs, traces across services. Also spend a bit more time building relevant dashboards on grafana, configuring alerts, and how the setup of the dashboards can be reused or shared across projects/with individuals
// Also don't forget to link aync operations like events into traces as well
// Refactor the depecrated grpc interceptors used in traceing
// - Implement Service 2 service communication with Rabbitmq and GRPC
// - Implement all services
// - Build the necessary dashboards for the services
// - Refactor app to use K8s and kubernetes
// - Also deploy app on aws with cloudformation, and then use terraform as well
// - Purchase server and deploy - use jenkins here
// - Once completed write a couple of articles: 1st article should be on the architecture
// - Share with relevant groups - discord server project showcase, whatsapp group, linkedin ...

// Current Flow
// Work on code to use air & Make file - Done
// Optimize Code here - Done
// ** Side Quests- Play around with jenkins for ci/cd
// ELK Stack
// ArgoCD
// gitlab
// Refactor to use go routines and take note of the time - Done, ignored time bench mark for the time being
// Current flow - prometheus, grafana, loki, opentelemetry(Begin Here), ElastiSearch-kibana?, Circuit breaker, service discovery(with plain api-gateway and no kubernetes), event-driven architecture (rabbitmq, kafka - agnostic is it possible), logging ..., grpc for streaming videos ... kubenetes, open source daily(the project picked earlier), MongoDB/Express/Service - say audit service, GraphQL - product service, elastisearch, AWS deployment after local setup via cloud formation, use lamda function to move log file though s3 storage ...,
// Start with the email service, also might need to refactor the use of the notification handler for grpc as not all actions will require a all types of notification -> registration just requires email, and not inapp , might need sms
// Also Flesh out the other service - sms(will require payment so skip for now), email (mail trap - and ensure you use templates to send the mail content - Done, but not using html for time being.), and in-app(use web sockets, and push notifications, as well as background jobs)
// Work on grpc security

// Next Project - should be an MCP.

type NotificationGrpcHandler struct {
	services []service.NotificationService
	trace    trace.Tracer
	logger   logger.Logger

	notification.UnimplementedNotificationServiceServer
}

func NewNotificiationGRPCHandler(s *grpc.Server, services []service.NotificationService, trace trace.Tracer, logger logger.Logger) {

	handler := &NotificationGrpcHandler{services: services, trace: trace, logger: logger}

	// register the NotificationServiceServer
	notification.RegisterNotificationServiceServer(s, handler)

}

func (n *NotificationGrpcHandler) Send(ctx context.Context, payload *notification.NotificationRequest) (*notification.Notification, error) {
	// md, ok := metadata.FromIncomingContext(ctx)
	// log.Println("Incoming", md)
	// if !ok {
	// 	return nil, errors.New("issue retrieving metdata from context")
	// }
	// ctx = observability.Propagator.Extract(ctx, propagation.HeaderCarrier(md))
	parentNotificationCtx, span := n.trace.Start(ctx, "send notification")
	defer span.End()
	n.logger.WithContext(ctx).Info("send notification starts")
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
	log.Println("Send Starts herere")
	var errs []error
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(n.services))
	for i, v := range n.services {
		go func(service service.NotificationService) {
			defer wg.Done()
			err := service.Send(ctx, res)
			_, span := n.trace.Start(parentNotificationCtx, fmt.Sprintf("sending notification by %v service", i))
			defer span.End()
			if err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}

		}(v)

	}
	wg.Wait()

	log.Println("Send ends herere bvbvbvbv")
	not := &notification.Notification{
		Id:      int32(res.ID),
		Email:   res.Email,
		Phone:   res.Phone,
		Title:   res.Title,
		Content: res.Content,
		IsRead:  res.IsRead,
	}
	// If there were any errors, return them
	if len(errs) > 0 {
		errs := errors.Join(errs...)
		// Return a single error that encapsulates all encountered errors
		span.RecordError(errs)
		span.SetStatus(codes.Error, errs.Error())
		return not, errs
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
