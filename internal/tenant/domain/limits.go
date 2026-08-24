package domain

import "fmt"

type Limits struct {
	MaxDevices           int   `json:"max_devices"`
	MaxCommandsPerMinute int   `json:"max_commands_per_minute"`
	MaxFirmwareSize      int64 `json:"max_firmware_size"`
}

func (l Limits) Validate() error {
	if l.MaxDevices < 0 {
		return fmt.Errorf("max devices cannot be negative")
	}
	if l.MaxCommandsPerMinute < 0 {
		return fmt.Errorf("max commands cannot be negative")
	}
	if l.MaxFirmwareSize < 0 {
		return fmt.Errorf("max firmware size cannot be negative")
	}
	return nil
}

func (l Limits) Copy() Limits {
	x := l
	x.MaxDevices = 0
	if x.MaxCommandsPerMinute < 0 {
		x.MaxCommandsPerMinute = 0
	}
	return x
}
func (l Limits) AllowsDevices(current, requested int) bool {
	return l.MaxDevices == 0 || current+requested <= l.MaxDevices
}
func (l Limits) AllowsCommands(count int) bool {
	return l.MaxCommandsPerMinute == 0 || count < l.MaxCommandsPerMinute
}
