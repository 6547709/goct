package service

import (
	"context"
	"errors"
	"testing"

	"github.com/6547709/goct/pkg/adapter"
)

// fakeVMOps 是 adapter.VMOps 的 test double。
type fakeVMOps struct {
	vms []adapter.VM
}

func (f *fakeVMOps) ListVMs(_ context.Context, opts adapter.ListOpts) ([]adapter.VM, error) {
	if opts.NameContains == "" {
		return f.vms, nil
	}
	var out []adapter.VM
	for _, v := range f.vms {
		if v.Name == opts.NameContains {
			out = append(out, v)
		}
	}
	return out, nil
}

func (f *fakeVMOps) GetVM(_ context.Context, id string) (*adapter.VM, error) {
	for _, v := range f.vms {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, adapter.ErrNotFound
}

func (f *fakeVMOps) CreateVM(_ context.Context, _ adapter.VMCreateSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-create"}, nil
}

func (f *fakeVMOps) CloneVM(_ context.Context, _ string, _ adapter.VMCloneSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-clone"}, nil
}

func (f *fakeVMOps) DestroyVM(_ context.Context, _ string, _ bool) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-destroy"}, nil
}

func (f *fakeVMOps) MigrateVM(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-migrate"}, nil
}

func (f *fakeVMOps) ExportVM(_ context.Context, _ string, _ adapter.VMExportSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-export"}, nil
}

func (f *fakeVMOps) PowerVM(_ context.Context, _ string, _ adapter.PowerAction, _ bool) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-power"}, nil
}

func newFakeVMOps() *fakeVMOps {
	return &fakeVMOps{
		vms: []adapter.VM{
			{ID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", Name: "test-vm-1", Status: "RUNNING"},
			{ID: "11111111-2222-3333-4444-555555555555", Name: "test-vm-2", Status: "STOPPED"},
		},
	}
}

func TestVMService_List(t *testing.T) {
	svc := NewVM(newFakeVMOps())
	vms, err := svc.List(context.Background(), adapter.ListOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(vms) != 2 {
		t.Fatalf("want 2, got %d", len(vms))
	}
}

func TestVMService_Resolve_ByUUID(t *testing.T) {
	svc := NewVM(newFakeVMOps())
	v, err := svc.Resolve(context.Background(), "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	if err != nil {
		t.Fatal(err)
	}
	if v.Name != "test-vm-1" {
		t.Fatalf("want test-vm-1, got %s", v.Name)
	}
}

func TestVMService_Resolve_ByName(t *testing.T) {
	svc := NewVM(newFakeVMOps())
	v, err := svc.Resolve(context.Background(), "test-vm-2")
	if err != nil {
		t.Fatal(err)
	}
	if v.ID != "11111111-2222-3333-4444-555555555555" {
		t.Fatalf("want 11111111..., got %s", v.ID)
	}
}

func TestVMService_Resolve_NotFound(t *testing.T) {
	svc := NewVM(newFakeVMOps())
	_, err := svc.Resolve(context.Background(), "nonexistent")
	if !errors.Is(err, adapter.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestVMService_Power(t *testing.T) {
	svc := NewVM(newFakeVMOps())
	ref, err := svc.Power(context.Background(), "test-vm-1", adapter.PowerOn, false)
	if err != nil {
		t.Fatal(err)
	}
	if ref.ID != "task-power" {
		t.Fatalf("want task-power, got %s", ref.ID)
	}
}

func TestVMService_Clone(t *testing.T) {
	svc := NewVM(newFakeVMOps())
	ref, err := svc.Clone(context.Background(), "test-vm-1", adapter.VMCloneSpec{Name: "clone-1"})
	if err != nil {
		t.Fatal(err)
	}
	if ref.ID != "task-clone" {
		t.Fatalf("want task-clone, got %s", ref.ID)
	}
}

func TestResolver_IsUUID(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", true},
		{"AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE", true},
		{"not-a-uuid", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsUUID(tt.in); got != tt.want {
			t.Errorf("IsUUID(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
