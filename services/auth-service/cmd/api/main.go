package main

import (
	"log"

	"github.com/kaasikodes/shop-ease/cmd/logger"
	"github.com/kaasikodes/shop-ease/internal/env"
)

const version = "0.0.0"




func main() {
	logger := logger.New("../../app.log")
	// logger := logger.NewZapLogger()
	cfg := config{
		addr: env.GetString("ADDR", ":3010"),
		apiURL: env.GetString("API_URL", "localhost:9010"),
		frontendUrl: env.GetString("FRONTEND_URL", "localhost:3000"),
		env: env.GetString("ENV", "development"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "mysql://root:root123$@localhost"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redis: redisConfig{},
		mail: mailConfig{},
		auth: authConfig{},
		
		
	}
	
	var app = &application{
		config: cfg,
		rateLimiter: rateLimiterConfig{},
		logger: logger,
		
	}
	mux := app.mount();

	log.Fatal(app.run(mux))

}
