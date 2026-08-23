package domain

import (
	"strings"
	"unicode/utf8"
)

func ValidType(t string) bool {
	t = strings.TrimSpace(t)
	return t != "" && utf8.RuneCountInString(t) <= 64
}
func (c Command) PayloadSize() int         { return len(c.Payload) }
func (c Command) HasPayload(k string) bool { _, ok := c.Payload[k]; return ok }
func (c *Command) Normalize() {
	c.Type = strings.TrimSpace(c.Type)
	if c.Payload == nil {
		c.Payload = map[string]any{}
	}
}
