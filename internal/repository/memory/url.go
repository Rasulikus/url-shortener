package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasulikus/url-shortener/internal/domain/model"
	"github.com/Rasulikus/url-shortener/internal/repository"
)

type Repo struct {
	m *Memory
}

func NewRepo(m *Memory) (*Repo, error) {
	if m == nil {
		return nil, fmt.Errorf("memory repository is nil")
	}
	return &Repo{
		m: m,
	}, nil
}

func (r *Repo) CreateOrGet(ctx context.Context, url *model.URL) (*model.URL, error) {
	r.m.mu.Lock()
	defer r.m.mu.Unlock()

	if existing, ok := r.m.byLong[url.LongURL]; ok {
		c := *existing
		return &c, nil
	}

	if _, ok := r.m.byAlias[url.Alias]; ok {
		return nil, repository.ErrAlreadyExists
	}

	u := &model.URL{
		ID:        r.m.nextID,
		LongURL:   url.LongURL,
		Alias:     url.Alias,
		CreatedAt: time.Now().UTC(),
	}
	r.m.nextID++

	r.m.byAlias[url.Alias] = u
	r.m.byLong[url.LongURL] = u
	c := *u
	return &c, nil
}
