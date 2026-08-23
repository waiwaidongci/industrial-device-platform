package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/example/industrial-device-platform/internal/remotecommand/domain"
	"time"
)

type LocalExecutor struct{}

func (LocalExecutor) Execute(ctx context.Context, c *domain.Command) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}
	b, e := json.Marshal(map[string]any{"command": c.Type, "payload": c.Payload, "device_id": c.DeviceID})
	if e != nil {
		return "", fmt.Errorf("encode result: %w", e)
	}
	return string(b), nil
}

type RetryPolicy struct {
	MaxAttempts int
	BaseDelay   time.Duration
}

func (p RetryPolicy) Delay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	d := p.BaseDelay
	for i := 1; i < attempt; i++ {
		d *= 2
	}
	return d
}
