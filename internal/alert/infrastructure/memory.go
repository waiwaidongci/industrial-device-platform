package infrastructure

import (
	"context"
	"github.com/example/industrial-device-platform/internal/alert/domain"
	"sync"
)

type EventStore struct {
	mu     sync.RWMutex
	events map[string]domain.Event
}

func NewEventStore() *EventStore { return &EventStore{events: map[string]domain.Event{}} }
func (s *EventStore) Put(ctx context.Context, e domain.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.ID] = e
}
func (s *EventStore) Get(ctx context.Context, id string) (domain.Event, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.events[id]
	return e, ok
}
func (s *EventStore) List(ctx context.Context, device string) []domain.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o := []domain.Event{}
	for _, e := range s.events {
		if device == "" || e.DeviceID == device {
			o = append(o, e)
		}
	}
	return o
}
