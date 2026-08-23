package score005_test

import (
	"context"
	firmwareinfra "github.com/example/industrial-device-platform/internal/firmware/infrastructure"
	"testing"

	"github.com/example/industrial-device-platform/internal/firmware/domain"
)

func TestFirmwareCanaryDoesNotAliasTargets(t *testing.T) {
	r := firmwareinfra.NewRepository()
	want := &domain.Release{ID: "rel-1", Version: "1.2.3", DeviceIDs: []string{"a", "b", "c", "d"}, Status: "draft"}
	if err := r.Save(context.Background(), want); err != nil {
		t.Fatal(err)
	}
	got, err := r.Get(context.Background(), want.ID)
	if err != nil {
		t.Fatal(err)
	}
	got.DeviceIDs[0] = "mutated-read"
	untouched, _ := r.Get(context.Background(), want.ID)
	if untouched.DeviceIDs[0] != "a" {
		t.Fatal("repository get leaked its backing slice")
	}
	canary := domain.SelectCanaryInto(got.DeviceIDs, 50)
	canary[0] = "changed"
	canary = append(canary, "extra")
	again, _ := r.Get(context.Background(), want.ID)
	if again.DeviceIDs[0] != "a" || len(again.DeviceIDs) != 4 {
		t.Fatalf("target list was aliased: %#v", again.DeviceIDs)
	}
	bytes := []byte("firmware")
	copyBytes := domain.CloneBytes(bytes)
	copyBytes[0] = 'F'
	if string(bytes) != "firmware" {
		t.Fatal("byte snapshot failed")
	}
	clone := again.Clone()
	clone.DeviceIDs[0] = "clone-only"
	if again.DeviceIDs[0] != "a" || clone.TargetCount() != 4 {
		t.Fatal("release clone lost targets")
	}
	source := []string{"x", "y", "z", "w"}
	selected := domain.SelectCanary(source, 50)
	selected[0] = "selected-only"
	if source[0] != "x" {
		t.Fatal("canary selection reused source backing array")
	}
	releaseClone := want.Clone()
	releaseClone.DeviceIDs[0] = "release-only"
	if want.DeviceIDs[0] != "a" {
		t.Fatal("release clone leaked device ids")
	}
}
