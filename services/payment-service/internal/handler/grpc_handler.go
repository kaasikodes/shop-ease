package handler

import (
	"context"

	"github.com/kaasikodes/shop-ease/services/payment-service/internal/model"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/types"

	"github.com/kaasikodes/shop-ease/shared/proto/payment"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type PaymentGrpcHandler struct {
	trace  trace.Tracer
	logger logger.Logger
	store  repository.PaymentRepo
	payment.UnimplementedPaymentServiceServer
}

func NewPaymentGRPCHandler(s *grpc.Server, store repository.PaymentRepo, trace trace.Tracer, logger logger.Logger) {

	handler := &PaymentGrpcHandler{trace: trace, logger: logger, store: store}

	// register the NotificationServiceServer
	payment.RegisterPaymentServiceServer(s, handler)

}

func (n *PaymentGrpcHandler) GetTransactions(ctx context.Context, payload *payment.GetTransactionsRequest) (*payment.TransactionList, error) {

	_, span := n.trace.Start(ctx, "retrieving transactions")
	defer span.End()
	n.logger.WithContext(ctx).Info("retrieving transactions starts")

	// get transactions from store
	result, total, err := n.store.GetTransactions(&types.PaginationPayload{Limit: int(payload.Pagination.Limit), Offset: int(payload.Pagination.Offset)}, &model.TransactionFilter{Status: model.PaymentStatus(payload.Filter.Status), Provider: model.PaymentProvider(payload.Filter.Provider)})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}
	// convert transactions to grpc transactions
	transactions := make([]*payment.Transaction, len(result))
	for i, row := range result {
		transactions[i].EntityId = int64(row.EntityId)
		transactions[i].EntityPaymentType = string(row.EntityPaymentType)
		transactions[i].Id = int64(row.ID)
		transactions[i].Provider = string(row.Provider)
		transactions[i].MetaData = row.MetaData
		transactions[i].Status = string(row.Status)

	}

	n.logger.WithContext(ctx).Info("retrieving transactions ends")
	return &payment.TransactionList{Total: int64(total), Transactions: transactions}, nil
}
func (n *PaymentGrpcHandler) GetTransactionById(ctx context.Context, payload *payment.GetByIdRequest) (*payment.Transaction, error) {
	_, span := n.trace.Start(ctx, "retrieving transactions")
	defer span.End()
	n.logger.WithContext(ctx).Info("retrieving transactions starts")
	data, err := n.store.GetTransactionById(int(payload.TransactionId))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &payment.Transaction{
		Id:                int64(data.ID),
		EntityId:          int64(data.EntityId),
		Status:            string(data.Status),
		EntityPaymentType: string(data.EntityPaymentType),
		Provider:          string(data.Provider),
		MetaData:          data.MetaData,
	}, nil

}
