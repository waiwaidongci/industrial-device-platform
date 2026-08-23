package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/deviceconnection/domain"
	"time"
)

type Service struct {
	repo domain.Repository
	now  func() time.Time
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
