package main

import (
	"sync"

	"github.com/kaasikodes/shop-ease/services/notification-service/config"
	"github.com/kaasikodes/shop-ease/services/notification-service/db"
	grpc_server "github.com/kaasikodes/shop-ease/services/notification-service/grpc"
	store "github.com/kaasikodes/shop-ease/services/notification-service/store/sql-store"
	"github.com/kaasikodes/shop-ease/shared/logger"
)

const version = "0.0.0"

func main() {
	logger := logger.New("../../app.log")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := grpc_server.NewNotificiationGRPCServer(config.ServiceConfig.GrpcAddr, logger)
		err := s.Run()
		if err != nil {
			logger.Info("error setting up grpc server: %v", err)
		}
		logger.Info("Grpc Notification Server is running .....End", config.ServiceConfig.GrpcAddr)

	}()
	wg.Wait()
	// logger := logger.NewZapLogger()
	cfg := config.ServiceConfig
	db, err := db.New(cfg.Db.Addr, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleConns, cfg.Db.MaxIdleTime)

	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection estatblished")

	sqlStore := store.NewSQLStorage(db)

	var app = &application{
		config: cfg,
		logger: logger,
		store:  sqlStore,
	}
	mux := app.mount()

	logger.Fatal(app.run(mux))

}
