package domain

import (
	"context"
	"errors"
	"time"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusDispatched Status = "dispatched"
	StatusSucceeded  Status = "succeeded"
	StatusFailed     Status = "failed"
)

type Command struct {
	ID             string         `json:"id"`
	DeviceID       string         `json:"device_id"`
	Type           string         `json:"type"`
	Payload        map[string]any `json:"payload"`
	IdempotencyKey string         `json:"idempotency_key"`
	Status         Status         `json:"status"`
	Result         string         `json:"result,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	History        []Status       `json:"history,omitempty"`
}

func (c Command) Validate() error {
	if c.DeviceID == "" {
		return errors.New("device_id is required")
	}
	if c.Type == "" {
		return errors.New("type is required")
	}
	return nil
}

type Repository interface {
	Create(context.Context, *Command) error
	Get(context.Context, string) (*Command, error)
	FindByKey(context.Context, string) (*Command, error)
	Update(context.Context, *Command) error
	List(context.Context, string) ([]*Command, error)
}
