package main

import (
	"log"
	"net"

	"github.com/kaasikodes/shop-ease/services/notification-service/db"
	"github.com/kaasikodes/shop-ease/services/notification-service/service"
	store "github.com/kaasikodes/shop-ease/services/notification-service/store/sql-store"
	"github.com/kaasikodes/shop-ease/shared/env"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
}

func NewNotificiationGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{addr}

}

func (s *gRPCServer) Run() error {

	cfg := config{

		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "mysql://root:root123$@localhost"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen to %v", err)
	}

	grpcServer := grpc.NewServer()
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStore := store.NewSQLStorage(db)
	notificationServices := []service.NotificationService{}
	notificationStore := sqlStore.Notification()
	inAppService := service.NewInAppNotificationService(notificationStore)
	emailService := service.NewEmailNotificationService(service.MailConfig{
		Addr: "",
	})
	notificationServices = append(notificationServices, inAppService)
	notificationServices = append(notificationServices, emailService)

	NewNotificiationGRPCHandler(grpcServer, notificationServices)

	return grpcServer.Serve(lis)

}

func main() {
	s := NewNotificiationGRPCServer(":5050")
	err := s.Run()
	if err != nil {
		log.Fatalf("Error setting up grpc server: %v", err)
	}
	log.Fatalln("Grpc Notification Server is running .....")
}
