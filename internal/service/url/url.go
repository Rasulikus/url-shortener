package url

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Rasulikus/url-shortener/internal/domain/model"
	"github.com/Rasulikus/url-shortener/internal/repository"
	"github.com/Rasulikus/url-shortener/internal/service"
	"github.com/Rasulikus/url-shortener/internal/utils/alias"
	"github.com/Rasulikus/url-shortener/internal/utils/validate"
	"github.com/rs/zerolog/log"
)

type Service struct {
	baseUrl string
	gen     *alias.Generator
	urlRepo repository.URLRepository
}

func NewService(baseUrl string, urlRepo repository.URLRepository) (*Service, error) {
	log.Info().Msg("Starting New URL Service")
	gen, err := alias.New(alias.DefaultLength)
	if err != nil {
		log.Error().Err(err).Msg("alias generator init failed")
		return nil, fmt.Errorf("service: failed to initialize alias generator: %w", err)
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
		log.Warn().
			Str("url", longURL).
			Err(err).
			Msg("invalid url")
		return "", service.ErrInvalidInput
	}
	for i := 0; i < 3; i++ {
		a, err := s.gen.NewAlias()
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to generate alias")
			return "", fmt.Errorf("service: failed to generate alias: %w", err)
		}
		u, err := s.urlRepo.CreateOrGet(ctx, &model.URL{
			LongURL: longURL,
			Alias:   a,
		})
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
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
			return "", fmt.Errorf("service: failed create alias in db: %w", err)
		}

		log.Info().
			Int64("id", u.ID).
			Str("long_url", u.LongURL).
			Str("alias", u.Alias).
			Str("created_at", u.CreatedAt.String()).
			Msg("alias url created")
		return u.Alias, nil
	}

	log.Warn().
		Str("url", longURL).
		Msg("alias collision limit reached")

	return "", service.ErrAliasCollision
}

func (s *Service) GetByAlias(ctx context.Context, a string) (*model.URL, error) {
	u, err := s.urlRepo.GetByAlias(ctx, a)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			log.Warn().
				Str("alias", a).
				Msg("alias not found")
			return nil, service.ErrNotFound
		}
		log.Error().
			Err(err).
			Str("alias", a).
			Msg("failed to get url by alias")
		return nil, fmt.Errorf("service: failed get a from db: %w", err)
	}
	log.Info().
		Int64("id", u.ID).
		Str("long_url", u.LongURL).
		Str("alias", u.Alias).
		Str("created_at", u.CreatedAt.String()).
		Msg("find url")
	return u, nil
}
