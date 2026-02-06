package app

import (
	"github.com/Rasulikus/url-shortener/internal/config"
	"github.com/Rasulikus/url-shortener/internal/repository/memory"
	"github.com/Rasulikus/url-shortener/internal/repository/postgres"
	urlService "github.com/Rasulikus/url-shortener/internal/service/url"
	"github.com/Rasulikus/url-shortener/internal/transport/http"
	"github.com/Rasulikus/url-shortener/internal/utils/alias"
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

	var urlRepo urlService.URLRepository

	switch cfg.Storage {
	case config.StoragePostgres:
		pool, err := postgres.NewPool(cfg.DB)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize postgres pool")
		}

		urlRepo, err = postgres.NewRepository(pool)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize postgres repository")
		}
	case config.StorageMemory:
		m := memory.New()

		urlRepo, err = memory.NewRepository(m)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize memory repository")
		}
	}

	gen, err := alias.New(alias.DefaultLength)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize alias generator")
	}

	urlServ, err := urlService.NewService(cfg.BaseURL, urlRepo, gen)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize url service")
	}

	urlHandler := http.NewURLHandler(urlServ)

	r := gin.Default()

	r.GET("/:alias", urlHandler.Redirect)

	urlApi := r.Group("/api")
	{
		urlApi.POST("", urlHandler.Create)
		urlApi.GET("/:alias", urlHandler.GetLongURLByAlias)
	}

	return r
}
