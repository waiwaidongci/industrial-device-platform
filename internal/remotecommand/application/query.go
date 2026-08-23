package application

import (
	"context"
	"github.com/example/industrial-device-platform/internal/remotecommand/domain"
	"sort"
	"time"
)

func (s *Service) Recent(ctx context.Context, device string, n int) ([]*domain.Command, error) {
	v, e := s.List(ctx, device)
	if e != nil {
		return nil, e
	}
	sort.Slice(v, func(i, j int) bool { return v[i].CreatedAt.After(v[j].CreatedAt) })
	if n > 0 && len(v) > n {
		v = v[:n]
	}
	return v, nil
}
func (s *Service) Age(ctx context.Context, id string) (time.Duration, error) {
	c, e := s.Get(ctx, id)
	if e != nil {
		return 0, e
	}
	return time.Since(c.CreatedAt), nil
}
