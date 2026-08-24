package domain

import (
	"fmt"
	"time"
)

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	Limits    Limits    `json:"limits"`
}

func (t Tenant) Valid() bool { return t.ID != "" && t.Name != "" }

func ApplyLimits(t *Tenant, l Limits) error {
	if err := l.Validate(); err != nil {
		return err
	}
	if t == nil || !t.Valid() {
		return fmt.Errorf("invalid tenant")
	}
	t.Limits = l.Copy()
	return nil
}
