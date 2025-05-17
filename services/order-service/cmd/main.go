package main

import (
	"github.com/kaasikodes/shop-ease/services/order-service/internal/handler"
	"github.com/kaasikodes/shop-ease/services/order-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/broker"
	"github.com/kaasikodes/shop-ease/shared/database"
	"github.com/kaasikodes/shop-ease/shared/env"
	"github.com/kaasikodes/shop-ease/shared/events"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/observability"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
)

var version = "0.0.0"
var serviceIdentifier = "payment_service"

func main() {

	shutdown := observability.InitTracer("product-service")

	defer shutdown()

	tr := otel.Tracer("example.com/trace")
	logCfg := logger.LogConfig{
		LogFilePath:       "../../logs/product-service.log",
		Format:            logger.DefaultLogFormat,
		PrimaryIdentifier: serviceIdentifier,
	}
	logger := logger.New(logCfg)
	// logger := logger.NewZapLogger(logCfg)
	cfg := config{
		addr: env.GetString("ADDR", ":3010"),

		env: env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "mysql://root:root123$@localhost"), //change to postgres db
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	db, err := database.NewSqlDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection estatblished")
	store := repository.NewPostgresOrderRepo(db)

	metricsReg := prometheus.NewRegistry()
	metrics := NewMetrics(metricsReg)

	broker := broker.NewKafkaHelper([]string{":9092"}, events.ProductTopic)
	defer broker.Close()
	var app = &application{
		config:  cfg,
		logger:  logger,
		metrics: metrics,
		trace:   tr,
		broker:  broker,

		store: store,
	}
	mux := app.mount(metricsReg)

	logger.Fatal(app.run(mux))

	// grpc server

	// event handler
	eventHandler := handler.InitEventHandler(store)

	go func() {
		broker.Subscribe(events.VendorTopic, eventHandler.HandleVendorEvents)

	}()

}
