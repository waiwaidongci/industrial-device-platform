package score007_test

import (
	"context"
	commandapp "github.com/example/industrial-device-platform/internal/remotecommand/application"
	"testing"

	"github.com/example/industrial-device-platform/internal/remotecommand/domain"
	"github.com/example/industrial-device-platform/internal/remotecommand/infrastructure"
)

func TestCommandTerminalTransition(t *testing.T) {
	repo := infrastructure.NewMemoryRepository()
	svc := commandapp.NewService(repo)
	c, err := svc.Submit(context.Background(), &domain.Command{DeviceID: "d1", Type: "stop"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Complete(context.Background(), c.ID, domain.StatusDispatched, "sent"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Complete(context.Background(), c.ID, domain.StatusSucceeded, "ok"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Complete(context.Background(), c.ID, domain.StatusFailed, "late"); err == nil {
		t.Fatal("terminal command accepted a late transition")
	}
	got, _ := svc.Get(context.Background(), c.ID)
	if got.Status != domain.StatusSucceeded || len(got.History) != 3 {
		t.Fatalf("bad state history: %#v", got)
	}
	active, err := svc.ActiveRecent(context.Background(), "d1", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(active) != 0 {
		t.Fatalf("terminal command leaked into active list: %#v", active)
	}
}
