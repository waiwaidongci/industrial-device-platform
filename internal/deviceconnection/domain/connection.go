package domain

import (
	"context"
	"errors"
	"time"
)

type Connection struct {
	DeviceID      string    `json:"device_id"`
	Protocol      string    `json:"protocol"`
	Address       string    `json:"address"`
	Connected     bool      `json:"connected"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	FailureCount  int       `json:"failure_count"`
}

func (c Connection) Validate() error {
	if c.DeviceID == "" {
		return errors.New("device id required")
	}
	if c.Protocol == "" {
		return errors.New("protocol required")
	}
	return nil
}

type Repository interface {
	Save(context.Context, *Connection) error
	Get(context.Context, string) (*Connection, error)
	ListOffline(context.Context, time.Time) ([]*Connection, error)
}
