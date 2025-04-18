package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	"github.com/kaasikodes/shop-ease/shared/logger"
	"github.com/kaasikodes/shop-ease/shared/proto/notification"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendUrl string
	auth        authConfig
	redis       redisConfig
}

type rateLimiterConfig struct {
}
type redisConfig struct {
}
type authConfig struct {
}
type mailConfig struct {
}
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type application struct {
	config              config
	rateLimiter         rateLimiterConfig
	logger              logger.Logger
	store               store.Storage
	notificationService notification.NotificationServiceClient
	metrics             *metrics
}

func (app *application) mount(reg *prometheus.Registry) http.Handler {
	log.Println("Api mounted ....")
	r := chi.NewRouter()
	// Add the metrics middleware
	r.Use(app.metricsMiddleware)

	r.Get("/healthz", app.healthzHandler)
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}).ServeHTTP(w, r)
		// promhttp.Handler().ServeHTTP(w, r)
	})
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.registerHandler) // customer(happy path), vendor
			r.Post("/login", app.loginHandler)
			r.Post("/verify", app.verifyHandler)
			r.Post("/forgot-password", app.forgotPasswordHandler)
			r.Post("/reset-password", app.resetPasswordHandler)
			// TODO: Add middleware here to check wether user is authenticated
			r.Group(func(r chi.Router) {
				r.Get("/me", app.retriveAuthAccountHandler)
			})
		})

		// They ought to be grpc endpoints that other services can call like
		// isUserVerified : vendor
		// listOfUsers with pagination params and all :
		// getUserById : vendor
		// doesUserHaveRoleAsActive
		// activateOrDeactivateUserRole: payment, admin
		// Note on login in: user ought to provide the code sent to mail, 2-FA authentication

		// Like wise
		// The login for vendor will have to check if the user has an active subscription or not, if not don't allow login
		//

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

	app.logger.Info("App running starting to run on .....", app.config.addr)

	err := server.ListenAndServe()
	log.Println("App running stopping to run on .....", app.config.addr)

	if err != nil {
		return err
	}

	return nil

}
