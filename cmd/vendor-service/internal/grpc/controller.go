package grpc_server

import (
	"context"

	"github.com/kaasikodes/shop-ease/services/vendor-service/internal/seller"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/proto/vendor_service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type Store struct {
	seller seller.SellerRepo
}
type GrpcHandler struct {
	store  Store
	trace  trace.Tracer
	logger logger.Logger

	vendor_service.UnimplementedVendorServiceServer
}

func NewGRPCHandler(s *grpc.Server, store Store, trace trace.Tracer, logger logger.Logger) {

	handler := &GrpcHandler{store: store, trace: trace, logger: logger}

	// register the VendorServiceServer
	vendor_service.RegisterVendorServiceServer(s, handler)

}

func (n *GrpcHandler) CreateVendor(ctx context.Context, payload *vendor_service.CreateVendorRequest) (*vendor_service.Vendor, error) {

	_, span := n.trace.Start(ctx, "Creating a vendor")
	defer span.End()
	n.logger.WithContext(ctx).Info("Creating a vendor starts")
	phone := ""
	if payload.Phone != nil {
		phone = *payload.Phone

	}
	request := seller.Seller{
		Email:  payload.Email,
		Phone:  phone,
		Name:   payload.Name,
		UserId: int(payload.UserId),
	}
	span.SetAttributes(
		attribute.Int("userId", int(payload.UserId)),
		attribute.String("email", payload.Email),
		attribute.String("name", payload.Name),
		attribute.String("phone", phone),
	)
	seller, err := n.store.seller.CreateVendor(request)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}
	vendor := &vendor_service.Vendor{
		Id:     int64(seller.ID),
		Phone:  seller.Phone,
		UserId: int64(seller.UserId),
		Email:  seller.Email,
		Name:   seller.Name,
	}
	span.SetAttributes(
		attribute.Int("vendorId", int(vendor.Id)),
	)
	n.logger.WithContext(ctx).Info("Created a vendor succesfully!")

	return vendor, nil

}
