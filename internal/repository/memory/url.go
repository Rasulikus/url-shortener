package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasulikus/url-shortener/internal/model"
	"github.com/Rasulikus/url-shortener/internal/repository"
)

type Repo struct {
	m *Memory
}

func NewRepository(m *Memory) (*Repo, error) {
	if m == nil {
		return nil, fmt.Errorf("memory repository is nil")
	}
	return &Repo{
		m: m,
	}, nil
}

func (r *Repo) CreateOrGet(_ context.Context, url *model.URL) (*model.URL, error) {
	r.m.mu.Lock()
	defer r.m.mu.Unlock()

	if existing, ok := r.m.byLong[url.LongURL]; ok {
		url.ID = existing.ID
		url.Alias = existing.Alias
		url.CreatedAt = existing.CreatedAt
		c := *url
		return &c, nil
	}

	if _, ok := r.m.byAlias[url.Alias]; ok {
		return nil, repository.ErrConflict
	}

	url.ID = r.m.nextID
	url.CreatedAt = time.Now().UTC()

	r.m.nextID++

	r.m.byAlias[url.Alias] = url
	r.m.byLong[url.LongURL] = url

	c := *url
	return &c, nil
}

func (r *Repo) GetByAlias(_ context.Context, alias string) (*model.URL, error) {
	r.m.mu.RLock()
	defer r.m.mu.RUnlock()

	u, ok := r.m.byAlias[alias]
	if !ok {
		return nil, repository.ErrNotFound
	}

	c := *u
	return &c, nil
}

func (r *Repo) GetLongURLByAlias(_ context.Context, alias string) (string, error) {
	r.m.mu.RLock()
	defer r.m.mu.RUnlock()

	u, ok := r.m.byAlias[alias]
	if !ok {
		return "", repository.ErrNotFound
	}

	return u.LongURL, nil
}
