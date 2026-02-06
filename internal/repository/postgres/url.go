package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rasulikus/url-shortener/internal/model"
	"github.com/Rasulikus/url-shortener/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) (*Repo, error) {
	if pool == nil {
		return nil, errors.New("repository: pgx pool is nil")
	}

	return &Repo{
		pool: pool,
	}, nil
}

func (r *Repo) CreateOrGet(ctx context.Context, u *model.URL) (*model.URL, error) {
	const q = `
	INSERT INTO urls (long_url, alias)
	VALUES ($1, $2)
	ON CONFLICT (long_url) DO UPDATE
	SET long_url = excluded.long_url
	RETURNING id, long_url, alias, created_at;
`

	err := r.pool.QueryRow(ctx, q, u.LongURL, u.Alias).Scan(&u.ID, &u.LongURL, &u.Alias, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, repository.ErrConflict
		}
		return nil, fmt.Errorf("repository: insert u: %w", err)
	}

	return u, nil
}

func (r *Repo) GetByAlias(ctx context.Context, alias string) (*model.URL, error) {
	const q = `
	SELECT id, long_url, alias, created_at FROM urls WHERE alias = $1;
`

	url := new(model.URL)

	err := r.pool.QueryRow(ctx, q, alias).Scan(&url.ID, &url.LongURL, &url.Alias, &url.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("repository: select url by alias: %w", err)
	}

	return url, nil
}

func (r *Repo) GetLongURLByAlias(ctx context.Context, alias string) (string, error) {
	const q = `
	SELECT long_url FROM urls WHERE alias = $1;
`

	var longURL string

	err := r.pool.QueryRow(ctx, q, alias).Scan(&longURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNotFound
		}
		return "", fmt.Errorf("repository: select longURL by alias: %w", err)
	}

	return longURL, nil
}
