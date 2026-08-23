package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/tenant/domain"
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	limits map[string]domain.Limits
	counts map[string]int
	window time.Time
}

func NewLimiter() *Limiter {
	return &Limiter{limits: map[string]domain.Limits{}, counts: map[string]int{}, window: time.Now()}
}
func (l *Limiter) Set(id string, v domain.Limits) error {
	if e := v.Validate(); e != nil {
		return e
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.limits[id] = v
	return nil
}
func (l *Limiter) AllowCommand(ctx context.Context, id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if time.Since(l.window) >= time.Minute {
		l.counts = map[string]int{}
		l.window = time.Now()
	}
	limit := l.limits[id]
	if !limit.AllowsCommands(l.counts[id]) {
		return fmt.Errorf("tenant command rate exceeded")
	}
	l.counts[id]++
	return nil
}
