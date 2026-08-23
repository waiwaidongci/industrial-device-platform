package domain

import "fmt"

func (r Release) ValidateCanary() error {
	if r.CanaryPercent < 0 || r.CanaryPercent > 100 {
		return fmt.Errorf("canary percent must be between 0 and 100")
	}
	if r.Version == "" {
		return fmt.Errorf("version required")
	}
	return nil
}
func (r *Release) Start() error {
	if r.Status != "draft" {
		return fmt.Errorf("release is not draft")
	}
	r.Status = "running"
	return nil
}
func (r *Release) Pause() error {
	if r.Status != "running" {
		return fmt.Errorf("release is not running")
	}
	r.Status = "paused"
	return nil
}
func (r *Release) Resume() error {
	if r.Status != "paused" {
		return fmt.Errorf("release is not paused")
	}
	r.Status = "running"
	return nil
}
func (r *Release) Complete() { r.Status = "completed" }
