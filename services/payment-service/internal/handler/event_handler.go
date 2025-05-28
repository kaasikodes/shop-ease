package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/kaasikodes/shop-ease/services/payment-service/internal/model"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/providers"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/env"
	"github.com/kaasikodes/shop-ease/shared/events"
)

type EventPayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}
type EventVendorSubscriptionPayload struct {
	UserId         int     `json:"userId"`
	VendorId       int     `json:"vendorId"`
	SubscriptionId int     `json:"subscriptionId"`
	Amount         float64 `json:"amount"`
}
type EventOrderPlacedPayload struct {
	UserId   int     `json:"userId"`
	VendorId int     `json:"vendorId"`
	OrderId  int     `json:"orderId"`
	Amount   float64 `json:"amount"`
}

type EventHandler struct {
	paymentRegistry map[model.PaymentProvider]providers.PaymentGateway
}

func InitEventHandler(store repository.PaymentRepo) *EventHandler {
	providers.RegisterProvider(model.PaymentProviderPaystack, providers.NewPaystackGateway(env.GetString("PAYSTACK_API_KEY", ""), store))
	providers.RegisterProvider(model.PaymentProviderFlutter, providers.NewFlutterkGateway(env.GetString("FLUTTER_API_KEY", ""), store))
	return &EventHandler{
		paymentRegistry: providers.ProviderRegistry,
	}

}

func (p *EventHandler) HandleSubscriptionEvents(msg []byte) error {
	// check the event type, and retrieve the msg convert to pay and then save the product
	var event EventPayload
	if err := json.Unmarshal(msg, &event.Data); err != nil {
		log.Printf("an error occured while unmarshaling the event: %v", err)
		return err
	}

	switch strings.ToLower(event.Event) {
	case events.VendorSubscriptionCreated:
		return p.payForVendorSubscription(msg)
	default:
		log.Printf("unhandled event type: %s", event.Event)

	}

	return nil

}
func (p *EventHandler) HandleOrderEvents(msg []byte) error {
	// check the event type, and retrieve the msg convert to pay and then save the product
	var event EventPayload
	if err := json.Unmarshal(msg, &event.Data); err != nil {
		log.Printf("an error occured while unmarshaling the event: %v", err)
		return err
	}

	switch strings.ToLower(event.Event) {
	case events.OrderCreated:
		return p.payForOrder(msg)
	default:
		log.Printf("unhandled event type: %s", event.Event)

	}

	return nil

}

func (p *EventHandler) payForOrder(msg []byte) error {
	var payload EventOrderPlacedPayload
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Printf("an error occured while unmarshaling the event payload: %v", err)
		return err
	}
	ctx := context.Background()
	p.paymentRegistry[model.PaymentProviderPaystack].InitiateTransaction(ctx, providers.PaymentRequest{
		Amount:     payload.Amount,
		EntityID:   strconv.Itoa(payload.OrderId),
		EntityType: model.EntityPaymentTypeVendorSubscriptionPayment,
		MetaData:   map[string]string{"vendorId": strconv.Itoa(payload.VendorId), "userId": strconv.Itoa(payload.UserId), "orderId": strconv.Itoa(payload.OrderId)},
	})
	return nil

}
func (p *EventHandler) payForVendorSubscription(msg []byte) error {
	var payload EventVendorSubscriptionPayload
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Printf("an error occured while unmarshaling the event payload: %v", err)
		return err
	}
	ctx := context.Background()
	p.paymentRegistry[model.PaymentProviderPaystack].InitiateTransaction(ctx, providers.PaymentRequest{
		Amount:     payload.Amount,
		EntityID:   strconv.Itoa(payload.SubscriptionId),
		EntityType: model.EntityPaymentTypeVendorSubscriptionPayment,
		MetaData:   map[string]string{"vendorId": strconv.Itoa(payload.VendorId), "userId": strconv.Itoa(payload.UserId), "subscriptionId": strconv.Itoa(payload.SubscriptionId)},
	})
	return nil

}
