package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kaasikodes/shop-ease/services/order-service/internal/model"
	"github.com/kaasikodes/shop-ease/services/order-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/utils"
	"go.opentelemetry.io/otel/codes"
)

type changeStatusPayload struct {
	Status string `json:"status"`
}

func (app *application) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "http.order.getOrders")
	defer span.End()

	var q = r.URL.Query()

	filter := repository.OrderFilter{
		Status:    model.OrderStatus(q.Get("status")),
		UserId:    utils.ParseInt(q.Get("user_id")),
		StoreId:   utils.ParseInt(q.Get("store_id")),
		ProductId: utils.ParseInt(q.Get("product_id")),
	}

	pagination := utils.GetPaginationFromQuery(r)

	orders, total, err := app.store.GetOrders(ctx, pagination, &filter)
	if err != nil {
		app.logger.Error("getOrders failed", err)
		app.internalServerError(w, r, err)
		return
	}
	var result = make([]any, len(orders))
	for i, o := range orders {
		result[i] = o

	}

	app.jsonResponse(w, http.StatusOK, "Orders retrieved successfully!", createPaginatedResponse(result, total))

}

func (app *application) getOrderByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "http.order.getOrderById")
	defer span.End()

	orderId := chi.URLParam(r, "orderId")
	id := utils.ParseInt(orderId)
	if id == 0 {

		app.badRequestResponse(w, r, errors.New("invalid order id"))
		return
	}

	order, err := app.store.GetOrderById(ctx, id)
	if err != nil {
		app.logger.Error("getOrderById failed", err)
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Order retrieved successfully!", order)

}

func (app *application) changeOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "http.order.changeOrderStatus")
	defer span.End()

	orderId := utils.ParseInt(chi.URLParam(r, "orderId"))
	if orderId == 0 {
		app.badRequestResponse(w, r, errors.New("invalid order id"))

		return
	}

	var payload changeStatusPayload
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.UpdateOrderStatus(ctx, orderId, model.OrderStatus(payload.Status))
	if err != nil {
		app.logger.Error("UpdateOrderStatus failed", err)
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Order status updated successfully!", nil)

}

func (app *application) changeOrderItemStatusHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := app.trace.Start(r.Context(), "http.order.changeOrderItemStatus")
	defer span.End()

	orderItemId := utils.ParseInt(chi.URLParam(r, "orderItemId"))
	if orderItemId == 0 {
		app.badRequestResponse(w, r, errors.New("invalid item id"))

		return
	}

	var payload changeStatusPayload
	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.store.UpdateOrderItemStatus(ctx, orderItemId, model.OrderStatus(payload.Status))
	if err != nil {
		app.logger.Error("UpdateOrderItemStatus failed", err)
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Order item status updated successfully!", nil)

}

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := getUserIdFromContext(ctx)
	ctx, span := app.trace.Start(ctx, "http.order.createOrder")
	defer span.End()
	if !ok {
		app.badRequestResponse(w, r, errors.New("unable to retrieve userId"))
		return

	}
	app.logger.Info("userId", userId)

	var body struct {
		Items []repository.CreateOrderInputItem `json:"items" validate:"min=1,dive,required"`
	}

	err := app.readJSON(w, r, &body)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(body); err != nil {
		app.logger.WithContext(ctx).Error("Error validating registration payload", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	orderId, err := app.store.CreateOrder(ctx, userId, body.Items)
	if err != nil {
		app.logger.Error("CreateOrder failed", err)
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, "Order created successfully!", map[string]any{"orderId": orderId, "itemCount": len(body.Items)})

}
