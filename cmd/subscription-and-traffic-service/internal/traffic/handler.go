package traffic

import (
	"encoding/json"
	"log"
	"strings"

	vendorplan "github.com/kaasikodes/shop-ease/services/subscription-and-traffic-service/internal/vendor-plan"
	"github.com/kaasikodes/shop-ease/shared/events"
)

type EventPayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}
type UserInteractionPayload struct {
	VendorId int                                  `json:"vendorId"`
	UserId   int                                  `json:"userId"`
	Type     vendorplan.VendorUserInteractionType `json:"type"`
}

type EventHandler struct {
	plan vendorplan.VendorPlanRepo
}

func InitEventHandler(plan vendorplan.VendorPlanRepo) *EventHandler {
	return &EventHandler{
		plan,
	}

}

func (p *EventHandler) HandleAuthEvents(msg []byte) error {
	// check the event type, and retrieve the msg convert to pay and then save the product
	var event EventPayload
	if err := json.Unmarshal(msg, &event.Data); err != nil {
		log.Printf("an error occured while unmarshaling the event: %v", err)
		return err
	}

	switch strings.ToLower(event.Event) {
	case events.UserOrderedItemEvent, events.UserInterestedInItemEvent:
		return p.saveInteraction(msg)
	default:
		log.Printf("unhandled event type: %s", event.Event)

	}

	return nil

}

func (p *EventHandler) saveInteraction(msg []byte) error {
	var payload UserInteractionPayload
	if err := json.Unmarshal(msg, &payload); err != nil {
		log.Printf("an error occured while unmarshaling the event payload: %v", err)
		return err
	}
	_, err := p.plan.CreateVendorUserInteractionRecord(payload.UserId, payload.Type)
	return err

}
