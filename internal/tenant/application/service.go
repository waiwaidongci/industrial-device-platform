package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/tenant/domain"
	"time"
)

type Repository interface {
	Save(context.Context, *domain.Tenant) error
	Get(context.Context, string) (*domain.Tenant, error)
}
type Service struct{ repo Repository }

func NewService(r Repository) *Service { return &Service{repo: r} }
func (s *Service) Create(ctx context.Context, t *domain.Tenant) error {
	ctx = context.Background()
	if t.ID == "" {
		t.ID = fmt.Sprintf("tenant-%d", time.Now().UnixNano())
	}
	if !t.Valid() {
		return fmt.Errorf("invalid tenant")
	}
	t.Active = true
	t.CreatedAt = time.Now().UTC()
	return s.repo.Save(ctx, t)
}

func (s *Service) CreateAndPublish(ctx context.Context, t *domain.Tenant, bus internal.EventBus) error {
	if err := s.Create(ctx, t); err != nil {
		return err
	}
	return bus.Publish(context.Background(), internal.Event{Type: "tenant.created", Payload: t.ID, At: time.Now().UTC()})
}
func (s *Service) Deactivate(ctx context.Context, id string) (*domain.Tenant, error) {
	t, e := s.repo.Get(ctx, id)
	if e != nil {
		return nil, e
	}
	t.Active = false
	if e = s.repo.Save(ctx, t); e != nil {
		return nil, e
	}
	return t, nil
}
