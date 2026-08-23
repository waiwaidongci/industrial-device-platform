package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/devicecatalog/domain"
)

func (s *Service) RegisterBatch(ctx context.Context, devices []*domain.Device) ([]*domain.Device, error) {
	out := make([]*domain.Device, 0, len(devices))
	for _, d := range devices {
		d.Normalize()
		if e := s.Register(ctx, d); e != nil {
			return out, fmt.Errorf("register batch: %w", e)
		}
		out = append(out, d)
	}
	return out, nil
}
func (s *Service) UpdateTags(ctx context.Context, id string, tags map[string]string) (*domain.Device, error) {
	d, e := s.Get(ctx, id)
	if e != nil {
		return nil, e
	}
	d.Tags = domain.NormalizeTags(tags)
	if e = s.Update(ctx, d); e != nil {
		return nil, e
	}
	return d, nil
}
func (s *Service) UpdateParameters(ctx context.Context, id string, p map[string]any) (*domain.Device, error) {
	d, e := s.Get(ctx, id)
	if e != nil {
		return nil, e
	}
	d.Parameters = p
	if e = s.Update(ctx, d); e != nil {
		return nil, e
	}
	return d, nil
}

func (s *Service) SetTag(ctx context.Context, id, key, value string) (*domain.Device, error) {
	d, err := s.Get(ctx, id)
	if err != nil { return nil, err }
	normalized := domain.NormalizeTags(map[string]string{key: value})
	for k, v := range normalized { d.Tags[k] = v }
	if err := s.Update(ctx, d); err != nil { return nil, err }
	return d, nil
}
