package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/firmware/adapter"
	"github.com/example/industrial-device-platform/internal/firmware/application"
	"github.com/example/industrial-device-platform/internal/firmware/domain"
)

type missingReleaseRepo struct{}

func (missingReleaseRepo) Save(context.Context, *domain.Release) error { return nil }
func (missingReleaseRepo) Get(context.Context, string) (*domain.Release, error) {
	return nil, internal.ErrNotFound
}
func (missingReleaseRepo) List(context.Context) ([]*domain.Release, error) { return nil, nil }

func TestReleaseLookupErrorClassification(t *testing.T) {
	svc := application.NewService(missingReleaseRepo{})
	if _, err := svc.Start(context.Background(), "missing"); !errors.Is(err, internal.ErrNotFound) {
		t.Fatalf("not-found chain lost: %v", err)
	} else if application.ClassifyRolloutError(err) != "not_found" {
		t.Fatalf("wrong classification: %v", err)
	}
	if adapter.RolloutStatus("not_found") != 404 {
		t.Fatal("adapter did not preserve not-found status")
	}
}
