package service

import (
	"context"
	"testing"

	"github.com/6547709/goct/pkg/adapter"
)

// fakeSnapshotClient 同时实现 VMOps 和 SnapshotOps。
type fakeSnapshotClient struct {
	*fakeVMOps
	snaps []adapter.Snapshot
}

func (f *fakeSnapshotClient) ListSnapshots(_ context.Context, vmID string) ([]adapter.Snapshot, error) {
	var out []adapter.Snapshot
	for _, s := range f.snaps {
		if s.VMID == vmID {
			out = append(out, s)
		}
	}
	return out, nil
}

func (f *fakeSnapshotClient) CreateSnapshot(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-snap-create"}, nil
}

func (f *fakeSnapshotClient) RevertSnapshot(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-snap-revert"}, nil
}

func (f *fakeSnapshotClient) DeleteSnapshot(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-snap-delete"}, nil
}

func newFakeSnapshotClient() *fakeSnapshotClient {
	return &fakeSnapshotClient{
		fakeVMOps: newFakeVMOps(),
		snaps: []adapter.Snapshot{
			{ID: "snap-1", Name: "daily-backup", VMID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"},
			{ID: "snap-2", Name: "before-upgrade", VMID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"},
		},
	}
}

func TestSnapshotService_List(t *testing.T) {
	svc := NewSnapshot(newFakeSnapshotClient())
	snaps, err := svc.List(context.Background(), "test-vm-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(snaps) != 2 {
		t.Fatalf("want 2, got %d", len(snaps))
	}
}

func TestSnapshotService_Create(t *testing.T) {
	svc := NewSnapshot(newFakeSnapshotClient())
	ref, err := svc.Create(context.Background(), "test-vm-1", "new-snap")
	if err != nil {
		t.Fatal(err)
	}
	if ref.ID != "task-snap-create" {
		t.Fatalf("want task-snap-create, got %s", ref.ID)
	}
}

func TestSnapshotService_Delete(t *testing.T) {
	svc := NewSnapshot(newFakeSnapshotClient())
	ref, err := svc.Delete(context.Background(), "snap-1")
	if err != nil {
		t.Fatal(err)
	}
	if ref.ID != "task-snap-delete" {
		t.Fatalf("want task-snap-delete, got %s", ref.ID)
	}
}
