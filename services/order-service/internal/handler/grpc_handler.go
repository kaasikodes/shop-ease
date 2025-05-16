package handler

import (
	"context"

	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/utils"

	"github.com/kaasikodes/shop-ease/services/order-service/internal/model"
	"github.com/kaasikodes/shop-ease/services/order-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/proto/order"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGrpcHandler struct {
	trace  trace.Tracer
	logger logger.Logger
	store  repository.OrderRepo
	order.UnimplementedOrderServiceServer
}

func NewOrderGrpcHandler(s *grpc.Server, store repository.OrderRepo, trace trace.Tracer, logger logger.Logger) {

	handler := &OrderGrpcHandler{trace: trace, logger: logger, store: store}

	// register the ProductServiceServer
	order.RegisterOrderServiceServer(s, handler)

}
func (h *OrderGrpcHandler) GetOrders(ctx context.Context, req *order.GetOrdersRequest) (*order.GetOrdersResponse, error) {
	ctx, span := h.trace.Start(ctx, "OrderGrpcHandler.GetOrders")
	defer span.End()

	// Pagination setup
	pagination := &utils.PaginationPayload{
		Limit:  int(req.Pagination.Limit),
		Offset: int(req.Pagination.Offset),
	}

	// Filter setup
	filter := &repository.OrderFilter{
		Status:    model.OrderStatus(req.Filter.Status),
		UserId:    int(req.Filter.UserId),
		StoreId:   int(req.Filter.UserId),
		ProductId: int(req.Filter.ProductId),
	}

	orders, total, err := h.store.GetOrders(ctx, pagination, filter)
	if err != nil {
		h.logger.Error("failed to get orders", err)
		return nil, status.Errorf(codes.Internal, "could not get orders: %v", err)
	}

	// Map to proto response
	var protoOrders []*order.OrderListItem
	for _, o := range orders {
		protoOrders = append(protoOrders, &order.OrderListItem{
			Id:         int32(o.Id),
			UserId:     int32(o.UserId),
			Status:     string(o.Status),
			ItemCount:  int32(o.ItemCount),
			IsPaid:     o.IsPaid,
			IsCanceled: o.IsCanceled,
			PaidAt:     o.PaidAt.String(),
			CanceledAt: o.CanceledAt.String(),
			CreatedAt:  o.Common.CreatedAt.String(),
			UpdatedAt:  o.Common.UpdatedAt.String(),
		})
	}

	return &order.GetOrdersResponse{
		Orders: protoOrders,
		Total:  int32(total),
	}, nil
}

func (h *OrderGrpcHandler) GetOrderById(ctx context.Context, req *order.GetOrderByIdRequest) (*order.GetOrderByIdResponse, error) {
	ctx, span := h.trace.Start(ctx, "OrderGrpcHandler.GetOrderById")
	defer span.End()

	ord, err := h.store.GetOrderById(ctx, int(req.OrderId))
	if err != nil {
		h.logger.Error("failed to get order by id", err)
		return nil, status.Errorf(codes.Internal, "could not get order: %v", err)
	}

	// Build order items
	var items []*order.OrderItem
	for _, item := range ord.Items {
		items = append(items, &order.OrderItem{
			Id:        int32(item.Id),
			ProductId: int32(item.ProductId),
			StoreId:   int32(item.StoreId),
			Quantity:  int32(item.Quantity),
			CreatedAt: item.CreatedAt.String(),
			UpdatedAt: item.UpdatedAt.String(),
		})
	}

	res := &order.GetOrderByIdResponse{
		Order: &order.Order{
			Id:         int32(ord.Id),
			UserId:     int32(ord.UserId),
			Status:     string(ord.Status),
			IsPaid:     ord.IsPaid,
			IsCanceled: ord.IsCanceled,
			PaidAt:     ord.PaidAt.String(),
			CanceledAt: (ord.CanceledAt.String()),
			CreatedAt:  ord.CreatedAt.String(),
			UpdatedAt:  ord.UpdatedAt.String(),
			Items:      items,
		},
	}

	return res, nil
}
