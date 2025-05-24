package main

import (
	"context"
	"time"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	jwttoken "github.com/kaasikodes/shop-ease/shared/jwt_token"
)

const (
	ExpiresAtVerificationToken = time.Hour * 24 * 5
	AccessTokenDuration        = time.Duration(time.Hour * 24 * 3)
)

type ContextKeyUser struct{}
type ContextKeyClaims struct{}

type paginatedResponse struct {
	Total  int   `json:"total"`
	Result []any `json:"result"`
}

func createPaginatedResponse(result []any, total int) paginatedResponse {
	return paginatedResponse{
		Total:  total,
		Result: result,
	}

}

func (app *application) isProduction() bool {
	return app.config.env == "production"
}
func getUserFromContext(ctx context.Context) (*store.User, bool) {
	user, ok := ctx.Value(ContextKeyUser{}).(*store.User)
	return user, ok
}

func getClaimsFromContext(ctx context.Context) (*jwttoken.CustomClaims, bool) {
	claims, ok := ctx.Value(ContextKeyClaims{}).(*jwttoken.CustomClaims)
	return claims, ok
}
