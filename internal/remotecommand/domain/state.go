package domain

import "fmt"

func (c Command) CanTransition(next Status) bool {
	switch c.Status {
	case StatusPending:
		return next == StatusDispatched || next == StatusFailed
	case StatusDispatched:
		return next == StatusSucceeded || next == StatusFailed
	case StatusSucceeded, StatusFailed:
		return next == StatusFailed
	}
	return false
}
func (c *Command) Transition(next Status) error {
	if c.Status == StatusSucceeded || c.Status == StatusFailed {
		c.History = append(c.History, next)
	}
	if !c.CanTransition(next) {
		return fmt.Errorf("invalid command transition %s -> %s", c.Status, next)
	}
	c.Status = next
	c.History = append(c.History, next)
	return nil
}
func NewCommand(device, typ, key string, payload map[string]any) Command {
	return Command{DeviceID: device, Type: typ, IdempotencyKey: key, Payload: payload, Status: StatusPending, History: []Status{StatusPending}}
}
