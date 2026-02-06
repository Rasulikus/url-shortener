package main

import (
	"fmt"
	"net/http"

	"github.com/Rasulikus/url-shortener/internal/app"
	"github.com/Rasulikus/url-shortener/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create config")
	}

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler: app.App(cfg),
	}

	log.Info().Msgf("starting server on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
