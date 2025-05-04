package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	grpc_server "github.com/kaasikodes/shop-ease/services/vendor-service/internal/grpc"
	"github.com/kaasikodes/shop-ease/services/vendor-service/internal/orders"
	"github.com/kaasikodes/shop-ease/services/vendor-service/internal/store"
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
	broker broker.MessageBroker
	// store
	store struct {
		store  store.StoreRepo
		orders orders.OrderRepo
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
		r.Route("/analytics", func(r chi.Router) {

			r.Get("/sales", app.salesAnalyticsHandler)         //sales report will have month filters
			r.Get("/inventory", app.inventoryAnalyticsHandler) //inventory report what items were stocked, returned ordered in a single order
		})
		r.Route("/returns", func(r chi.Router) {
			r.Get("/", app.getReturnsHandler) //get a list of orders that the vendor is part of
			r.Get("/:id", app.getReturnsByIdHandler)
			r.Post("/:id/approve", app.acceptReturnHandler)
			r.Post("/:id/reject", app.rejectReturnHandler)
			r.Post("/:id/receive", app.rejectReturnHandler)

		})
		r.Route("/orders", func(r chi.Router) {
			r.Get("/", app.getOrdersHandler)       //get a list of orders that the vendor is part of
			r.Get("/:id", app.getOrderByIdHandler) //get a list of items in an order (can filter by order they belong to ) - shows the status(pending,accepted/rejected/processed, shipped, delivered, fulfilled, returned_by_user) of each item
			r.Post("/:id/accept", app.acceptOrderHandler)
			r.Post("/:id/reject", app.rejectOrderHandler)
			r.Post("/:id/ship", app.shipOrderHandler)
			r.Post("/:id/reject/bulk", app.bulkRejectOrderHandler)
			r.Post("/:id/accept/bulk", app.bulkAcceptOrderHandler)
			r.Post("/:id/reject/bulk", app.bulkRejectOrderHandler)
			r.Post("/:id/ship/bulk", app.bulkShipOrderHandler)

		})
		r.Route("/seller", func(r chi.Router) {
			r.Post("/", app.createSellerHandler)
			r.Get("/:sellerId", app.GetSellerHandler)

		})
		r.Route("/store", func(r chi.Router) {
			r.Post("/", app.createStoreHandler)
			r.Get("/{storeId}/performance", app.getStorePerformanceHandler) //get performance score of the store: lets grade by the sold stock in a month/total stock in a month
			r.Patch("/{storeId}", app.updateStoreHandler)
			r.Get("/{storeId}", app.getStoreHandler)
			r.Get("/{storeId}/products", app.getProductsHandler)                       //products in the store
			r.Get("/{storeId}/inventory", app.getInventoriesHandler)                   //inventories in the store
			r.Post("/{storeId}/inventory/bulk", app.bulkAddInventoryHandler)           // add a multitude of product inventories
			r.Post("/{storeId}/inventory/", app.addInventoryHandler)                   //add inventory for a single product
			r.Put("/{storeId}/inventory/{inventoryId}", app.updateInventoryHandler)    //update inventory for a single product, audit
			r.Delete("/{storeId}/inventory/{inventoryId}", app.deleteInventoryHandler) //delete inventory for a single product, only if not been used - update audit

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
		server := grpc_server.NewVendorGRPCServer(app.config.grpcAddr, app.logger)
		server.Run() //has a graceful shutdown built in, consider revisting ...

	}()

	if err != nil {
		return err
	}

	return nil

}
