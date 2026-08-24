package domain

import "time"

type DeviceEvent struct {
	ID       string         `json:"id"`
	DeviceID string         `json:"device_id"`
	Type     string         `json:"type"`
	Payload  map[string]any `json:"payload"`
	At       time.Time      `json:"at"`
}

func NewDeviceEvent(id, t string, p map[string]any) DeviceEvent {
	return DeviceEvent{ID: id, DeviceID: id, Type: t, Payload: p, At: time.Now().UTC()}
}
func (e DeviceEvent) Valid() bool        { return e.DeviceID != "" && e.Type != "" }
func (e DeviceEvent) Age() time.Duration { return time.Since(e.At) }
