package domain

import "time"

type Rule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Metric    string    `json:"metric"`
	Threshold float64   `json:"threshold"`
	Severity  string    `json:"severity"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

func (r Rule) Matches(v float64) bool { return r.Enabled && v >= r.Threshold }
