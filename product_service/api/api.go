package api

import (
	_ "embed"
	"fmt"
	"net/http"
	"product_service/config"

	"github.com/flowchartsman/swaggerui"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func ServePublicServer(cfg config.ServerConfig) {
	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		w.Write([]byte("service is up!"))
	})

	httpPort := fmt.Sprintf(":%d", cfg.Port)
	go func() {
		if err := http.ListenAndServe(httpPort, r); err != nil {
			log.Panic().AnErr("ServePublicServer http.ListenAndServe failed", err)
		}
	}()

	log.Info().Msgf("[HTTP] server is running at port: \t%d\n", cfg.Port)
}

var apiDocs []byte

func ServeAPIDocs(cfg config.ServerConfig) {
	mux := http.NewServeMux()
	mux.Handle("/api-docs/", http.StripPrefix("/api-docs", swaggerui.Handler(apiDocs)))

	httpPort := fmt.Sprintf(":%d", cfg.APIDocsPort)
	go func() {
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			log.Panic().AnErr("ServeAPIDocs http.ListenAndServe failed", err)
		}
	}()

	log.Info().Msgf("[HTTP] api docs server is running at port: \t%d\n", cfg.APIDocsPort)
}
