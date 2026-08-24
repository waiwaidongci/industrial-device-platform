package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal"
	"github.com/example/industrial-device-platform/internal/devicecatalog/domain"
)

type EventPublisher interface {
	Publish(context.Context, internal.Event) error
}

func PublishStatus(ctx context.Context, p EventPublisher, d *domain.Device) error {
	return p.Publish(ctx, internal.Event{Type: "device.status_changed", DeviceID: d.ID, Payload: map[string]any{"status": d.Status}})
}
func EnsureTenant(d *domain.Device) error {
	if d.TenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}
	return nil
}

func EnsureTags(d *domain.Device) error {
	if d.Tags == nil {
		d.Tags = map[string]string{}
	}
	return nil
}
