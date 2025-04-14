package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi"
	"github.com/kaasikodes/shop-ease/shared/env"
	"github.com/kaasikodes/shop-ease/shared/logger"
)

type GatewayService struct {
	Auth         struct{}
	Notification struct{}
	Order        struct{}
	// payment                struct{}
	// product                struct{}
	// reviewAndRating        struct{}
	// searchAndRecommend     struct{}
	// subscriptionAndTraffic struct{}
	// vendor                 struct{}
}

func NewAuthProxyHandler() http.Handler {
	target := &url.URL{Scheme: "http", Host: "localhost:3000"}
	proxy := httputil.NewSingleHostReverseProxy(target)
	return http.StripPrefix("/auth", proxy)
}
func NewNotificationProxyHandler() http.Handler {
	target := &url.URL{Scheme: "http", Host: "localhost:3020"}
	proxy := httputil.NewSingleHostReverseProxy(target)
	return http.StripPrefix("/notification", proxy)
}

func run() error {
	r := chi.NewRouter()
	addr := env.GetString("ADDR", ":3070")

	r.Handle("/auth/*", NewAuthProxyHandler())
	r.Handle("/notification/*", NewNotificationProxyHandler())

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	logger := logger.New("../../app.log")
	logger.Info("App running starting to run on .....", addr)
	err := server.ListenAndServe()

	if err != nil {
		return err
	}

	return nil

}
func main() {

	log.Fatal(run())

}
