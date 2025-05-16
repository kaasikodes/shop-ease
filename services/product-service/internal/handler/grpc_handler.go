package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/kaasikodes/shop-ease/services/product-service/internal/model"
	"github.com/kaasikodes/shop-ease/services/product-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/types"
	"github.com/kaasikodes/shop-ease/shared/utils"

	"github.com/kaasikodes/shop-ease/shared/proto/product"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type ProductGrpcHandler struct {
	trace  trace.Tracer
	logger logger.Logger
	store  repository.ProductRepo
	product.UnimplementedProductServiceServer
}

func NewProductGrpcHandler(s *grpc.Server, store repository.ProductRepo, trace trace.Tracer, logger logger.Logger) {

	handler := &ProductGrpcHandler{trace: trace, logger: logger, store: store}

	// register the ProductServiceServer
	product.RegisterProductServiceServer(s, handler)

}
func (n *ProductGrpcHandler) CreateDiscount(ctx context.Context, req *product.CreateDiscountRequest) (*product.Discount, error) {
	parentCtx, span := n.trace.Start(ctx, "creating discount")
	defer span.End()
	n.logger.WithContext(ctx).Info("creating discount starts")

	// Parse time fields
	effectiveAt, err := time.Parse(time.RFC3339, req.EffectiveAt)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid effectiveAt format")
		return nil, fmt.Errorf("invalid effectiveAt format: %w", err)
	}

	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "invalid expiresAt format")
			return nil, fmt.Errorf("invalid expiresAt format: %w", err)
		}
		expiresAt = &t
	}

	// Map string enums to internal types
	valueType := types.DiscountValueType(req.ValueType)
	paidBy := types.PaidBy(req.PaidBy)

	// Map to internal payload
	discount := &model.Discount{
		ValueType:   valueType,
		EffectiveAt: effectiveAt,
		ExpiresAt:   expiresAt,
		PaidBy:      paidBy,
		Value:       int16(req.Value),

		CommonDescriptiveModel: types.CommonDescriptiveModel{Name: req.Name, Description: req.Description},
		// Description: req.Description,
		ApplicableTo: types.DiscountApplicability{
			ProductIds:               req.ApplicableTo.ProductIds,
			StoreProductIds:          req.ApplicableTo.StoreProductIds,
			StoreProductInventoryIds: req.ApplicableTo.StoreProductInventoryIds,
		},
	}

	// Call repository
	err = n.store.CreateDiscount(parentCtx, *discount)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Map response
	res := &product.Discount{
		Id:          int64(discount.ID),
		Value:       int32(discount.Value),
		ValueType:   string(discount.ValueType),
		EffectiveAt: discount.EffectiveAt.Format(time.RFC3339),
		ExpiresAt:   "",
		PaidBy:      string(discount.PaidBy),
		Name:        discount.Name,
		Description: discount.Description,
		CreatedAt:   discount.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   discount.UpdatedAt.Format(time.RFC3339),
		ApplicableTo: &product.DiscountApplicability{
			ProductIds:               discount.ApplicableTo.ProductIds,
			StoreProductIds:          discount.ApplicableTo.StoreProductIds,
			StoreProductInventoryIds: discount.ApplicableTo.StoreProductInventoryIds,
		},
	}

	if discount.ExpiresAt != nil {
		res.ExpiresAt = discount.ExpiresAt.Format(time.RFC3339)
	}

	n.logger.WithContext(ctx).Info("creating discount ends")
	return res, nil
}

func (n *ProductGrpcHandler) GetDiscounts(ctx context.Context, req *product.GetDiscountsRequest) (*product.DiscountList, error) {
	parentCtx, span := n.trace.Start(ctx, "GetDiscounts")
	defer span.End()
	n.logger.WithContext(ctx).Info("Getting discounts")

	// Step 1: Build pagination and filter payloads
	pagination := &utils.PaginationPayload{
		Limit:  int(req.Pagination.Limit),
		Offset: int(req.Pagination.Offset),
	}

	var expiresAt *time.Time
	if req.Filter.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.Filter.ExpiresAt)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "invalid expiresAt format")
			return nil, fmt.Errorf("invalid expiresAt format: %w", err)
		}
		expiresAt = &t
	}

	filter := &repository.DiscountFilter{

		ExpiresAt: expiresAt,
		Applicability: types.DiscountApplicability{
			ProductIds:               req.Filter.ApplicableTo.ProductIds,
			StoreProductIds:          req.Filter.ApplicableTo.StoreProductIds,
			StoreProductInventoryIds: req.Filter.ApplicableTo.StoreProductInventoryIds,
		},
	}

	// Step 2: Query discounts from store
	discounts, total, err := n.store.GetDiscounts(parentCtx, pagination, filter)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Step 3: Map to proto response
	resp := make([]*product.Discount, 0, len(discounts))
	for _, d := range discounts {
		var expiresAt string
		if d.ExpiresAt != nil {
			expiresAt = d.ExpiresAt.Format(time.RFC3339)
		}

		resp = append(resp, &product.Discount{
			Id:          int64(d.Id),
			Value:       int32(d.Value),
			ValueType:   string(d.ValueType),
			EffectiveAt: d.EffectiveAt.Format(time.RFC3339),
			ExpiresAt:   expiresAt,
			PaidBy:      string(d.PaidBy),
			Name:        d.Name,
			Description: d.Description,
			ApplicableTo: &product.DiscountApplicability{
				ProductIds:               d.ApplicableTo.ProductIds,
				StoreProductIds:          d.ApplicableTo.StoreProductIds,
				StoreProductInventoryIds: d.ApplicableTo.StoreProductInventoryIds,
			},
			CreatedAt: d.CreatedAt.Format(time.RFC3339),
			UpdatedAt: d.UpdatedAt.Format(time.RFC3339),
		})
	}

	n.logger.WithContext(ctx).Info("Discount retrieval completed")
	return &product.DiscountList{
		Discounts: resp,
		Total:     int64(total),
	}, nil
}
