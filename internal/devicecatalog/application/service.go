package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/devicecatalog/domain"
	"strings"
	"time"
)

type Service struct {
	repo domain.Repository
	now  func() time.Time
}

func NewService(r domain.Repository) *Service {
	return &Service{repo: r, now: func() time.Time { return time.Now().UTC() }}
}
func (s *Service) Register(ctx context.Context, d *domain.Device) error {
	if d.ID == "" {
		d.ID = fmt.Sprintf("dev-%d", s.now().UnixNano())
	}
	if d.Tags == nil {
		d.Tags = map[string]string{}
	}
	if d.Parameters == nil {
		d.Parameters = map[string]any{}
	}
	if d.Status == "" {
		d.Status = domain.StatusOffline
	}
	if err := d.Validate(); err != nil {
		return fmt.Errorf("validate device: %w", err)
	}
	d.CreatedAt = s.now()
	d.UpdatedAt = d.CreatedAt
	return s.repo.Create(ctx, d)
}
func (s *Service) Get(ctx context.Context, id string) (*domain.Device, error) {
	return s.repo.Get(ctx, strings.TrimSpace(id))
}
func (s *Service) List(ctx context.Context, t, g string) ([]*domain.Device, error) {
	return s.repo.List(ctx, t, g)
}
func (s *Service) Update(ctx context.Context, d *domain.Device) error {
	if err := d.Validate(); err != nil {
		return fmt.Errorf("validate device: %w", err)
	}
	return s.repo.Update(ctx, d)
}
func (s *Service) Heartbeat(ctx context.Context, id string) (*domain.Device, error) {
	d, e := s.Get(ctx, id)
	if e != nil {
		return nil, e
	}
	d.Status = domain.StatusOnline
	d.LastSeen = s.now()
	if e = s.repo.Update(ctx, d); e != nil {
		return nil, fmt.Errorf("heartbeat: %w", e)
	}
	return d, nil
}
func (s *Service) SetStatus(ctx context.Context, id string, status domain.Status) (*domain.Device, error) {
	d, e := s.Get(ctx, id)
	if e != nil {
		return nil, e
	}
	d.Status = status
	if e = s.repo.Update(ctx, d); e != nil {
		return nil, fmt.Errorf("set status: %w", e)
	}
	return d, nil
}
