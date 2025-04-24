package grpc_server

import (
	"fmt"
	"log"
	"net"

	"github.com/kaasikodes/shop-ease/services/notification-service/config"
	"github.com/kaasikodes/shop-ease/services/notification-service/db"
	"github.com/kaasikodes/shop-ease/services/notification-service/service"
	store "github.com/kaasikodes/shop-ease/services/notification-service/store/sql-store"
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

func NewNotificiationGRPCServer(addr string, logger logger.Logger) *gRPCServer {
	log.Println("Isssssss", addr)
	return &gRPCServer{addr, logger}

}

func (s *gRPCServer) Run() error {

	cfg := config.ServiceConfig

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	db, err := db.New(cfg.Db.Addr, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleConns, cfg.Db.MaxIdleTime)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStore := store.NewSQLStorage(db)
	notificationServices := make([]service.NotificationService, 2)
	notificationStore := sqlStore.Notification()
	inAppService := service.NewInAppNotificationService(notificationStore, s.logger)
	emailService := service.NewEmailNotificationService(cfg.Mail, s.logger)
	notificationServices[0] = inAppService
	notificationServices[1] = emailService

	// tracer
	shutdown := observability.InitTracer("notification-service")
	defer shutdown()

	trace := otel.Tracer("app.notification/trace")

	NewNotificiationGRPCHandler(grpcServer, notificationServices, trace)
	s.logger.Info("The GRPC SERVER IS UP >>>>>>")

	return grpcServer.Serve(lis)

}
