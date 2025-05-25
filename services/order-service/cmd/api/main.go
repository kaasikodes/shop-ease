package main

import (
	"time"

	"github.com/kaasikodes/shop-ease/services/order-service/internal/cache"
	"github.com/kaasikodes/shop-ease/services/order-service/internal/handler"
	"github.com/kaasikodes/shop-ease/services/order-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/broker"
	"github.com/kaasikodes/shop-ease/shared/database"
	"github.com/kaasikodes/shop-ease/shared/env"
	"github.com/kaasikodes/shop-ease/shared/events"
	jwttoken "github.com/kaasikodes/shop-ease/shared/jwt_token"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/observability"
	"github.com/kaasikodes/shop-ease/shared/proto/auth"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
)

var version = "0.0.0"
var serviceIdentifier = "order_service"

func main() {

	shutdown := observability.InitTracer("order-service")

	defer shutdown()

	tr := otel.Tracer("example.com/trace")
	logCfg := logger.LogConfig{
		LogFilePath:       "../../logs/order-service.log",
		Format:            logger.DefaultLogFormat,
		PrimaryIdentifier: serviceIdentifier,
	}
	logger := logger.New(logCfg)
	// logger := logger.NewZapLogger(logCfg)
	cfg := config{
		addr: env.GetString("ADDR", ":3010"),

		env: env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""), //change to postgres db
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	db, err := database.NewPostgresSqlDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection estatblished")
	store := repository.NewPostgresOrderRepo(db)

	metricsReg := prometheus.NewRegistry()
	metrics := NewMetrics(metricsReg)

	broker := broker.NewKafkaHelper([]string{env.GetString("KAFKA_BROKER_ADDR", ":9092")}, events.ProductTopic)
	defer broker.Close()

	// grpc clients
	authConn := NewGRPCClient(env.GetString("AUTH_GRPC_SERVER_ADDR", ":4040"), logger)
	defer authConn.Close()
	authClient := auth.NewAuthServiceClient(authConn)

	// set up jwt
	jwt := jwttoken.NewJwtMaker(env.GetString("JWT_SECRET", ""))

	// cache
	inMemoryCache := cache.NewInMemoryCache(time.Duration(time.Hour*24*1), time.Duration(time.Hour*24*3))
	redisCache := cache.NewRedisCache(env.GetString("REDIS_ADDR", ""), env.GetString("REDIS_PWD", ""), env.GetInt("REDIS_LOGICAL_DB", 1), serviceIdentifier, time.Duration(time.Hour*24*1))
	var app = &application{
		config:  cfg,
		logger:  logger,
		metrics: metrics,
		trace:   tr,
		broker:  broker,
		store:   store,
		jwt:     jwt,
		clients: Clients{
			auth: authClient,
		},
		cache: Cache{
			memory: inMemoryCache,
			redis:  redisCache,
		},
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
