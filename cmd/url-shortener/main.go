package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Rasulikus/url-shortener/internal/app"
	"github.com/Rasulikus/url-shortener/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler: app.App(cfg),
	}
	log.Printf("Starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
