package handler

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/kaasikodes/shop-ease/services/order-service/internal/model"
	"github.com/kaasikodes/shop-ease/services/order-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/events"
)

type EventPayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}
type EventVendorAcceptedOrderItemPayload struct {
	OrderItemId int `json:"orderItemId"`
}
type EventHandler struct {
	store repository.OrderRepo
}

func InitEventHandler(store repository.OrderRepo) *EventHandler {

	return &EventHandler{
		store,
	}

}

func (p *EventHandler) HandleVendorEvents(msg []byte) error {
	// check the event type, and retrieve the msg convert to pay and then save the product
	var event EventPayload
	if err := json.Unmarshal(msg, &event.Data); err != nil {
		log.Printf("an error occured while unmarshaling the event: %v", err)
		return err
	}

	switch strings.ToLower(event.Event) {
	case events.VendorAcceptedOrderItem:
		return p.vendorAccepetedOrderItem(msg)
	default:
		log.Printf("unhandled event type: %s", event.Event)

	}

	return nil

}

func (p *EventHandler) vendorAccepetedOrderItem(msg []byte) error {
	var payload EventVendorAcceptedOrderItemPayload
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Printf("an error occured while unmarshaling the event payload: %v", err)
		return err
	}
	ctx := context.Background()
	err := p.store.UpdateOrderItemStatus(ctx, payload.OrderItemId, model.ProcessingOrderStatus)
	return err

}
