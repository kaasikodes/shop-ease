package main

import (
	"github.com/kaasikodes/shop-ease/services/notification-service/db"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/handler"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/model"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/providers"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/broker"
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
	store := repository.NewSqlPaymentRepo()

	// register payment provider
	providers.RegisterProvider(model.PaystackPaymentProvider, providers.NewPaystackGateway(env.GetString("PAYSTACK_API_KEY", ""), store))
	providers.RegisterProvider(model.FlutterPaymentProvider, providers.NewFlutterkGateway(env.GetString("FLUTTER_API_KEY", ""), store))

	shutdown := observability.InitTracer("payment-service")

	defer shutdown()

	tr := otel.Tracer("example.com/trace")
	logCfg := logger.LogConfig{
		LogFilePath:       "../../logs/payment-service.log",
		Format:            logger.DefaultLogFormat,
		PrimaryIdentifier: serviceIdentifier,
	}
	logger := logger.New(logCfg)
	// logger := logger.NewZapLogger(logCfg)
	cfg := config{
		addr: env.GetString("ADDR", ":3010"),

		env: env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "mysql://root:root123$@localhost"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection estatblished")

	metricsReg := prometheus.NewRegistry()
	metrics := NewMetrics(metricsReg)

	broker := broker.NewKafkaHelper([]string{":9092"}, events.PaymentTopic)
	defer broker.Close()
	var app = &application{
		config:  cfg,
		logger:  logger,
		metrics: metrics,
		trace:   tr,
		broker:  broker,

		paymentRegistry: providers.ProviderRegistry,
		store:           store,
	}
	mux := app.mount(metricsReg)

	logger.Fatal(app.run(mux))

	// grpc server

	// event handler
	eventHandler := handler.InitEventHandler(store)

	go func() {
		broker.Subscribe(events.SubscriptionTopic, eventHandler.HandleSubscriptionEvents)
		broker.Subscribe(events.OrderTopic, eventHandler.HandleSubscriptionEvents)

	}()

}
