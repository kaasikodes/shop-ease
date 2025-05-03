package seller

import (
	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/types"
)

type Seller struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	types.Common
}
