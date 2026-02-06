package url

import (
	"context"
	"errors"
	"strings"

	"github.com/Rasulikus/url-shortener/internal/model"
	"github.com/Rasulikus/url-shortener/internal/repository"
	"github.com/Rasulikus/url-shortener/internal/service"
	"github.com/Rasulikus/url-shortener/internal/utils/alias"
	"github.com/Rasulikus/url-shortener/internal/utils/validate"
	"github.com/rs/zerolog/log"
)

type URLRepository interface {
	CreateOrGet(ctx context.Context, u *model.URL) (*model.URL, error)
	GetLongURLByAlias(ctx context.Context, alias string) (string, error)
}

type AliasGenerator interface {
	NewAlias() (string, error)
}

type Service struct {
	baseUrl string
	gen     AliasGenerator
	urlRepo URLRepository
}

func NewService(baseUrl string, urlRepo URLRepository, gen AliasGenerator) (*Service, error) {
	log.Info().Msg("starting new LongURL Service")

	gen, err := alias.New(alias.DefaultLength)
	if err != nil {
		log.Error().Err(err).Msg("alias generator init failed")

		return nil, service.ErrInternalError
	}

	log.Info().
		Str("base_url", baseUrl).
		Int("alias_length", alias.DefaultLength).
		Msg("url service initialized")

	return &Service{
		baseUrl: baseUrl,
		gen:     gen,
		urlRepo: urlRepo,
	}, nil
}

func (s *Service) CreateOrGet(ctx context.Context, longURL string) (string, error) {
	longURL = strings.TrimSpace(longURL)

	err := validate.URL(longURL)
	if err != nil {
		log.Debug().
			Str("url", longURL).
			Err(err).
			Msg("invalid url")
		return "", service.ErrInvalidInput
	}

	for i := 0; i < 2; i++ {
		a, err := s.gen.NewAlias()
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to generate alias")

			return "", service.ErrInternalError
		}

		u, err := s.urlRepo.CreateOrGet(ctx, &model.URL{
			LongURL: longURL,
			Alias:   a,
		})
		if err != nil {
			if errors.Is(err, repository.ErrConflict) {
				log.Warn().
					Str("alias", a).
					Int("attempt", i+1).
					Msg("alias collision, retrying")

				continue
			}
			log.Error().
				Err(err).
				Str("url", longURL).
				Str("alias", a).
				Int("attempt", i+1).
				Msg("failed to create url in db")

			return "", service.ErrInternalError
		}
		log.Info().
			Int64("id", u.ID).
			Str("long_url", u.LongURL).
			Str("alias", u.Alias).
			Str("created_at", u.CreatedAt.String()).
			Msg("alias url created")

		return s.baseUrl + "/" + u.Alias, nil
	}
	log.Warn().
		Str("url", longURL).
		Msg("alias collision limit reached")

	return "", service.ErrAliasCollision
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
