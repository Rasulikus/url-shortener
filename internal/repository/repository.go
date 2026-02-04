package repository

import (
	"context"
	"errors"

	"github.com/Rasulikus/url-shortener/internal/domain/model"
)

var (
	ErrNotFound      = errors.New("repository: not found")
	ErrAlreadyExists = errors.New("repository: already exists")
)

type URLRepository interface {
	CreateOrGet(ctx context.Context, u *model.URL) (*model.URL, error)
	GetByAlias(ctx context.Context, alias string) (*model.URL, error)
}
