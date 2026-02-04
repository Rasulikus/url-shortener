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
)

type Service struct {
	gen     *alias.Generator
	urlRepo repository.URLRepository
}

func NewService(urlRepo repository.URLRepository) (*Service, error) {
	gen, err := alias.New(alias.DefaultLength)
	if err != nil {
		return nil, fmt.Errorf("service: failed to generate alias generator: %w", err)
	}
	return &Service{
		gen:     gen,
		urlRepo: urlRepo,
	}, nil
}

func (s *Service) CreateOrGet(ctx context.Context, longURL string) (*model.URL, error) {
	longURL = strings.TrimSpace(longURL)
	err := validate.URL(longURL)
	if err != nil {
		return nil, service.ErrInvalidInput
	}
	for i := 0; i < 3; i++ {
		a, err := s.gen.NewAlias()
		if err != nil {
			return nil, fmt.Errorf("service: failed to generate alias: %w", err)
		}

		u, err := s.urlRepo.CreateOrGet(ctx, &model.URL{
			LongURL: longURL,
			Alias:   a,
		})
		if err != nil {
			if errors.Is(err, repository.ErrAlreadyExists) {
				continue
			}
			return nil, fmt.Errorf("service: failed create alias in db: %w", err)
		}
		return u, nil
	}
	return nil, service.ErrAliasCollision
}

func (s *Service) GetByAlias(ctx context.Context, a string) (*model.URL, error) {
	u, err := s.urlRepo.GetByAlias(ctx, a)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, service.ErrNotFound
		}
		return nil, fmt.Errorf("service: failed get a from db: %w", err)
	}
	return u, nil
}
