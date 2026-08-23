package infrastructure

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/remotecommand/domain"
	"sync"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Command
	keys map[string]string
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{data: map[string]*domain.Command{}, keys: map[string]string{}}
}
func cp(c *domain.Command) *domain.Command { x := *c; return &x }
func (r *MemoryRepository) Create(ctx context.Context, c *domain.Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[c.ID]; ok {
		return internal.ErrConflict
	}
	if c.IdempotencyKey != "" {
		if _, ok := r.keys[c.IdempotencyKey]; ok {
			return fmt.Errorf("idempotency key: %w", internal.ErrConflict)
		}
		r.keys[c.IdempotencyKey] = c.ID
	}
	r.data[c.ID] = cp(c)
	return nil
}
func (r *MemoryRepository) Get(ctx context.Context, id string) (*domain.Command, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.data[id]
	if !ok {
		return nil, internal.ErrNotFound
	}
	return cp(c), nil
}
func (r *MemoryRepository) FindByKey(ctx context.Context, k string) (*domain.Command, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.keys[k]
	if !ok {
		return nil, internal.ErrNotFound
	}
	return cp(r.data[id]), nil
}
func (r *MemoryRepository) Update(ctx context.Context, c *domain.Command) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[c.ID]; !ok {
		return internal.ErrNotFound
	}
	r.data[c.ID] = cp(c)
	return nil
}
func (r *MemoryRepository) List(ctx context.Context, id string) ([]*domain.Command, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o := []*domain.Command{}
	for _, c := range r.data {
		if id == "" || c.DeviceID == id {
			o = append(o, cp(c))
		}
	}
	return o, nil
}
