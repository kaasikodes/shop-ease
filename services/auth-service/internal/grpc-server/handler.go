package grpc_server

import (
	"context"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	"github.com/kaasikodes/shop-ease/shared/logger"

	"github.com/kaasikodes/shop-ease/shared/proto/auth"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type AuthGrpcHandler struct {
	trace  trace.Tracer
	logger logger.Logger
	store  store.Storage
	auth.UnimplementedAuthServiceServer
}

func NewAuthGRPCHandler(s *grpc.Server, store store.Storage, trace trace.Tracer, logger logger.Logger) {

	handler := &AuthGrpcHandler{trace: trace, logger: logger, store: store}

	// register the AuthServiceServer
	auth.RegisterAuthServiceServer(s, handler)

}

func (n *AuthGrpcHandler) GetUserById(ctx context.Context, payload *auth.GetUserByIdRequest) (*auth.GetUserByIdResponse, error) {

	_, span := n.trace.Start(ctx, "retrieving user")
	defer span.End()
	n.logger.WithContext(ctx).Info("retrieving user starts")
	user, err := n.store.Users().GetByEmailOrId(ctx, &store.User{ID: int(payload.UserId)})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	roles := make([]*auth.Role, len(user.Roles))
	for i, r := range user.Roles {
		roles[i] = &auth.Role{
			Id:       int32(r.ID),
			Name:     string(r.Name),
			IsActive: r.IsActive,
		}

	}
	return &auth.GetUserByIdResponse{
		User: &auth.User{
			Id:    int32(user.ID),
			Name:  user.Name,
			Email: user.Email,
			Roles: roles,
		},
	}, nil

}
