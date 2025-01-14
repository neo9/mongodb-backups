package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/neo9/mongodb-backups/pkg/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
	err := http.ListenAndServe(fmt.Sprintf(":%v", server.Port), r)
	if err != nil {
		log.Error("Could not start server in port %s", server.Port)
	}
}
