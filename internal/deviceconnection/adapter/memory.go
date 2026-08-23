package adapter

import (
	"context"
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/deviceconnection/domain"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Connection
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{data: map[string]*domain.Connection{}}
}
func (r *MemoryRepository) Save(ctx context.Context, c *domain.Connection) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	x := *c
	r.data[c.DeviceID] = &x
	return nil
}
func (r *MemoryRepository) Get(ctx context.Context, id string) (*domain.Connection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.data[id]
	if !ok {
		return nil, internal.ErrNotFound
	}
	x := *c
	return &x, nil
}
func (r *MemoryRepository) ListOffline(ctx context.Context, t time.Time) ([]*domain.Connection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o := []*domain.Connection{}
	for _, c := range r.data {
		if c.LastHeartbeat.Before(t) {
			x := *c
			o = append(o, &x)
		}
	}
	return o, nil
}
