package memory

import (
	"context"
	"movie-app/metadata/internal/v1/repository"
	"movie-app/metadata/pkg/model"
	"sync"
)

// Repository is a struct that represents a rating handler
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

func New() *Repository {
	return &Repository{data: map[string]*model.Metadata{}}
}

// Get retrieves api rating for by api id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
