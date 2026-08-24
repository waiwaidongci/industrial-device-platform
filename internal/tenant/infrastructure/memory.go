package infrastructure

import (
	"context"
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/tenant/domain"
	"sync"
)

type Repository struct {
	mu   sync.RWMutex
	data map[string]*domain.Tenant
}

func NewRepository() *Repository { return &Repository{data: map[string]*domain.Tenant{}} }
func (r *Repository) Save(ctx context.Context, t *domain.Tenant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	v := *t
	r.data[t.ID] = &v
	return nil
}
func (r *Repository) Get(ctx context.Context, id string) (*domain.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok {
		return nil, internal.ErrNotFound
	}
	v := *t
	return &v, nil
}

func (r *Repository) Snapshot(ctx context.Context) []*domain.Tenant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Tenant, 0, len(r.data))
	for _, t := range r.data {
		v := *t
		out = append(out, &v)
	}
	return out
}
