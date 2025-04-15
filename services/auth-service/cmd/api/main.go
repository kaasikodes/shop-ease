package main

import (
	grpc_client "github.com/kaasikodes/shop-ease/services/auth-service/cmd/grpc"
	"github.com/kaasikodes/shop-ease/services/auth-service/internal/db"
	"github.com/kaasikodes/shop-ease/services/auth-service/internal/env"
	store "github.com/kaasikodes/shop-ease/services/auth-service/internal/store/sql-store"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/proto/notification"
	"github.com/prometheus/client_golang/prometheus"
)

const version = "0.0.0"

func main() {
	logger := logger.New("../../app.log")
	// logger := logger.NewZapLogger()
	cfg := config{
		addr:        env.GetString("ADDR", ":3010"),
		apiURL:      env.GetString("API_URL", "localhost:9010"),
		frontendUrl: env.GetString("FRONTEND_URL", "localhost:3000"),
		env:         env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "mysql://root:root123$@localhost"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redis: redisConfig{},
		mail:  mailConfig{},
		auth:  authConfig{},
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection estatblished")
	// grpcConn
	notificationConn := grpc_client.NewGRPCClient(":5050")
	defer notificationConn.Close()
	n := notification.NewNotificationServiceClient(notificationConn)
	metricsReg := prometheus.NewRegistry()
	metrics := NewMetrics(metricsReg)
	var app = &application{
		config:              cfg,
		rateLimiter:         rateLimiterConfig{},
		logger:              logger,
		store:               store.NewSQLStorage(db),
		notificationService: n,
		metrics:             metrics,
	}
	mux := app.mount(metricsReg)

	logger.Fatal(app.run(mux))

}
