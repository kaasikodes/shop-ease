package main

import (
	grpc_client "github.com/kaasikodes/shop-ease/services/auth-service/cmd/grpc"
	"github.com/kaasikodes/shop-ease/services/auth-service/internal/db"
	grpc_server "github.com/kaasikodes/shop-ease/services/auth-service/internal/grpc-server"
	"github.com/kaasikodes/shop-ease/services/auth-service/internal/oauth/provider"
	"github.com/kaasikodes/shop-ease/shared/env"
	jwttoken "github.com/kaasikodes/shop-ease/shared/jwt_token"

	store "github.com/kaasikodes/shop-ease/services/auth-service/internal/store/sql-store"
	"github.com/kaasikodes/shop-ease/shared/broker"
	"github.com/kaasikodes/shop-ease/shared/events"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/observability"
	"github.com/kaasikodes/shop-ease/shared/proto/notification"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
)

const version = "0.0.0"
const serviceIdentifier = "auth-service"

func main() {
	shutdown := observability.InitTracer("auth-service")
	// TODO: Create a auth topic, and also find out about the straming and the offset, refactor to compose kafka to the app as an queue system so you can switch between rabbitmq to kafka, then move to the consumers like notification
	// Define the events that each service should emit - say for example auth.user.verrifed, auth.customer_created

	defer shutdown()

	tr := otel.Tracer("example.com/trace")
	logCfg := logger.LogConfig{
		LogFilePath:       "../../logs/auth-service.log",
		Format:            logger.DefaultLogFormat,
		PrimaryIdentifier: serviceIdentifier,
	}
	logger := logger.New(logCfg)
	// logger := logger.NewZapLogger(logCfg)
	cfg := config{
		grpcAddr:    env.GetString("GRPC_ADDR", ":4010"),
		addr:        env.GetString("ADDR", ":3010"),
		apiURL:      env.GetString("API_URL", "localhost:9010"),
		frontendUrl: env.GetString("FRONTEND_URL", "localhost:3000"),
		env:         env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", ""),
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
	broker := broker.NewKafkaHelper([]string{":9092"}, events.AuthTopic)
	defer broker.Close()
	// set up oauth providers
	githubOauthProvider := provider.NewGithubOauthProvider(env.GetString("GITHUB_CLIENT_ID", ""), env.GetString("GITHUB_CLIENT_SECRET", ""), env.GetString("GITHUB_REDIRECT_URL", ""))
	provider.OauthProviderRegistry = make(map[provider.OauthProviderType]provider.OauthProvider)
	provider.OauthProviderRegistry[provider.OauthProviderTypeGithub] = githubOauthProvider
	// set up jwt
	jwt := jwttoken.NewJwtMaker(env.GetString("JWT_SECRET", ""))
	var app = &application{
		config:                cfg,
		rateLimiter:           rateLimiterConfig{},
		logger:                logger,
		store:                 store.NewSQLStorage(db),
		notificationService:   n,
		metrics:               metrics,
		trace:                 tr,
		broker:                broker,
		oauthProviderRegistry: provider.OauthProviderRegistry,
		jwt:                   jwt,
	}
	mux := app.mount(metricsReg)

	// grpc server
	go func() {
		authGrpcServer := grpc_server.NewAuthGRPCServer(cfg.grpcAddr, grpc_server.Config{
			Db: grpc_server.DbConfig{
				Addr:         cfg.db.addr,
				MaxOpenConns: cfg.db.maxOpenConns,
				MaxIdleConns: cfg.db.maxIdleConns,
				MaxIdleTime:  cfg.db.maxIdleTime,
			},
		}, logger)
		logger.Fatal(authGrpcServer.Run()) //has a graceful shutdown built in, consider revisting ...

	}()

	// http server
	logger.Fatal(app.run(mux))

}
