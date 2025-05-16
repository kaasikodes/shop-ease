package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/kaasikodes/shop-ease/services/product-service/internal/repository"
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
	store  repository.ProductRepo
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

	r.Route("/v1", func(r chi.Router) {
		r.Route("/category", func(r chi.Router) {
			r.Get("/", app.getCategoriesHandler)
			r.Post("/bulk", app.bulkAddCategoriesHandler)
			r.Delete("/:categoryId", app.deleteCategoryHandler)
			r.Patch("/:categoryId", app.updateCategoryHandler)

		})
		r.Route("/discount", func(r chi.Router) {
			r.Post("/", app.createDiscountHandler)
			r.Get("/", app.getDiscountsHandler)
			r.Patch("/:discountId", app.updateDiscountApplicabilityHandler)
			r.Patch("/:discountId", app.updateDiscountExpiryDateHandler)

		})
		r.Route("/product-policy", func(r chi.Router) {
			r.Post("/", app.saveProductPolicyHandler)
			r.Get("/", app.getProductPolicyHandler)

		})

		r.Route("/products", func(r chi.Router) {
			r.Get("/", app.getProductsHandler)
			r.Post("/bulk", app.bulkAddProductsHandler)
			r.Delete("/:productId", app.deleteProductHandler)
			r.Delete("/:productId", app.deleteProductHandler)
			r.Patch("/:inventoryId", app.updateProductInventoryHandler)

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
		server := NewProductGRPCServer(app.config.grpcAddr, app.config, app.logger)
		server.Run() //has a graceful shutdown built in, consider revisting ...

	}()

	if err != nil {
		return err
	}

	return nil

}
