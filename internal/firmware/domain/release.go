package domain

import "time"

type Release struct {
	ID            string    `json:"id"`
	Version       string    `json:"version"`
	DeviceIDs     []string  `json:"device_ids"`
	CanaryPercent int       `json:"canary_percent"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

func (r Release) Clone() Release { return r }

func (r Release) Valid() bool {
	return r.Version != "" && r.CanaryPercent >= 0 && r.CanaryPercent <= 100
}
