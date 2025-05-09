package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	grpc_server "github.com/kaasikodes/shop-ease/services/subscription-and-traffic-service/internal/grpc"
	"github.com/kaasikodes/shop-ease/services/subscription-and-traffic-service/internal/traffic"
	vendorplan "github.com/kaasikodes/shop-ease/services/subscription-and-traffic-service/internal/vendor-plan"
	"github.com/kaasikodes/shop-ease/shared/broker"
	"github.com/kaasikodes/shop-ease/shared/events"
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
	broker broker.MessageBroker
	// store
	store struct {
		plan vendorplan.VendorPlanRepo
	}
}

func (app *application) mount(reg *prometheus.Registry) http.Handler {
	log.Println("Api mounted ....")
	r := chi.NewRouter()
	// Add the metrics middleware
	r.Use(app.metricsMiddleware)
	// middleware to get the vendor id from headers, and that they are only accessing and modifying the data they own

	r.Get("/healthz", app.healthzHandler)
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP(w, r)

	})
	r.Route("/v1", func(r chi.Router) {

		// create vendor plans
		// retrieve vendor plans
		// activate or deactivate vendor plans in bulk
		// retrieve user interactions and filter based on vendor

		r.Route("/plan", func(r chi.Router) {

			r.Post("/", app.createVendorPlan)
			r.Patch("/toggle-activation", app.toggleVendorPlanActivation)
			r.Get("/", app.getVendorPlans)
		})
		r.Route("/interaction", func(r chi.Router) {
			r.Get("/user", app.getUserInteractions) //filter - vendorId

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
		server := grpc_server.NewSubscriptionGRPCServer(app.config.grpcAddr, app.logger)
		server.Run() //has a graceful shutdown built in, consider revisting ...

	}()
	// run event handler in background
	// handler to listen to the following events - user interactions: order made for product(order service); item added to wishlist(search & recommend - can change ), payment made for subscription
	// emits the following events
	eventHandler := traffic.InitEventHandler(app.store.plan)
	broker := broker.NewKafkaHelper([]string{":9092"}, events.SubscriptionTopic)
	defer broker.Close()
	go func() {
		broker.Subscribe(events.AuthTopic, eventHandler.HandleAuthEvents)

	}()

	if err != nil {
		return err
	}

	return nil

}
