package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}


type HttpServer struct {
	Port int32
}

func (server *HttpServer) Start() {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// r.Mount("/metrics", metricsRouter())

	log.Error(http.ListenAndServe(fmt.Sprintf(":%v", server.Port), r))
}

