package service

import (
	"context"
	"errors"

	"github.com/Rasulikus/url-shortener/internal/domain/model"
)

var (
	ErrInvalidInput   = errors.New("service: invalid input")
	ErrNotFound       = errors.New("service: not found")
	ErrAliasCollision = errors.New("service: failed to generate unique alias")
)

type URLService interface {
	CreateOrGet(ctx context.Context, longURL string) (string, error)
	GetByAlias(ctx context.Context, alias string) (*model.URL, error)
}
