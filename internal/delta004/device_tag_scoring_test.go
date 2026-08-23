package application_test

import (
	"context"
	"github.com/example/industrial-device-platform/internal"
	"testing"

	"github.com/example/industrial-device-platform/internal/devicecatalog/adapter"
	"github.com/example/industrial-device-platform/internal/devicecatalog/application"
	"github.com/example/industrial-device-platform/internal/devicecatalog/domain"
	"github.com/example/industrial-device-platform/internal/devicecatalog/infrastructure"
)

type nilTagRepo struct{ device *domain.Device }

func (r *nilTagRepo) Create(context.Context, *domain.Device) error { return nil }
func (r *nilTagRepo) Get(context.Context, string) (*domain.Device, error) {
	x := *r.device
	return &x, nil
}
func (r *nilTagRepo) List(context.Context, string, string) ([]*domain.Device, error) { return nil, nil }
func (r *nilTagRepo) Update(context.Context, *domain.Device) error                   { return nil }
func (r *nilTagRepo) Delete(context.Context, string) error                           { return nil }

func TestDeviceTagUpdateNilMap(t *testing.T) {
	if got := domain.NormalizeTags(nil); got == nil {
		t.Fatal("normalization returned nil map")
	}
	if err := adapter.ValidateTagKey(" "); err == nil {
		t.Fatal("blank tag key was accepted")
	}
	repo := infrastructure.NewMemoryRepository()
	svc := application.NewService(repo)
	d := &domain.Device{TenantID: "plant", Name: "press", Kind: "plc"}
	if err := svc.Register(context.Background(), d); err != nil {
		t.Fatal(err)
	}
	updated, err := svc.SetTag(context.Background(), d.ID, " Area ", " line-1 ")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Tags["area"] != "line-1" {
		t.Fatalf("unexpected tags: %#v", updated.Tags)
	}
	if err := application.EnsureTags(updated); err != nil || updated.Tags == nil {
		t.Fatalf("tags not writable: %v", err)
	}
	empty := &domain.Device{}
	if err := application.EnsureTags(empty); err != nil || empty.Tags == nil {
		t.Fatalf("zero device tags were not initialized: %v", err)
	}
	nilRepo := &nilTagRepo{device: &domain.Device{ID: "nil-tags", Name: "press", Kind: "plc"}}
	nilSvc := application.NewService(nilRepo)
	if _, err := nilSvc.SetTag(context.Background(), "nil-tags", "zone", "line-2"); err != nil {
		t.Fatalf("nil tag path still failed: %v", err)
	}
	if internal.ErrNotFound == nil {
		t.Fatal("sentinel missing")
	}
}
