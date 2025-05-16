package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/model"
	"github.com/kaasikodes/shop-ease/shared/types"
	"github.com/kaasikodes/shop-ease/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (app *application) getTransactionsHandler(w http.ResponseWriter, r *http.Request) {

	initialTraceCtx, span := app.trace.Start(r.Context(), "Get Payment Transactions")

	defer span.End()

	app.logger.WithContext(initialTraceCtx).Info("getting transactions")
	amountStr := r.URL.Query().Get("amount")
	provider := r.URL.Query().Get("provider")
	entityPaymentType := r.URL.Query().Get("entityPaymentType")
	status := r.URL.Query().Get("status")
	span.SetAttributes(
		attribute.String("filter.amount", amountStr),
		attribute.String("filter.provider", provider),
		attribute.String("filter.entityPaymentType", entityPaymentType),
		attribute.String("filter.status", status),
	)
	var amount int
	if amountStr != "" {
		amount, _ = strconv.Atoi(amountStr)

	}
	pagination := utils.GetPaginationFromQuery(r)

	transactions, total, err := app.store.GetTransactions(&types.PaginationPayload{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	}, &model.TransactionFilter{Provider: model.PaymentProvider(provider), Amount: float64(amount), EntityPaymentType: model.EntityPaymentType(entityPaymentType), Status: model.PaymentStatus(status)})
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error getting transactions", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed getting transactions")

	var result = make([]any, len(transactions))
	for i, product := range transactions {
		result[i] = product

	}

	app.jsonResponse(w, http.StatusOK, "Payment Transactions retrieved successfully!", createPaginatedResponse(result, total))
	return

}
func determinePaymentProvider(r *http.Request) (model.PaymentProvider, error) {
	// Case 1 -> From Header
	if provider := r.Header.Get("X-Payment-Provider"); provider != "" {
		return model.PaymentProvider(provider), nil
	}

	// Case 2 -> Infer from unique JSON fields
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read request body")
	}
	defer r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(body)) // Reset

	if bytes.Contains(body, []byte(`"event":"charge.success"`)) {
		return model.PaystackPaymentProvider, nil
	} else if bytes.Contains(body, []byte(`"data":"flutter_event"`)) {
		return model.FlutterPaymentProvider, nil
	}

	return "", fmt.Errorf("could not determine payment provider")
}

func (app *application) webHookHandler(w http.ResponseWriter, r *http.Request) {
	initialTraceCtx, span := app.trace.Start(r.Context(), "Get Payment Transactions")

	defer span.End()
	provider, err := determinePaymentProvider(r)
	span.SetAttributes(attribute.String("provider", string(provider)))
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error determining payment provider", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.badRequestResponse(w, r, err)
		return
	}
	handler, ok := app.paymentRegistry[provider]
	if !ok {
		app.logger.WithContext(initialTraceCtx).Error("Unregistered payment provider:", provider)
		span.SetStatus(codes.Error, "Provider not registered")
		app.badRequestResponse(w, r, fmt.Errorf("unsupported payment provider"))
		return
	}
	err = handler.HandleWebhook(w, r)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error handling webhook", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	span.SetStatus(codes.Ok, "Webhook processed")
	app.jsonResponse(w, http.StatusOK, "Payment web hook processed successfully!", nil)
	return
}
func (app *application) getTransactionByIdHandler(w http.ResponseWriter, r *http.Request) {

	initialTraceCtx, span := app.trace.Start(r.Context(), "Get Transaction")

	defer span.End()
	transactionIdStr := chi.URLParam(r, "transactionId")
	transactionId, err := strconv.Atoi(transactionIdStr)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error reading transactionId from url", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	span.SetAttributes(
		attribute.Int("transactionId", transactionId),
	)

	app.logger.WithContext(initialTraceCtx).Info("getting single transaction")

	transaction, err := app.store.GetTransactionById(transactionId)
	if err != nil {
		app.logger.WithContext(initialTraceCtx).Error("Error getting transaction", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		app.internalServerError(w, r, err)
		return
	}
	app.logger.WithContext(initialTraceCtx).Info("Completed getting inventory for vendor/seller")

	app.jsonResponse(w, http.StatusOK, "Transaction retrieved successfully!", transaction)
	return

}
