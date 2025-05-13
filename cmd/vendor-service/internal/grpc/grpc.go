package grpc_server

import (
	"fmt"
	"net"

	"github.com/kaasikodes/shop-ease/services/notification-service/config"
	"github.com/kaasikodes/shop-ease/services/vendor-service/internal/seller"
	"github.com/kaasikodes/shop-ease/services/vendor-service/pkg/database"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/observability"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr   string
	logger logger.Logger
}

func NewVendorGRPCServer(addr string, logger logger.Logger) *gRPCServer {
	logger.Info("Initializing Grpc Server .....")
	return &gRPCServer{addr, logger}

}

func (s *gRPCServer) Run() error {

	cfg := config.ServiceConfig

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	defer grpcServer.GracefulStop()
	db, err := database.New(cfg.Db.Addr, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleConns, cfg.Db.MaxIdleTime)
	if err != nil {
		return err
	}
	defer db.Close()

	// tracer
	shutdown := observability.InitTracer("vendor-service")
	defer shutdown()

	trace := otel.Tracer("app.notification/trace")
	store := Store{
		seller: seller.NewSqlSellerRepo(db),
	}
	NewGRPCHandler(grpcServer, store, trace, s.logger)
	s.logger.Info("The Vendor GRPC SERVER IS UP .....")

	return grpcServer.Serve(lis)

}
