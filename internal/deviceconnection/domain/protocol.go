package domain

import "strings"

var SupportedProtocols = []string{"mqtt", "http", "modbus", "opcua"}

func Supported(p string) bool {
	p = strings.ToLower(strings.TrimSpace(p))
	for _, x := range SupportedProtocols {
		if p == x {
			return true
		}
	}
	return false
}
func NormalizeProtocol(p string) string { return strings.ToLower(strings.TrimSpace(p)) }
func (c *Connection) MarkFailure()      { c.Connected = false; c.FailureCount++ }
func (c *Connection) MarkHealthy()      { c.Connected = true; c.FailureCount = 0 }
