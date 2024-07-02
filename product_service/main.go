package main

import (
	"os"
	"os/signal"
	"product_service/api"
	"product_service/config"
	"product_service/database"
	"product_service/logger"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	logger.SetupZeroLog()
	log.Info().Msg("[APP] initialize app")

	cfg := config.LoadConfig(".")
	log.Info().Msgf("[APP] current env: \t%s", cfg.App.Env)

	database.Setup(&cfg.Database)
	// database.EnsureMigrations(database.Migrations)

	api.ServePublicServer(cfg.Server)
	api.ServeAPIDocs(cfg.Server)

	gracefulShutdown(
		func() error {
			return database.DBConnection.Close()
		},
		func() error {
			log.Warn().Msg("[APP] shutting down service")
			os.Exit(0)
			return nil
		},
	)
}

func gracefulShutdown(ops ...func() error) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	if <-shutdown != nil {
		for _, op := range ops {
			if err := op(); err != nil {
				log.Panic().AnErr("gracefullShutdown op failed", err)
			}
		}
	}
}
