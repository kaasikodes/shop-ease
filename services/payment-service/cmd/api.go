package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/model"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/providers"
	"github.com/kaasikodes/shop-ease/services/payment-service/internal/repository"
	"github.com/kaasikodes/shop-ease/shared/broker"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	addr     string
	grpcAddr string
	db       dbConfig
	env      string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type application struct {
	config config
	// observability & monitoring
	metrics *metrics
	trace   trace.Tracer
	logger  logger.Logger
	// message broker
	broker          broker.MessageBroker
	store           repository.PaymentRepo
	paymentRegistry map[model.PaymentProvider]providers.PaymentGateway // done this way, so if certain types of payments are to be made with a certain provider we can flexibly implement this, lets say vendor payment to flutter and order payment to paystack based on customer requirements

}

func (app *application) mount(reg *prometheus.Registry) http.Handler {
	r := chi.NewRouter()
	// Add the metrics middleware
	r.Use(app.metricsMiddleware)
	// middleware to get the vendor id from headers, and that they are only accessing and modifying the data they own

	r.Get("/healthz", app.healthzHandler)
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP(w, r)

	})
	//TODO: Add middleware for webhook, so you ensure that only certain ip addresses that belong to a provider can call it
	r.Get("/webhook", app.webHookHandler) //will be called by providers and in here will check and update transaction record, and then send event
	r.Route("/v1", func(r chi.Router) {

		r.Route("/transactions", func(r chi.Router) {
			r.Get("/:transactionId", app.getTransactionByIdHandler)
			r.Get("/", app.getTransactionsHandler)

		})

	})

	return r

}

func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("Api running starting to run on .....", app.config.addr)

	err := server.ListenAndServe()

	go func() {
		app.logger.Info("Grpc server running in the background on .....", app.config.addr)
		grpcServer := NewPaymentGRPCServer(app.config.grpcAddr, app.config, app.logger)
		grpcServer.Run() //has a graceful shutdown built in, consider revisting ...

	}()

	if err != nil {
		return err
	}

	return nil

}
