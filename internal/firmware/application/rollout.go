package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/firmware/domain"
)

func (s *Service) Start(ctx context.Context, id string) (*domain.Release, error) {
	r, e := s.repo.Get(ctx, id)
	if e != nil {
		return nil, e
	}
	if e = r.Start(); e != nil {
		return nil, fmt.Errorf("start rollout: %w", e)
	}
	if e = s.repo.Save(ctx, r); e != nil {
		return nil, e
	}
	return r, nil
}
func (s *Service) Pause(ctx context.Context, id string) (*domain.Release, error) {
	r, e := s.repo.Get(ctx, id)
	if e != nil {
		return nil, e
	}
	if e = r.Pause(); e != nil {
		return nil, e
	}
	if e = s.repo.Save(ctx, r); e != nil {
		return nil, e
	}
	return r, nil
}
