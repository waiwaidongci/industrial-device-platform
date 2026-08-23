package infrastructure

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/devicecatalog/domain"
	"sort"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Device
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{data: map[string]*domain.Device{}}
}
func clone(d *domain.Device) *domain.Device {
	x := *d
	x.Tags = map[string]string{}
	for k, v := range d.Tags {
		x.Tags[k] = v
	}
	x.Parameters = map[string]any{}
	for k, v := range d.Parameters {
		x.Parameters[k] = v
	}
	return &x
}
func (r *MemoryRepository) Create(ctx context.Context, d *domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[d.ID]; ok {
		return fmt.Errorf("create device: %w", internal.ErrConflict)
	}
	r.data[d.ID] = clone(d)
	return nil
}
func (r *MemoryRepository) Get(ctx context.Context, id string) (*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.data[id]
	if !ok {
		return nil, internal.ErrNotFound
	}
	return clone(d), nil
}
func (r *MemoryRepository) List(ctx context.Context, tenant, group string) ([]*domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []*domain.Device{}
	for _, d := range r.data {
		if tenant != "" && d.TenantID != tenant {
			continue
		}
		if group != "" && d.Group != group {
			continue
		}
		out = append(out, clone(d))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}
func (r *MemoryRepository) Update(ctx context.Context, d *domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[d.ID]; !ok {
		return internal.ErrNotFound
	}
	d.UpdatedAt = time.Now().UTC()
	r.data[d.ID] = clone(d)
	return nil
}
func (r *MemoryRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return internal.ErrNotFound
	}
	delete(r.data, id)
	return nil
}
