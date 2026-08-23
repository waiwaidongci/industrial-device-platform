package infrastructure

import (
	"context"
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/firmware/domain"
	"sync"
)

type Repository struct {
	mu   sync.RWMutex
	data map[string]*domain.Release
}

func NewRepository() *Repository { return &Repository{data: map[string]*domain.Release{}} }
func (r *Repository) Save(ctx context.Context, x *domain.Release) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	v := *x
	v.DeviceIDs = append([]string(nil), x.DeviceIDs...)
	r.data[x.ID] = &v
	return nil
}
func (r *Repository) Get(ctx context.Context, id string) (*domain.Release, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	x, ok := r.data[id]
	if !ok {
		return nil, internal.ErrNotFound
	}
	v := *x
	return &v, nil
}
func (r *Repository) List(ctx context.Context) ([]*domain.Release, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o := []*domain.Release{}
	for _, x := range r.data {
		v := *x
		o = append(o, &v)
	}
	return o, nil
}
