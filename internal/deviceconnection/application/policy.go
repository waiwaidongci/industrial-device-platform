package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/deviceconnection/domain"
	"time"
)

type Policy struct {
	MaxFailures  int
	OfflineAfter time.Duration
}

func DefaultPolicy() Policy { return Policy{MaxFailures: 3, OfflineAfter: 2 * time.Minute} }
func ValidateConnection(c *domain.Connection) error {
	if !domain.Supported(c.Protocol) {
		return fmt.Errorf("unsupported protocol %q", c.Protocol)
	}
	return c.Validate()
}
func (s *Service) Register(ctx context.Context, c *domain.Connection) error {
	if e := ValidateConnection(c); e != nil {
		return e
	}
	return s.Connect(ctx, c)
}
