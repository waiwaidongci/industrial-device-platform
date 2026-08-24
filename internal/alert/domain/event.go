package domain

import "time"

type Event struct {
	ID         string     `json:"id"`
	RuleID     string     `json:"rule_id"`
	DeviceID   string     `json:"device_id"`
	Value      float64    `json:"value"`
	Status     string     `json:"status"`
	Message    string     `json:"message"`
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

func (e *Event) Resolve()    { now := time.Now().UTC(); e.Status = "resolved"; e.ResolvedAt = &now }
func (e Event) Active() bool { return e.Status == "open" || e.Status == "acknowledged" }
