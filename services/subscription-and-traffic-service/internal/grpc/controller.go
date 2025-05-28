package grpc_server

import (
	"context"
	"errors"
	"time"

	vendorplan "github.com/kaasikodes/shop-ease/services/subscription-and-traffic-service/internal/vendor-plan"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/proto/subscription"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type Store struct {
	plan vendorplan.VendorPlanRepo
}
type GrpcHandler struct {
	store  Store
	trace  trace.Tracer
	logger logger.Logger

	subscription.UnimplementedSubscriptionServiceServer
}

func NewGRPCHandler(s *grpc.Server, store Store, trace trace.Tracer, logger logger.Logger) {

	handler := &GrpcHandler{store: store, trace: trace, logger: logger}

	// register the SubscriptionServiceServer
	subscription.RegisterSubscriptionServiceServer(s, handler)

}

// func (n *GrpcHandler) VerifyVendorSubscriptionStatus(ctx context.Context, payload *subscription.VerifyVendorSubscriptionStatusRequest) (*subscription.VerifyVendorSubscriptionStatusResponse, error)
// func (n *GrpcHandler) MarkVendorSubscriptionAsPaid(ctx context.Context, payload *subscription.MarkVendorSubscriptionAsPaidRequest) (*subscription.VendorSubscription, error)
// func (n *GrpcHandler) CreateVendorSubscription(ctx context.Context, payload *subscription.CreateVendorSubscriptionRequest) (*subscription.VendorSubscription, error)

func (n *GrpcHandler) VerifyVendorSubscriptionStatus(ctx context.Context, payload *subscription.VerifyVendorSubscriptionStatusRequest) (*subscription.VerifyVendorSubscriptionStatusResponse, error) {
	subscriptions, err := n.store.plan.GetActiveSubscriptionsForVendor(payload.VendorId)
	if err != nil {
		n.logger.Error("failed to get subscriptions", err)
		return nil, err
	}

	now := time.Now()

	// This example assumes you have a way to get user interaction count for the vendor.
	// Here we just do pseudo logic:
	for _, sub := range subscriptions {
		// Check if expired
		if sub.ExpiresAt.Before(now) {
			continue
		}
		// Check if limit exceeded already marked
		if !sub.LimitExceededAt.IsZero() {
			return &subscription.VerifyVendorSubscriptionStatusResponse{
				IsValid: false,
				Message: "Subscription limit has been exceeded",
			}, nil
		}
		// totalUserInteractionsUsed += ... // fetch from your logic or DB
	}

	return &subscription.VerifyVendorSubscriptionStatusResponse{
		IsValid: true,
		Message: "Subscription is valid",
	}, nil
}

func (n *GrpcHandler) MarkVendorSubscriptionAsPaid(ctx context.Context, payload *subscription.MarkVendorSubscriptionAsPaidRequest) (*subscription.VendorSubscription, error) {
	// Find subscription by transaction id (assuming transactionId maps to subscription)
	sub, err := n.store.plan.GetVendorSubscriptionID(payload.TransactionId) //correct the TrasactionId to a better name
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("subscription not found for transaction")
	}

	err = n.store.plan.MarkSubscriptionPaid(int64(sub.ID))
	if err != nil {
		return nil, err
	}

	return &subscription.VendorSubscription{
		Id:       int64(sub.ID),
		PlanId:   int64(sub.PlanId),
		VendorId: int64(sub.VendorId),
	}, nil
}

func (n *GrpcHandler) CreateVendorSubscription(ctx context.Context, payload *subscription.CreateVendorSubscriptionRequest) (*subscription.VendorSubscription, error) {
	sub, err := n.store.plan.CreateVendorPlanSubscription(int(payload.PlanId), int(payload.VendorId))
	if err != nil {
		return nil, err
	}

	return &subscription.VendorSubscription{
		Id:       int64(sub.ID),
		PlanId:   int64(sub.PlanId),
		VendorId: int64(sub.VendorId),
	}, nil
}
