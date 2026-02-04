package app

import (
	"github.com/Rasulikus/url-shortener/internal/config"
	"github.com/Rasulikus/url-shortener/internal/repository"
	"github.com/Rasulikus/url-shortener/internal/repository/memory"
	"github.com/Rasulikus/url-shortener/internal/repository/postgres"
	"github.com/Rasulikus/url-shortener/internal/repository/postgres/url"
	urlService "github.com/Rasulikus/url-shortener/internal/service/url"
	"github.com/Rasulikus/url-shortener/internal/transport/http/handler"
	"github.com/Rasulikus/url-shortener/internal/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func App(cfg *config.Config) *gin.Engine {
	err := logger.Init(logger.Config{
		Level: cfg.LogLevel,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize logger")
	}

	var urlRepo repository.URLRepository

	switch cfg.Storage {
	case "postgresql":
		pool, err := postgres.NewPool(cfg.DB)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize postgres pool")
		}
		urlRepo, err = url.NewRepository(pool)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize postgres repository")
		}
	case "memory":
		m := memory.New()
		urlRepo, err = memory.NewRepository(m)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize memory repository")
		}
	}

	urlServ, err := urlService.NewService(cfg.BaseURL, urlRepo)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize url service")
	}
	urlHandler := handler.NewURLHandler(urlServ)

	r := gin.Default()
	urlApi := r.Group("/urls")
	{
		urlApi.POST("", urlHandler.Create)
		urlApi.GET("/:alias", urlHandler.GetByAlias)
	}

	return r
}
