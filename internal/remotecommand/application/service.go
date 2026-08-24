package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/remotecommand/domain"
	"time"
)

type Service struct {
	repo domain.Repository
	now  func() time.Time
}

func NewService(r domain.Repository) *Service {
	return &Service{repo: r, now: func() time.Time { return time.Now().UTC() }}
}
func (s *Service) Submit(ctx context.Context, c *domain.Command) (*domain.Command, error) {
	if e := c.Validate(); e != nil {
		return nil, fmt.Errorf("validate command: %w", e)
	}
	if c.IdempotencyKey != "" {
		if old, e := s.repo.FindByKey(ctx, c.IdempotencyKey); e == nil {
			return old, nil
		}
	}
	c.ID = fmt.Sprintf("cmd-%d", s.now().UnixNano())
	c.Status = domain.StatusPending
	c.History = []domain.Status{domain.StatusPending}
	c.CreatedAt = s.now()
	c.UpdatedAt = c.CreatedAt
	if e := s.repo.Create(ctx, c); e != nil {
		return nil, fmt.Errorf("create command: %w", e)
	}
	return c, nil
}
func (s *Service) Get(ctx context.Context, id string) (*domain.Command, error) {
	return s.repo.Get(ctx, id)
}
func (s *Service) List(ctx context.Context, id string) ([]*domain.Command, error) {
	return s.repo.List(ctx, id)
}
func (s *Service) Complete(ctx context.Context, id string, status domain.Status, result string) (*domain.Command, error) {
	c, e := s.Get(ctx, id)
	if e != nil {
		return nil, e
	}
	if e = c.Transition(status); e != nil {
		return nil, fmt.Errorf("complete command: %w", e)
	}
	c.Result = result
	c.UpdatedAt = s.now()
	if e = s.repo.Update(ctx, c); e != nil {
		return nil, fmt.Errorf("complete command: %w", e)
	}
	return c, nil
}
