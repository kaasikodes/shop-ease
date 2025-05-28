package providers

import (
	"context"
	"net/http"

	"github.com/kaasikodes/shop-ease/services/payment-service/internal/model"
)

type PaymentRequest struct {
	Amount     float64
	EntityID   string
	EntityType model.EntityPaymentType
	MetaData   map[string]string
}

type PaymentGateway interface {
	InitiateTransaction(ctx context.Context, req PaymentRequest) (transactionID string, paymentUrl string, meta map[string]string, err error)
	HandleWebhook(w http.ResponseWriter, r *http.Request) error
}

var ProviderRegistry = make(map[model.PaymentProvider]PaymentGateway)

func RegisterProvider(providerType model.PaymentProvider, gateway PaymentGateway) {
	ProviderRegistry[providerType] = gateway

}
