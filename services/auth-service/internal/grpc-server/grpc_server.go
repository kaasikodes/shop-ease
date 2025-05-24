package grpc_server

import (
	"fmt"
	"net"

	store "github.com/kaasikodes/shop-ease/services/auth-service/internal/store/sql-store"

	"github.com/kaasikodes/shop-ease/shared/database"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/observability"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type gRPCServer struct {
	addr   string
	config Config
	logger logger.Logger
}

func NewAuthGRPCServer(addr string, config Config, logger logger.Logger) *gRPCServer {
	logger.Info("addr for payment grpc server", addr)
	return &gRPCServer{addr, config, logger}

}

func (s *gRPCServer) Run() error {

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	db, err := database.NewMySqlDB(s.config.Db.Addr, s.config.Db.MaxOpenConns, s.config.Db.MaxIdleConns, s.config.Db.MaxIdleTime)
	if err != nil {
		return err
	}
	defer db.Close()

	// tracer
	store := store.NewSQLStorage(db)
	shutdown := observability.InitTracer("notification-service")
	defer shutdown()

	trace := otel.Tracer("app.notification/trace")

	NewAuthGRPCHandler(grpcServer, store, trace, s.logger)
	s.logger.Info("The GRPC SERVER IS UP >>>>>>")

	return grpcServer.Serve(lis)

}
