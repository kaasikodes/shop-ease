package main

import (
	"github.com/kaasikodes/shop-ease/services/vendor-service/internal/products"
	"github.com/kaasikodes/shop-ease/shared/broker"
	"github.com/kaasikodes/shop-ease/shared/env"
	"github.com/kaasikodes/shop-ease/shared/events"
	jwttoken "github.com/kaasikodes/shop-ease/shared/jwt_token"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
)

const serviceIdentifier = "vendor-service"

// store = name, address, contact details - phone, email,vendor_id
// sales = product_id, paid_at, cutomer_order_id, store_id, - Remove
// inventory = product_id, amount, addedAt
// customer_orders - order_id, paidAt, deliveredAt, store_id

// So vendor service (needs primary actor vendor as table that has stores and all) and then the product service

// plan db and create api's update postman
// then move to products service
// then to payment service
// and then back to the auth register vendor service, go to subscription service first
// and then to order service

const version = "0.0.0"

func main() {
	// background
	productStore := products.NewInMemoryProductRepo()
	productHandler := products.InitProductHandler(productStore)
	broker := broker.NewKafkaHelper([]string{":9092"}, events.VendorTopic)
	defer broker.Close()
	go func() {
		broker.Subscribe(events.ProductTopic, productHandler.HandleProductEvents)
		broker.Subscribe(events.AuthTopic, productHandler.HandleAuthEvents)

	}()

	// main app - api

	tr := otel.Tracer("example.com/trace")
	// logger
	logCfg := logger.LogConfig{
		LogFilePath:       "../../logs/vendor-service.log",
		Format:            logger.DefaultLogFormat,
		PrimaryIdentifier: serviceIdentifier,
	}
	logger := logger.New(logCfg)
	//  jwt
	jwt := jwttoken.NewJwtMaker(env.GetString("JWT_SECRET", ""))
	app := &application{
		jwt:    jwt,
		logger: logger,
		trace:  tr,
	}

	// metrics
	metricsReg := prometheus.NewRegistry()

	mux := app.mount(metricsReg)

	logger.Fatal(app.run(mux))

}
