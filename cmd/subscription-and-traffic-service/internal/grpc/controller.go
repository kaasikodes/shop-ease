package grpc_server

import (
	"context"

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

func (n *GrpcHandler) VerifyVendorSubscriptionStatus(ctx context.Context, payload *subscription.VerifyVendorSubscriptionStatusRequest) (*subscription.VerifyVendorSubscriptionStatusResponse, error)
func (n *GrpcHandler) MarkVendorSubscriptionAsPaid(ctx context.Context, payload *subscription.MarkVendorSubscriptionAsPaidRequest) (*subscription.VendorSubscription, error)
func (n *GrpcHandler) CreateVendorSubscription(ctx context.Context, payload *subscription.CreateVendorSubscriptionRequest) (*subscription.VendorSubscription, error)
