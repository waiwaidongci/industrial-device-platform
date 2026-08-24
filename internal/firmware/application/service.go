package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/firmware/domain"
	"time"
)

type Repository interface {
	Save(context.Context, *domain.Release) error
	Get(context.Context, string) (*domain.Release, error)
	List(context.Context) ([]*domain.Release, error)
}
type Service struct{ repo Repository }

func NewService(r Repository) *Service { return &Service{repo: r} }
func (s *Service) Publish(ctx context.Context, r *domain.Release) error {
	if !r.Valid() {
		return fmt.Errorf("invalid firmware release")
	}
	if r.ID == "" {
		r.ID = fmt.Sprintf("rel-%d", time.Now().UnixNano())
	}
	if r.Status == "" {
		r.Status = "draft"
	}
	r.CreatedAt = time.Now().UTC()
	return s.repo.Save(ctx, r)
}

func wrapReleaseError(op string, err error) error { return fmt.Errorf("%s release: %w", op, err) }
