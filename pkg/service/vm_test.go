package service

import (
	"context"
	"errors"
	"testing"

	"github.com/6547709/goct/pkg/adapter"
)

// fakeClient 是 adapter.Client 的 test double。
type fakeClient struct {
	vms []adapter.VM
}

func (f *fakeClient) About(_ context.Context) (adapter.TowerInfo, error) {
	return adapter.TowerInfo{Version: "test"}, nil
}

func (f *fakeClient) ListVMs(_ context.Context, opts adapter.ListOpts) ([]adapter.VM, error) {
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

func (f *fakeClient) GetVM(_ context.Context, id string) (*adapter.VM, error) {
	for _, v := range f.vms {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, adapter.ErrNotFound
}

func (f *fakeClient) CreateVM(_ context.Context, _ adapter.VMCreateSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-create"}, nil
}

func (f *fakeClient) CloneVM(_ context.Context, _ string, _ adapter.VMCloneSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-clone"}, nil
}

func (f *fakeClient) DestroyVM(_ context.Context, _ string, _ bool) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-destroy"}, nil
}

func (f *fakeClient) MigrateVM(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-migrate"}, nil
}

func (f *fakeClient) ExportVM(_ context.Context, _ string, _ adapter.VMExportSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-export"}, nil
}

func (f *fakeClient) PowerVM(_ context.Context, _ string, _ adapter.PowerAction, _ bool) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-power"}, nil
}

func (f *fakeClient) UpdateVM(_ context.Context, _ string, _ adapter.VMUpdateSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-update"}, nil
}

func (f *fakeClient) MoveToRecycle(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-recycle"}, nil
}

func (f *fakeClient) RecoverFromRecycle(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-recover"}, nil
}

func (f *fakeClient) ShutDownVM(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-shutdown"}, nil
}

func (f *fakeClient) AddDisk(_ context.Context, _ string, _ adapter.DiskAddSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-disk-add"}, nil
}

func (f *fakeClient) ExpandDisk(_ context.Context, _, _ string, _ int64) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-disk-expand"}, nil
}

func (f *fakeClient) RemoveDisk(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-disk-rm"}, nil
}

func (f *fakeClient) AddCdRom(_ context.Context, _ string, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-cdrom-add"}, nil
}

func (f *fakeClient) EjectCdRom(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-cdrom-eject"}, nil
}

func (f *fakeClient) RemoveCdRom(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-cdrom-rm"}, nil
}

func (f *fakeClient) AddNic(_ context.Context, _ string, _ adapter.NicAddSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-nic-add"}, nil
}

func (f *fakeClient) RemoveNic(_ context.Context, _ string, _ int32) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-nic-rm"}, nil
}

func (f *fakeClient) AddGpuDevice(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-gpu-add"}, nil
}

func (f *fakeClient) RemoveGpuDevice(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-gpu-rm"}, nil
}

func (f *fakeClient) InstallVmtools(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-vmtools"}, nil
}

func (f *fakeClient) GetVNCInfo(_ context.Context, _ string) (*adapter.VNCInfo, error) {
	return &adapter.VNCInfo{ClusterIP: "192.168.1.100", Redirect: "localhost:5900"}, nil
}

func (f *fakeClient) MigrateAcrossCluster(_ context.Context, _, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-migrate-across"}, nil
}

func (f *fakeClient) CreateVMFromTemplate(_ context.Context, _ adapter.VMCreateFromTemplateSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-create-from-template"}, nil
}

func (f *fakeClient) ListVMNics(_ context.Context, _ string) ([]adapter.VMNic, error) {
	return nil, nil
}

func (f *fakeClient) UpdateNic(_ context.Context, _ string, _ adapter.VMNicUpdateSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-nic-update"}, nil
}

func (f *fakeClient) ListVMDisks(_ context.Context, _ string) ([]adapter.VMDisk, error) {
	return nil, nil
}

func (f *fakeClient) UpdateDisk(_ context.Context, _ string, _ adapter.DiskUpdateSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-disk-update"}, nil
}

func (f *fakeClient) ToggleCdRom(_ context.Context, _ string, _ adapter.CdRomToggleSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-cdrom-toggle"}, nil
}

