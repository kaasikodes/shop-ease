package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kaasikodes/shop-ease/services/notification-service/config"
	"github.com/kaasikodes/shop-ease/services/notification-service/store"
	"github.com/kaasikodes/shop-ease/shared/logger"
)

type application struct {
	config config.Config
	logger logger.Logger
	store  store.Storage
}

func (app *application) mount() http.Handler {
	log.Println("Api mounted ....")
	r := chi.NewRouter()

	r.Get("/healthz", app.healthzHandler)
	r.Route("/v1", func(r chi.Router) {
		r.Route("/notification", func(r chi.Router) {
			r.Post("/", app.getAllNotifications)

		})

	})

	return r

}
func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr:         app.config.ApiAddr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("App running starting to run on .....", app.config.ApiAddr)
	err := server.ListenAndServe()
	log.Println("App running stopping to run on .....", app.config.ApiAddr)

	if err != nil {
		return err
	}

	return nil

}
