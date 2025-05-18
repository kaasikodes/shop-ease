package main

import (
	"fmt"
	"net"

	"github.com/kaasikodes/shop-ease/services/order-service/internal/handler"
	"github.com/kaasikodes/shop-ease/services/order-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/database"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/observability"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr   string
	config config
	logger logger.Logger
}

func NewProductGRPCServer(addr string, config config, logger logger.Logger) *gRPCServer {
	logger.Info("addr for product grpc server", addr)
	return &gRPCServer{addr, config, logger}

}

func (s *gRPCServer) Run() error {

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	db, err := database.NewMySqlDB(s.config.db.addr, s.config.db.maxOpenConns, s.config.db.maxOpenConns, s.config.db.maxIdleTime)
	if err != nil {
		return err
	}
	defer db.Close()

	// tracer
	store := repository.NewPostgresOrderRepo(db)
	shutdown := observability.InitTracer("notification-service")
	defer shutdown()

	trace := otel.Tracer("app.notification/trace")

	handler.NewOrderGrpcHandler(grpcServer, store, trace, s.logger)
	s.logger.Info("The GRPC SERVER IS UP >>>>>>")

	return grpcServer.Serve(lis)

}
