package rest

import (
	"net/http"
	"time"

	"github.com/andrelince/github-proxy/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewRest(
	router *mux.Router,
	r Handler,
	conf config.Config,
) *http.Server {

	router.
		Path("/health").
		Methods(http.MethodGet).
		HandlerFunc(r.Health)

	corsOpts := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}

	return &http.Server{
		Handler:      cors.New(corsOpts).Handler(router),
		Addr:         ":" + conf.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}
