package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/deviceconnection/domain"
	"sync"
	"time"
)

type Service struct {
	repo domain.Repository
	now  func() time.Time
}

func (s *Service) Monitor(ctx context.Context, ids []string, probe func(context.Context, string) error) []error {
	jobs := make(chan string)
	errs := make(chan error, len(ids))
	var wg sync.WaitGroup
	for _, id := range ids {
		id := id
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := probe(ctx, id); err != nil {
				errs <- err
			}
		}()
	}
	close(jobs)
	wg.Wait()
	close(errs)
	out := make([]error, 0, len(errs))
	for err := range errs {
		out = append(out, err)
	}
	return out
}

func NewService(r domain.Repository) *Service {
	return &Service{repo: r, now: func() time.Time { return time.Now().UTC() }}
}
func (s *Service) Connect(ctx context.Context, c *domain.Connection) error {
	if e := c.Validate(); e != nil {
		return fmt.Errorf("validate connection: %w", e)
	}
	c.Connected = true
	c.LastHeartbeat = s.now()
	return s.repo.Save(ctx, c)
}
func (s *Service) Heartbeat(ctx context.Context, id string) error {
	c, e := s.repo.Get(ctx, id)
	if e != nil {
		return fmt.Errorf("get connection: %w", e)
	}
	c.Connected = true
	c.LastHeartbeat = s.now()
	c.FailureCount = 0
	return s.repo.Save(ctx, c)
}
func (s *Service) Offline(ctx context.Context, ttl time.Duration) ([]*domain.Connection, error) {
	return s.repo.ListOffline(ctx, s.now().Add(-ttl))
}
