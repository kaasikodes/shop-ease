package providers

import (
	"context"
	"net/http"

	"github.com/kaasikodes/shop-ease/services/payment-service/internal/repository"
)

type PaystackGateway struct {
	apiKey string
	store  repository.PaymentRepo
}

func NewPaystackGateway(apiKey string, store repository.PaymentRepo) *PaystackGateway {
	return &PaystackGateway{apiKey, store}
}
func (p *PaystackGateway) InitiateTransaction(ctx context.Context, req PaymentRequest) (transactionID string, paymentUrl string, meta map[string]string, err error)
func (p *PaystackGateway) HandleWebhook(w http.ResponseWriter, r *http.Request) error
