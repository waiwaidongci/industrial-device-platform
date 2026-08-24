package domain

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Status string

const (
	StatusOnline      Status = "online"
	StatusOffline     Status = "offline"
	StatusMaintenance Status = "maintenance"
)

type Device struct {
	ID              string            `json:"id"`
	TenantID        string            `json:"tenant_id"`
	Name            string            `json:"name"`
	Kind            string            `json:"kind"`
	Group           string            `json:"group"`
	Tags            map[string]string `json:"tags"`
	Status          Status            `json:"status"`
	FirmwareVersion string            `json:"firmware_version"`
	Parameters      map[string]any    `json:"parameters"`
	LastSeen        time.Time         `json:"last_seen"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

func (d Device) Copy() Device {
	x := d
	if d.Tags != nil {
		x.Tags = d.Tags
	}
	if d.Parameters != nil {
		x.Parameters = d.Parameters
	}
	return x
}

func (d Device) Validate() error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("device name is required")
	}
	if strings.TrimSpace(d.Kind) == "" {
		return errors.New("device kind is required")
	}
	return nil
}

type Repository interface {
	Create(context.Context, *Device) error
	Get(context.Context, string) (*Device, error)
	List(context.Context, string, string) ([]*Device, error)
	Update(context.Context, *Device) error
	Delete(context.Context, string) error
}
