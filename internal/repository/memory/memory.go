package memory

import (
	"sync"

	"github.com/Rasulikus/url-shortener/internal/domain/model"
)

type Memory struct {
	mu      sync.RWMutex
	byAlias map[string]*model.URL
	byLong  map[string]*model.URL
	nextID  int64
}

func New() *Memory {
	return &Memory{
		byAlias: make(map[string]*model.URL),
		byLong:  make(map[string]*model.URL),
		nextID:  1,
	}
}