func (f *fakeClient) ResetPassword(_ context.Context, _ string, _ adapter.ResetPasswordSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-reset-password"}, nil
}

func (f *fakeClient) RebuildVM(_ context.Context, _ string, _ adapter.RebuildVMSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-rebuild"}, nil
}

func (f *fakeClient) AbortMigrateAcrossCluster(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-abort-migrate"}, nil
}

func (f *fakeClient) ConvertToVM(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-convert-to-vm"}, nil
}

func (f *fakeClient) ListSnapshots(_ context.Context, _ string) ([]adapter.Snapshot, error) {
	return nil, nil
}
func (f *fakeClient) GetSnapshot(_ context.Context, _ string) (*adapter.Snapshot, error) {
	return nil, nil
}
func (f *fakeClient) CreateSnapshot(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-snapshot-create"}, nil
}
func (f *fakeClient) DeleteSnapshot(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-snapshot-delete"}, nil
}
func (f *fakeClient) RevertSnapshot(_ context.Context, _, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-snapshot-revert"}, nil
}
func (f *fakeClient) ListHosts(_ context.Context, _ adapter.ListOpts) ([]adapter.Host, error) {
	return nil, nil
}
func (f *fakeClient) GetHost(_ context.Context, _ string) (*adapter.Host, error) {
	return nil, nil
}
func (f *fakeClient) GetHostByName(_ context.Context, _ string) (*adapter.Host, error) {
	return nil, nil
}
func (f *fakeClient) ListHostsByCluster(_ context.Context, _ string) ([]adapter.Host, error) {
	return nil, nil
}
func (f *fakeClient) EnterMaintenanceMode(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-host-enter-maint"}, nil
}
func (f *fakeClient) ExitMaintenanceMode(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-host-exit-maint"}, nil
}
func (f *fakeClient) ShutdownHost(_ context.Context, _ string, _ bool) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-host-shutdown"}, nil
}
func (f *fakeClient) RebootHost(_ context.Context, _ string, _ bool) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-host-reboot"}, nil
}
func (f *fakeClient) ListClusters(_ context.Context, _ adapter.ListOpts) ([]adapter.Cluster, error) {
	return nil, nil
}
func (f *fakeClient) GetCluster(_ context.Context, _ string) (*adapter.Cluster, error) {
	return nil, nil
}
func (f *fakeClient) GetClusterByName(_ context.Context, _ string) (*adapter.Cluster, error) {
	return nil, nil
}
func (f *fakeClient) ListDatastores(_ context.Context, _ adapter.ListOpts) ([]adapter.Datastore, error) {
	return nil, nil
}
func (f *fakeClient) GetDatastore(_ context.Context, _ string) (*adapter.Datastore, error) {
	return nil, nil
}
func (f *fakeClient) ListDisks(_ context.Context, _ string) ([]adapter.Disk, error) {
	return nil, nil
}
func (f *fakeClient) ListDiskPools(_ context.Context, _ adapter.ListOpts) ([]adapter.DiskPool, error) {
	return nil, nil
}
func (f *fakeClient) ListNetworks(_ context.Context, _ adapter.ListOpts) ([]adapter.Network, error) {
	return nil, nil
}
func (f *fakeClient) GetNetwork(_ context.Context, _ string) (*adapter.Network, error) {
	return nil, nil
}
func (f *fakeClient) ListVLANs(_ context.Context, _ adapter.ListOpts) ([]adapter.VLAN, error) {
	return nil, nil
}
func (f *fakeClient) GetVLAN(_ context.Context, _ string) (*adapter.VLAN, error) {
	return nil, nil
}
func (f *fakeClient) GetVLANByName(_ context.Context, _ string) (*adapter.VLAN, error) {
	return nil, nil
}
func (f *fakeClient) CreateVLAN(_ context.Context, _ adapter.VLANCreateSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-vlan-create"}, nil
}
func (f *fakeClient) DeleteVLAN(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-vlan-delete"}, nil
}
func (f *fakeClient) ListTasks(_ context.Context, _ adapter.ListOpts) ([]adapter.Task, error) {
	return nil, nil
}
func (f *fakeClient) GetTask(_ context.Context, _ string) (*adapter.Task, error) {
	return nil, nil
}
func (f *fakeClient) ListAlerts(_ context.Context, _ adapter.ListOpts) ([]adapter.Alert, error) {
	return nil, nil
}
func (f *fakeClient) GetAlert(_ context.Context, _ string) (*adapter.Alert, error) {
	return nil, nil
}
func (f *fakeClient) ListUsers(_ context.Context, _ adapter.ListOpts) ([]adapter.User, error) {
	return nil, nil
}
func (f *fakeClient) GetUser(_ context.Context, _ string) (*adapter.User, error) {
	return nil, nil
}
func (f *fakeClient) CreateUser(_ context.Context, _ adapter.UserCreateSpec) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-user-create"}, nil
}
func (f *fakeClient) DeleteUser(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-user-delete"}, nil
}
func (f *fakeClient) ListContentLibraryTemplates(_ context.Context, _ adapter.ListOpts) ([]adapter.ContentLibraryTemplate, error) {
	return nil, nil
}
func (f *fakeClient) GetContentLibraryTemplateByName(_ context.Context, _ string) (*adapter.ContentLibraryTemplate, error) {
	return nil, nil
}
func (f *fakeClient) GetTaskProgress(_ context.Context, _ string) (int, string, error) {
	return 0, "", nil
}

