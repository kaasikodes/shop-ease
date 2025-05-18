package main

import (
	"fmt"
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

type serviceDetail struct {
	scheme string
	host   string
	prefix string
}

func NewProxyHandler(scheme string, host string, prefix string) http.Handler {
	target := &url.URL{Scheme: scheme, Host: host}
	proxy := httputil.NewSingleHostReverseProxy(target)
	return http.StripPrefix(prefix, proxy)
}

var services []serviceDetail = []serviceDetail{
	{"http", "localhost:3010", "/auth"}, //TODO: ALternatively read from the .env and use default values where needed, still not guarateed to have api gateway service be aware of change, cosider your own implemetation on way could be service server notifying the api gateway but how ... lets say in the memory of the api gateway server, what are the downsides. So when up send request and when off send request is it possible ...
	{"http", "localhost:3020", "/notification"},
	{"http", "localhost:3030", "/order"},
	{"http", "localhost:3040", "/payment"},
	{"http", "localhost:3050", "/product"},
	{"http", "localhost:3060", "/subscription"},
	{"http", "localhost:3070", "/vendor"},
}

func run() error {
	r := chi.NewRouter()
	addr := env.GetString("ADDR", ":3000")

	for _, s := range services {
		// TODO: add auth middleware for all services asides from auth, and push in the details of the user to service as headers or in acontext and test
		r.Handle(fmt.Sprintf("%s/*", s.prefix), NewProxyHandler(s.scheme, s.host, s.prefix))

	}

	r.Handle("/notification/*", NewNotificationProxyHandler())
	// TODO: Implement the following middelwares
	// auth, logging, tracing and then prometheus and grafan (for each service ?) -> then rabbit mq
	// Then refactor on another branch to use proper setup for micro services & why its preferred to this, should also have service discovery

	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	logCfg := logger.LogConfig{
		LogFilePath: "../../logs/api-gateway.log",
		Format:      "",
	}
	logger := logger.New(logCfg)
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
