package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

type HttpServer struct {
	Port int32
}

func metricsRouter() http.Handler {
	promHandler := func(next http.Handler) http.Handler { return promhttp.Handler() }
	emptyHandler := func(w http.ResponseWriter, r *http.Request) {}
	r := chi.NewRouter()
	r.Use(promHandler)
	r.Get("/", emptyHandler)
	return r
}

func (server *HttpServer) Start() {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Mount("/metrics", metricsRouter())

	log.Error(http.ListenAndServe(fmt.Sprintf(":%v", server.Port), r))
}
