package repository

import (
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/model"
	"github.com/kaasikodes/shop-ease/shared/types"
)

type PaymentRepo interface {
	GetTransactions(pagination *types.PaginationPayload, filter *model.TransactionFilter) (result []model.Transaction, total int, err error)
	CreateTransaction(model.Transaction) (data *model.Transaction, err error)
	UpdateTransaction(id int, payload model.Transaction) (data *model.Transaction, err error)
	GetTransactionById(id int) (data *model.Transaction, err error)
}