func (f *fakeClient) AckAlert(_ context.Context, _ string) (adapter.TaskRef, error) {
	return adapter.TaskRef{ID: "task-alert-ack"}, nil
}

func newFakeClient() *fakeClient {
	return &fakeClient{
		vms: []adapter.VM{
			{ID: "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", Name: "test-vm-1", Status: "RUNNING"},
			{ID: "11111111-2222-3333-4444-555555555555", Name: "test-vm-2", Status: "STOPPED"},
		},
	}
}

func TestVMService_List(t *testing.T) {
	svc := NewVM(newFakeClient())
	vms, err := svc.List(context.Background(), adapter.ListOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(vms) != 2 {
		t.Fatalf("want 2, got %d", len(vms))
	}
}

func TestVMService_Resolve_ByUUID(t *testing.T) {
	svc := NewVM(newFakeClient())
	v, err := svc.Resolve(context.Background(), "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	if err != nil {
		t.Fatal(err)
	}
	if v.Name != "test-vm-1" {
		t.Fatalf("want test-vm-1, got %s", v.Name)
	}
}

func TestVMService_Resolve_ByName(t *testing.T) {
	svc := NewVM(newFakeClient())
	v, err := svc.Resolve(context.Background(), "test-vm-2")
	if err != nil {
		t.Fatal(err)
	}
	if v.ID != "11111111-2222-3333-4444-555555555555" {
		t.Fatalf("want 11111111..., got %s", v.ID)
	}
}

func TestVMService_Resolve_NotFound(t *testing.T) {
	svc := NewVM(newFakeClient())
	_, err := svc.Resolve(context.Background(), "nonexistent")
	if !errors.Is(err, adapter.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestVMService_Power(t *testing.T) {
	svc := NewVM(newFakeClient())
	ref, err := svc.Power(context.Background(), "test-vm-1", adapter.PowerOn, false)
	if err != nil {
		t.Fatal(err)
	}
	if ref.ID != "task-power" {
		t.Fatalf("want task-power, got %s", ref.ID)
	}
}

func TestVMService_Clone(t *testing.T) {
	svc := NewVM(newFakeClient())
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

func TestResolver_IsID(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		// UUID
		{"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee", true},
		{"5e52cf6e-1e8c-4a0a-9e3a-1b2c3d4e5f6a", true},
		// cuid (CloudTower style)
		{"cl5k7g2xo04070822fhxjfsev9q", true},
		{"cl0000000000000000000000000", true},
		// Not IDs
		{"my-vm-name", false},
		{"web-server-01", false},
		{"", false},
		{"cl", false},       // too short
		{"CL5K7G2XO04070822FHXJFSEV9Q", false}, // cuid is lowercase
	}
	for _, tt := range tests {
		if got := IsID(tt.in); got != tt.want {
			t.Errorf("IsID(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}
