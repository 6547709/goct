package service

import (
	"context"
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
)

// VMService 封装 VM 相关的业务逻辑。
type VMService struct{ c adapter.Client }

// NewVM 构造 VMService。c 一般来自 client.From(ctx)（实现 adapter.Client）。
func NewVM(c adapter.Client) *VMService { return &VMService{c: c} }

// List 列出虚拟机。
func (s *VMService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.VM, error) {
	return s.c.ListVMs(ctx, opts)
}

// Resolve 按 name 或 UUID 解析到单个 VM。
func (s *VMService) Resolve(ctx context.Context, idOrName string) (*adapter.VM, error) {
	return Resolve(ctx, s.c.ListVMs, s.c.GetVM,
		func(v adapter.VM) (string, string) { return v.ID, v.Name },
		idOrName)
}

// Create 创建 VM，返回 task 引用。
func (s *VMService) Create(ctx context.Context, spec adapter.VMCreateSpec) (adapter.TaskRef, error) {
	return s.c.CreateVM(ctx, spec)
}

// Clone 克隆 VM。srcIDOrName 先 Resolve 再调 adapter。
func (s *VMService) Clone(ctx context.Context, srcIDOrName string, spec adapter.VMCloneSpec) (adapter.TaskRef, error) {
	src, err := s.Resolve(ctx, srcIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.CloneVM(ctx, src.ID, spec)
}

// Destroy 销毁 VM。
func (s *VMService) Destroy(ctx context.Context, idOrName string, force bool) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.DestroyVM(ctx, v.ID, force)
}

// Migrate 迁移 VM 到另一台主机。
// 如果 hostID 为空，则在同集群内随机选择一台主机（不能是当前主机）。
// hostID 可以是 host ID 或 host 名称。
func (s *VMService) Migrate(ctx context.Context, vmIDOrName, hostIDOrName string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}

	// 解析 host 名称或 ID
	targetHostID, err := s.resolveHostID(ctx, hostIDOrName, v.ClusterID, v.HostID)
	if err != nil {
		return adapter.TaskRef{}, err
	}

	return s.c.MigrateVM(ctx, v.ID, targetHostID)
}

// resolveHostID 解析 host 名称或 ID。
//
// 行为约定（v0.2.1 修复）：
//   - nameOrID 非空 → 解析名称到 ID（IsID 直接透传）。
//   - nameOrID 为空 → **不再客户端随机选 host**：govc / vCenter 的语义是让调度器自动决策，
//     CloudTower 也应当走集群侧调度（cluster scheduler）。当前 SDK 的 MigrateVM 接受空 hostID 时
//     由 CloudTower 自行选择目标主机。
//
// excludeHostID 仅在用户给了名字、并且名字解析后等于源主机时用于校验阻断（避免迁移到自己）。
func (s *VMService) resolveHostID(ctx context.Context, nameOrID string, clusterID, excludeHostID string) (string, error) {
	_ = clusterID // 保留参数以便后续做"集群范围内校验"
	if nameOrID == "" {
		// 空 hostID → 让 CloudTower 调度选目标，符合 govc 的行为预期。
		return "", nil
	}
	if IsID(nameOrID) {
		if nameOrID == excludeHostID {
			return "", fmt.Errorf("target host %s is the current host of the VM", nameOrID)
		}
		return nameOrID, nil
	}
	host, err := s.c.GetHostByName(ctx, nameOrID)
	if err != nil {
		return "", fmt.Errorf("resolve host %q: %w", nameOrID, err)
	}
	if host.ID == excludeHostID {
		return "", fmt.Errorf("target host %q resolves to the current host of the VM", nameOrID)
	}
	return host.ID, nil
}

// Export 导出 VM。
func (s *VMService) Export(ctx context.Context, idOrName string, spec adapter.VMExportSpec) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.ExportVM(ctx, v.ID, spec)
}

// Power 执行 VM 电源操作。
func (s *VMService) Power(ctx context.Context, idOrName string, action adapter.PowerAction, force bool) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.PowerVM(ctx, v.ID, action, force)
}

// Update 更新 VM 基本信息。
func (s *VMService) Update(ctx context.Context, idOrName string, spec adapter.VMUpdateSpec) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.UpdateVM(ctx, v.ID, spec)
}

// MoveToRecycle 把 VM 放入回收站。
func (s *VMService) MoveToRecycle(ctx context.Context, idOrName string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.MoveToRecycle(ctx, v.ID)
}

// RecoverFromRecycle 从回收站恢复 VM。
func (s *VMService) RecoverFromRecycle(ctx context.Context, idOrName string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.RecoverFromRecycle(ctx, v.ID)
}

// ShutDown 优雅关闭 VM。
func (s *VMService) ShutDown(ctx context.Context, idOrName string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, idOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.ShutDownVM(ctx, v.ID)
}

