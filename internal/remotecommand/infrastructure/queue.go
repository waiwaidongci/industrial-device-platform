package infrastructure

import (
	"context"
	"github.com/example/industrial-device-platform/internal/remotecommand/domain"
	"sync"
)

type Queue struct {
	mu    sync.Mutex
	items []*domain.Command
}

func NewQueue() *Queue { return &Queue{items: []*domain.Command{}} }
func (q *Queue) Enqueue(ctx context.Context, c *domain.Command) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	x := *c
	q.items = append(q.items, &x)
	return nil
}
func (q *Queue) Dequeue(ctx context.Context) (*domain.Command, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil, false
	}
	c := q.items[0]
	q.items = q.items[1:]
	return c, true
}
func (q *Queue) Len() int { q.mu.Lock(); defer q.mu.Unlock(); return len(q.items) }

func (q *Queue) Drain(ctx context.Context, fn func(*domain.Command) error) error {
	for {
		c, ok := q.Dequeue(ctx)
		if !ok {
			return nil
		}
		if err := fn(c); err != nil {
			return err
		}
	}
}
