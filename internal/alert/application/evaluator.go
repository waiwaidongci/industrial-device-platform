package application

import (
	"fmt"
	"github.com/example/industrial-device-platform/internal/alert/domain"
	"sync"
	"time"
)

type Evaluator struct {
	mu    sync.RWMutex
	rules map[string]domain.Rule
}

func NewEvaluator() *Evaluator { return &Evaluator{rules: map[string]domain.Rule{}} }
func (e *Evaluator) Put(r domain.Rule) error {
	if r.ID == "" || r.Name == "" || r.Metric == "" {
		return fmt.Errorf("rule id, name and metric required")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}
	e.rules[r.ID] = r
	return nil
}
func (e *Evaluator) Get(id string) (domain.Rule, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	r, ok := e.rules[id]
	return r, ok
}
func (e *Evaluator) Evaluate(id string, v float64) (domain.Event, bool) {
	r, ok := e.Get(id)
	if !ok || !r.Matches(v) {
		return domain.Event{}, false
	}
	return domain.Event{ID: fmt.Sprintf("evt-%d", time.Now().UnixNano()), RuleID: id, Value: v, Status: "open", Message: fmt.Sprintf("%s crossed %.2f", r.Metric, r.Threshold), CreatedAt: time.Now().UTC()}, true
}
