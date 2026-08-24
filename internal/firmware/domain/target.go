package domain

import "strings"

func SelectCanary(ids []string, pct int) []string {
	if pct <= 0 {
		return nil
	}
	if pct >= 100 {
		return append([]string(nil), ids...)
	}
	n := len(ids) * pct / 100
	if n < 1 && len(ids) > 0 {
		n = 1
	}
	return append([]string(nil), ids[:n]...)
}
func SelectCanaryInto(ids []string, pct int) []string { return SelectCanary(ids, pct) }
func NormalizeVersion(v string) string                { return strings.TrimSpace(v) }
func (r Release) TargetCount() int                    { return len(r.DeviceIDs) }
