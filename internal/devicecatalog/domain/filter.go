package domain

import "strings"

type Filter struct {
	TenantID string
	Group    string
	Status   Status
	Kind     string
	TagKey   string
	TagValue string
}

func (f Filter) Match(d *Device) bool {
	if f.TenantID != "" && d.TenantID != f.TenantID {
		return false
	}
	if f.Group != "" && d.Group != f.Group {
		return false
	}
	if f.Status != "" && d.Status != f.Status {
		return false
	}
	if f.Kind != "" && d.Kind != f.Kind {
		return false
	}
	if f.TagKey != "" && d.Tags[f.TagKey] != f.TagValue {
		return false
	}
	return true
}
func NormalizeTags(tags map[string]string) map[string]string {
	out := make(map[string]string, len(tags))
	for k, v := range tags {
		key := strings.ToLower(strings.TrimSpace(k))
		val := strings.TrimSpace(v)
		if key != "" && val != "" {
			out[key] = val
		}
	}
	return out
}
func (d *Device) Normalize() {
	d.Name = strings.TrimSpace(d.Name)
	d.Kind = strings.ToLower(strings.TrimSpace(d.Kind))
	d.Group = strings.TrimSpace(d.Group)
	d.Tags = NormalizeTags(d.Tags)
}
