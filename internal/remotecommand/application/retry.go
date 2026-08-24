package application

import (
	"context"
	"fmt"
	"github.com/example/industrial-device-platform/internal/remotecommand/domain"
	"time"
)

type Executor interface {
	Execute(context.Context, *domain.Command) (string, error)
}

func (s *Service) Dispatch(ctx context.Context, id string, exec Executor) (*domain.Command, error) {
	c, e := s.Get(ctx, id)
	if e != nil {
		return nil, e
	}
	if e = c.Transition(domain.StatusDispatched); e != nil {
		return nil, e
	}
	if e = s.repo.Update(ctx, c); e != nil {
		return nil, e
	}
	result, e := exec.Execute(ctx, c)
	if e != nil {
		c.Status = domain.StatusFailed
		c.Result = e.Error()
	} else {
		c.Status = domain.StatusSucceeded
		c.Result = result
	}
	c.UpdatedAt = time.Now().UTC()
	if up := s.repo.Update(ctx, c); up != nil {
		return c, combineDispatchError(e, up)
	}
	return c, e
}

func combineDispatchError(primary, persist error) error {
	if primary == nil {
		return fmt.Errorf("persist dispatch: %w", persist)
	}
	return fmt.Errorf("dispatch failed: %v (persist dispatch: %w)", primary, persist)
}
