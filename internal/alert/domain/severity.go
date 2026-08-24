package domain

type Severity int

const (
	SeverityInfo Severity = iota
	SeverityWarning
	SeverityCritical
)

func ParseSeverity(v string) Severity {
	switch v {
	case "critical":
		return SeverityCritical
	case "warning":
		return SeverityWarning
	default:
		return SeverityInfo
	}
}
func (s Severity) String() string {
	switch s {
	case SeverityCritical:
		return "critical"
	case SeverityWarning:
		return "warning"
	default:
		return "info"
	}
}
