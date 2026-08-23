package application

import (
	"github.com/example/industrial-device-platform/internal/alert/domain"
	"sync"
	"time"
)

type Deduplicator struct {
	mu   sync.Mutex
	seen map[string]time.Time
	ttl  time.Duration
}

func NewDeduplicator(ttl time.Duration) *Deduplicator {
	return &Deduplicator{seen: map[string]time.Time{}, ttl: ttl}
}
func (d *Deduplicator) Accept(e domain.Event) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	key := e.RuleID + ":" + e.DeviceID
	if t, ok := d.seen[key]; ok && time.Since(t) < d.ttl {
		return false
	}
	d.seen[key] = time.Now()
	return true
}
func (d *Deduplicator) Cleanup() {
	d.mu.Lock()
	defer d.mu.Unlock()
	for k, t := range d.seen {
		if time.Since(t) >= d.ttl {
			delete(d.seen, k)
		}
	}
}