// AddDisk 添加磁盘到 VM。
func (s *VMService) AddDisk(ctx context.Context, vmIDOrName string, spec adapter.DiskAddSpec) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.AddDisk(ctx, v.ID, spec)
}

// ExpandDisk 扩容 VM 磁盘。
func (s *VMService) ExpandDisk(ctx context.Context, vmIDOrName, diskID string, sizeBytes int64) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.ExpandDisk(ctx, v.ID, diskID, sizeBytes)
}

// RemoveDisk 从 VM 移除磁盘。
func (s *VMService) RemoveDisk(ctx context.Context, vmIDOrName, diskID string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.RemoveDisk(ctx, v.ID, diskID)
}

// AddCdRom 添加 CD-ROM 到 VM。
func (s *VMService) AddCdRom(ctx context.Context, vmIDOrName, isoPath string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.AddCdRom(ctx, v.ID, isoPath)
}

// EjectCdRom 弹出 CD-ROM ISO。
func (s *VMService) EjectCdRom(ctx context.Context, cdromID string) (adapter.TaskRef, error) {
	return s.c.EjectCdRom(ctx, cdromID)
}

// RemoveCdRom 从 VM 移除 CD-ROM。
func (s *VMService) RemoveCdRom(ctx context.Context, vmIDOrName, cdromID string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.RemoveCdRom(ctx, v.ID, cdromID)
}

// AddNic 添加 NIC 到 VM。
func (s *VMService) AddNic(ctx context.Context, vmIDOrName string, spec adapter.NicAddSpec) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.AddNic(ctx, v.ID, spec)
}

// RemoveNic 从 VM 移除 NIC（按索引）。
func (s *VMService) RemoveNic(ctx context.Context, vmIDOrName string, nicIndex int32) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.RemoveNic(ctx, v.ID, nicIndex)
}

// AddGpuDevice 添加 GPU 设备到 VM。
func (s *VMService) AddGpuDevice(ctx context.Context, vmIDOrName, gpuDeviceID string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.AddGpuDevice(ctx, v.ID, gpuDeviceID)
}

// RemoveGpuDevice 从 VM 移除 GPU 设备。
func (s *VMService) RemoveGpuDevice(ctx context.Context, vmIDOrName, gpuDeviceID string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.RemoveGpuDevice(ctx, v.ID, gpuDeviceID)
}

// InstallVmtools 在 VM 上安装 VMware Tools。
func (s *VMService) InstallVmtools(ctx context.Context, vmIDOrName string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	if v.VMToolsStatus != "" && v.VMToolsStatus != "NOT_INSTALLED" {
		return adapter.TaskRef{ID: "already-installed"}, nil
	}
	return s.c.InstallVmtools(ctx, v.ID)
}

// GetVNCInfo 获取 VM 的 VNC 连接信息。
func (s *VMService) GetVNCInfo(ctx context.Context, vmIDOrName string) (*adapter.VNCInfo, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return nil, err
	}
	return s.c.GetVNCInfo(ctx, v.ID)
}

// MigrateAcrossCluster 跨集群迁移 VM。
// targetClusterID 可以是集群 ID 或名称。
// hostID 可以是 host ID、host 名称，或空（由 CloudTower 自动选择）。
func (s *VMService) MigrateAcrossCluster(ctx context.Context, vmIDOrName, targetClusterIDOrName, hostIDOrName string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}

	// 解析集群名称到 ID
	clusterID, err := s.resolveClusterID(ctx, targetClusterIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}

	// 解析 host 名称到 ID（可选）
	var hostID string
	if hostIDOrName != "" {
		hostID, err = s.resolveHostID(ctx, hostIDOrName, clusterID, "")
		if err != nil {
			return adapter.TaskRef{}, err
		}
	}

	return s.c.MigrateAcrossCluster(ctx, v.ID, clusterID, hostID)
}

// CreateFromTemplate 从模板创建 VM。
func (s *VMService) CreateFromTemplate(ctx context.Context, spec adapter.VMCreateFromTemplateSpec) (adapter.TaskRef, error) {
	// Resolve cluster name to ID
	clusterID, err := s.resolveClusterID(ctx, spec.ClusterID)
	if err != nil {
		return adapter.TaskRef{}, err
	}

	// Resolve template name to ID
	templateID, err := s.resolveTemplateID(ctx, spec.TemplateID)
	if err != nil {
		return adapter.TaskRef{}, err
	}

	// Resolve VLAN name to ID if specified
	if spec.NIC.VlanID != "" {
		vlanID, err := s.resolveVLANID(ctx, spec.NIC.VlanID)
		if err != nil {
			return adapter.TaskRef{}, err
		}
		spec.NIC.VlanID = vlanID
	}

	spec.ClusterID = clusterID
	spec.TemplateID = templateID

	return s.c.CreateVMFromTemplate(ctx, spec)
}

