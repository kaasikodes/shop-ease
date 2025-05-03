package products

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/kaasikodes/shop-ease/shared/events"
)

type EventPayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}
type EventProductPayload struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductEventHandler struct {
	products ProductRepo
}

func InitProductHandler(products ProductRepo) *ProductEventHandler {
	return &ProductEventHandler{
		products,
	}

}

func (p *ProductEventHandler) HandleProductEvents(msg []byte) error {
	// check the event type, and retrieve the msg convert to pay and then save the product
	var event EventPayload
	if err := json.Unmarshal(msg, &event.Data); err != nil {
		log.Printf("an error occured while unmarshaling the event: %v", err)
		return err
	}

	switch strings.ToLower(event.Event) {
	case events.ProductCreatedEvent, events.ProductUpdatedEvent:
		return p.updateProducts(msg)
	default:
		log.Printf("unhandled event type: %s", event.Event)

	}

	return nil

}
func (p *ProductEventHandler) HandleAuthEvents(msg []byte) error {
	// check the event type, and retrieve the msg convert to pay and then save the product
	var event EventPayload
	if err := json.Unmarshal(msg, &event.Data); err != nil {
		log.Printf("an error occured while unmarshaling the event: %v", err)
		return err
	}

	switch strings.ToLower(event.Event) {
	case events.ProductCreatedEvent, events.ProductUpdatedEvent:
		return p.updateProducts(msg)
	default:
		log.Printf("unhandled event type: %s", event.Event)

	}

	return nil

}

func (p *ProductEventHandler) updateProducts(msg []byte) error {
	var payload EventProductPayload
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Printf("an error occured while unmarshaling the event payload: %v", err)
		return err
	}
	p.products.Save(payload.ID, Product{
		ID:          payload.ID,
		Name:        payload.Name,
		Description: payload.Description,
	})
	return nil

}
