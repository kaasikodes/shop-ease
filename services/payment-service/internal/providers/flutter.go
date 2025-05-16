package providers

import (
	"context"
	"net/http"

	"github.com/kaasikodes/shop-ease/services/payment-service/internal/repository"
)

type FlutterGateway struct {
	apiKey string
	store  repository.PaymentRepo
}

func NewFlutterkGateway(apiKey string, store repository.PaymentRepo) *FlutterGateway {
	return &FlutterGateway{apiKey, store}
}
func (p *FlutterGateway) InitiateTransaction(ctx context.Context, req PaymentRequest) (transactionID string, meta map[string]string, err error)
func (p *FlutterGateway) HandleWebhook(w http.ResponseWriter, r *http.Request) error