// resolveClusterID resolves cluster name to ID. If already an ID, returns it as-is.
func (s *VMService) resolveClusterID(ctx context.Context, nameOrID string) (string, error) {
	if IsID(nameOrID) {
		return nameOrID, nil
	}
	cluster, err := s.c.GetClusterByName(ctx, nameOrID)
	if err != nil {
		return "", fmt.Errorf("resolve cluster %q: %w", nameOrID, err)
	}
	return cluster.ID, nil
}

// resolveTemplateID resolves template name to ID. If already an ID, returns it as-is.
func (s *VMService) resolveTemplateID(ctx context.Context, nameOrID string) (string, error) {
	if IsID(nameOrID) {
		return nameOrID, nil
	}
	tpl, err := s.c.GetContentLibraryTemplateByName(ctx, nameOrID)
	if err != nil {
		return "", fmt.Errorf("resolve template %q: %w", nameOrID, err)
	}
	return tpl.ID, nil
}

// resolveVLANID resolves VLAN name to ID. If already an ID, returns it as-is.
func (s *VMService) resolveVLANID(ctx context.Context, nameOrID string) (string, error) {
	if IsID(nameOrID) {
		return nameOrID, nil
	}
	vlan, err := s.c.GetVLANByName(ctx, nameOrID)
	if err != nil {
		return "", fmt.Errorf("resolve VLAN %q: %w", nameOrID, err)
	}
	return vlan.ID, nil
}

// ListNics lists all NICs of a VM.
func (s *VMService) ListNics(ctx context.Context, vmIDOrName string) ([]adapter.VMNic, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return nil, err
	}
	return s.c.ListVMNics(ctx, v.ID)
}

// UpdateNic updates a NIC configuration.
// 必须指定 vmIDOrName，因 CloudTower update-vm-nic API 的 Where 是 VM 维度。
func (s *VMService) UpdateNic(ctx context.Context, vmIDOrName, nicID string, spec adapter.VMNicUpdateSpec) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	spec.VMID = v.ID
	return s.c.UpdateNic(ctx, nicID, spec)
}

// ListDisks lists all disks of a VM.
func (s *VMService) ListDisks(ctx context.Context, vmIDOrName string) ([]adapter.VMDisk, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return nil, err
	}
	return s.c.ListVMDisks(ctx, v.ID)
}

// UpdateDisk updates disk configuration.
//
// 必须先 Resolve VM，然后把 VMID 透传给 adapter（CloudTower update-vm-disk
// 的 Where 是 VM 维度，不是 disk 维度）。
func (s *VMService) UpdateDisk(ctx context.Context, vmIDOrName, diskID string, spec adapter.DiskUpdateSpec) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	spec.VMID = v.ID
	return s.c.UpdateDisk(ctx, diskID, spec)
}

// ToggleCdRom enables or disables a CD-ROM.
func (s *VMService) ToggleCdRom(ctx context.Context, cdromID string, disabled bool) (adapter.TaskRef, error) {
	return s.c.ToggleCdRom(ctx, cdromID, adapter.CdRomToggleSpec{Disabled: disabled})
}

// ResetPassword resets the guest OS password for a VM.
func (s *VMService) ResetPassword(ctx context.Context, vmIDOrName, username, password string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.ResetPassword(ctx, v.ID, adapter.ResetPasswordSpec{Username: username, Password: password})
}

// RebuildVM rebuilds a VM from a snapshot.
func (s *VMService) RebuildVM(ctx context.Context, vmIDOrName string, spec adapter.RebuildVMSpec) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.RebuildVM(ctx, v.ID, spec)
}

// AbortMigrateAcrossCluster aborts a cross-cluster migration in progress.
// 先把 name 解析成 VM ID 再传给 adapter，避免 adapter 拿到一个 name 当成 resource_id 过滤为空。
func (s *VMService) AbortMigrateAcrossCluster(ctx context.Context, vmIDOrName string) (adapter.TaskRef, error) {
	v, err := s.Resolve(ctx, vmIDOrName)
	if err != nil {
		return adapter.TaskRef{}, err
	}
	return s.c.AbortMigrateAcrossCluster(ctx, v.ID)
}

// ConvertToVM converts a template to a VM.
//
// 如果 newName 为空，则回退到 "<template-name>-vm"，避免下发空字符串导致 CloudTower 422。
func (s *VMService) ConvertToVM(ctx context.Context, templateIDOrName, newName string) (adapter.TaskRef, error) {
	tpl, err := s.c.GetContentLibraryTemplateByName(ctx, templateIDOrName)
	if err != nil && !IsID(templateIDOrName) {
		return adapter.TaskRef{}, fmt.Errorf("resolve template %q: %w", templateIDOrName, err)
	}

	templateID := templateIDOrName
	tplName := ""
	if tpl != nil {
		templateID = tpl.ID
		tplName = tpl.Name
	}
	if newName == "" {
		if tplName != "" {
			newName = tplName + "-vm"
		} else {
			newName = "converted-vm"
		}
	}
	return s.c.ConvertToVM(ctx, templateID, newName)
}
