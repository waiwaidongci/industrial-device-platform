package infrastructure

import (
	"context"
	"github.com/example/industrial-device-platform/internal/deviceconnection/application"
	"time"
)

type Monitor struct {
	service           *application.Service
	interval, timeout time.Duration
}

func NewMonitor(s *application.Service, i, t time.Duration) *Monitor {
	return &Monitor{service: s, interval: i, timeout: t}
}
func (m *Monitor) Run(ctx context.Context, onOffline func(string)) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			items, _ := m.service.Offline(ctx, m.timeout)
			for _, c := range items {
				onOffline(c.DeviceID)
			}
		}
	}
}
