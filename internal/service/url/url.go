package url

import (
	"context"
	"errors"
	"strings"

	"github.com/Rasulikus/url-shortener/internal/model"
	"github.com/Rasulikus/url-shortener/internal/repository"
	"github.com/Rasulikus/url-shortener/internal/service"
	"github.com/Rasulikus/url-shortener/internal/utils/generator"
	"github.com/Rasulikus/url-shortener/internal/utils/validate"
	"github.com/rs/zerolog/log"
)

type URLRepository interface {
	// GetLastID возвращает последний использованный ID.
	// Используется для инициализации генератора алиасов.
	GetLastID(ctx context.Context) (uint64, error)

	// CreateOrGet создаёт новую запись с длинным URL и алиасом.
	// Если такой long URL уже существует, возвращает существующую запись.
	// Может вернуть ErrConflict при конфликте уникальности.
	CreateOrGet(ctx context.Context, u *model.URL) (*model.URL, error)

	// GetLongURLByAlias возвращает длинный URL по алиасу.
	// Если алиас не найден, возвращает ErrNotFound.
	GetLongURLByAlias(ctx context.Context, alias string) (string, error)
}

type AliasGenerator interface {
	// NewAlias генерирует новый алиас для ссылки.
	NewAlias() (string, error)
}

type Service struct {
	baseURL string
	gen     AliasGenerator
	urlRepo URLRepository
}

func NewService(baseUrl string, gen AliasGenerator, urlRepo URLRepository) (*Service, error) {
	log.Info().Msg("starting new URL Service")

	log.Info().
		Str("base_url", baseUrl).
		Int("alias_length", generator.DefaultLength).
		Msg("url service initialized")

	return &Service{
		baseURL: baseUrl,
		gen:     gen,
		urlRepo: urlRepo,
	}, nil
}

func (s *Service) CreateOrGet(ctx context.Context, longURL string) (string, error) {
	longURL = strings.TrimSpace(longURL)

	if err := validate.URL(longURL); err != nil {
		log.Debug().
			Str("url", longURL).
			Err(err).
			Msg("invalid url")

		return "", service.ErrInvalidInput
	}

	alias, err := s.gen.NewAlias()
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to generate alias")

		return "", service.ErrInternalError
	}

	u, err := s.urlRepo.CreateOrGet(ctx, &model.URL{
		LongURL: longURL,
		Alias:   alias,
	})
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			log.Warn().
				Err(err).
				Str("alias", alias).
				Str("url", longURL).
				Msg("conflict while creating url")

			return "", service.ErrConflict
		}

		log.Error().
			Err(err).
			Str("alias", alias).
			Str("url", longURL).
			Msg("failed to create url")

		return "", service.ErrInternalError
	}

	log.Info().
		Int64("id", u.ID).
		Str("alias", u.Alias).
		Str("long_url", u.LongURL).
		Msg("url created")

	return s.baseURL + "/" + u.Alias, nil
}

func (s *Service) GetLongURLByAlias(ctx context.Context, a string) (string, error) {
	longURL, err := s.urlRepo.GetLongURLByAlias(ctx, a)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn().
				Str("alias", a).
				Msg("alias not found")

			return "", service.ErrNotFound
		}
		log.Error().
			Err(err).
			Str("alias", a).
			Msg("failed to get url by alias")

		return "", service.ErrInternalError
	}
	log.Info().
		Str("long_url", longURL).
		Msg("find url")

	return longURL, nil
}
