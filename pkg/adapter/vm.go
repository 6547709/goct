package adapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm_disk"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm_nic"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// VMOps 定义 VM 相关的 SDK 操作。
type VMOps interface {
	ListVMs(ctx context.Context, opts ListOpts) ([]VM, error)
	GetVM(ctx context.Context, id string) (*VM, error)
	CreateVM(ctx context.Context, spec VMCreateSpec) (TaskRef, error)
	CloneVM(ctx context.Context, srcID string, spec VMCloneSpec) (TaskRef, error)
	DestroyVM(ctx context.Context, id string, force bool) (TaskRef, error)
	MigrateVM(ctx context.Context, id, hostID string) (TaskRef, error)
	ExportVM(ctx context.Context, id string, spec VMExportSpec) (TaskRef, error)
	PowerVM(ctx context.Context, id string, action PowerAction, force bool) (TaskRef, error)
	UpdateVM(ctx context.Context, id string, spec VMUpdateSpec) (TaskRef, error)
	MoveToRecycle(ctx context.Context, id string) (TaskRef, error)
	RecoverFromRecycle(ctx context.Context, id string) (TaskRef, error)
	ShutDownVM(ctx context.Context, id string) (TaskRef, error)
	AddDisk(ctx context.Context, vmID string, spec DiskAddSpec) (TaskRef, error)
	ExpandDisk(ctx context.Context, vmID, diskID string, sizeBytes int64) (TaskRef, error)
	RemoveDisk(ctx context.Context, vmID, diskID string) (TaskRef, error)
	AddCdRom(ctx context.Context, vmID string, isoPath string) (TaskRef, error)
	EjectCdRom(ctx context.Context, cdromID string) (TaskRef, error)
	RemoveCdRom(ctx context.Context, vmID, cdromID string) (TaskRef, error)
	// NIC 操作
	AddNic(ctx context.Context, vmID string, spec NicAddSpec) (TaskRef, error)
	RemoveNic(ctx context.Context, vmID string, nicIndex int32) (TaskRef, error)
	ListVMNics(ctx context.Context, vmID string) ([]VMNic, error)
	UpdateNic(ctx context.Context, nicID string, spec VMNicUpdateSpec) (TaskRef, error)
	// GPU 操作
	AddGpuDevice(ctx context.Context, vmID, gpuDeviceID string) (TaskRef, error)
	RemoveGpuDevice(ctx context.Context, vmID, gpuDeviceID string) (TaskRef, error)
	// 磁盘操作
	ListVMDisks(ctx context.Context, vmID string) ([]VMDisk, error)
	UpdateDisk(ctx context.Context, diskID string, spec DiskUpdateSpec) (TaskRef, error)
	// CD-ROM 操作
	ToggleCdRom(ctx context.Context, cdromID string, spec CdRomToggleSpec) (TaskRef, error)
	// 其他
	InstallVmtools(ctx context.Context, vmID string) (TaskRef, error)
	GetVNCInfo(ctx context.Context, vmID string) (*VNCInfo, error)
	CreateVMFromTemplate(ctx context.Context, spec VMCreateFromTemplateSpec) (TaskRef, error)
	MigrateAcrossCluster(ctx context.Context, vmID, targetClusterID string, hostID string) (TaskRef, error)
	ResetPassword(ctx context.Context, vmID string, spec ResetPasswordSpec) (TaskRef, error)
	RebuildVM(ctx context.Context, vmID string, spec RebuildVMSpec) (TaskRef, error)
	AbortMigrateAcrossCluster(ctx context.Context, vmID string) (TaskRef, error)
	ConvertToVM(ctx context.Context, templateID string) (TaskRef, error)
}

// ---------- ListVMs ----------

func (c *defaultClient) ListVMs(ctx context.Context, opts ListOpts) ([]VM, error) {
	params := vm.NewGetVmsParams()
	params.SetContext(ctx)

	where := &models.VMWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if opts.ClusterID != "" {
		where.Cluster = &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}
		hasWhere = true
	}

	body := &models.GetVmsRequestBody{}
	if hasWhere {
		body.Where = where
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	if opts.Skip > 0 {
		body.Skip = pointy.Int32(opts.Skip)
	}
	params.SetRequestBody(body)

	resp, err := c.api.VM.GetVms(params)
	if err != nil {
		return nil, fmt.Errorf("list vms: %w", err)
	}
	out := make([]VM, 0, len(resp.Payload))
	for _, v := range resp.Payload {
		out = append(out, toVM(v))
	}
	return out, nil
}

// ---------- GetVM ----------

func (c *defaultClient) GetVM(ctx context.Context, id string) (*VM, error) {
	params := vm.NewGetVmsParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetVmsRequestBody{
		Where: &models.VMWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.VM.GetVms(params)
	if err != nil {
		return nil, fmt.Errorf("get vm %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get vm %s: %w", id, ErrNotFound)
	}
	v := toVM(resp.Payload[0])
	return &v, nil
}

// ---------- CreateVM ----------

func (c *defaultClient) CreateVM(ctx context.Context, spec VMCreateSpec) (TaskRef, error) {
	fw := models.VMFirmwareBIOS
	if strings.EqualFold(spec.Firmware, "UEFI") {
		fw = models.VMFirmwareUEFI
	}

	sockets := int32(1)
	cores := spec.VCPU
	if cores == 0 {
		cores = 1
	}

	// 默认创建一个 10GB SCSI 磁盘
	defaultDiskSize := int64(10 * 1024 * 1024 * 1024) // 10GB
	bus := models.BusSCSI

	p := &models.VMCreationParams{
		ClusterID:  pointy.String(spec.ClusterID),
		Name:       pointy.String(spec.Name),
		Ha:        pointy.Bool(true),
		CPUSockets: pointy.Int32(sockets),
		CPUCores:   pointy.Int32(cores),
		Memory:     pointy.Int64(spec.MemoryBytes),
		Firmware:   &fw,
		Status:     modelVMStatusStopped(),
		// 必须提供 vm_disks 和 vm_nics
		VMDisks: &models.VMDiskParams{
			MountNewCreateDisks: []*models.MountNewCreateDisksParams{
				{
					Boot: pointy.Int32(0),
					Bus:   &bus,
					Index: pointy.Int32(0),
					VMVolume: &models.MountNewCreateDisksParamsVMVolume{
						Name: pointy.String("disk0"),
						Size: pointy.Int64(defaultDiskSize),
					},
				},
			},
		},
		VMNics: []*models.VMNicParams{
			{
				Type:  models.VMNicTypeVLAN.Pointer(),
				Model: models.VMNicModelVIRTIO.Pointer(),
			},
		},
	}
	if spec.Description != "" {
		p.Description = pointy.String(spec.Description)
	}

	params := vm.NewCreateVMParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.VMCreationParams{p})

	resp, err := c.api.VM.CreateVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("create vm %s: %w", spec.Name, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- CloneVM ----------

func (c *defaultClient) CloneVM(ctx context.Context, srcID string, spec VMCloneSpec) (TaskRef, error) {
	p := &models.VMCloneParams{
		SrcVMID: pointy.String(srcID),
		Name:    pointy.String(spec.Name),
	}
	if spec.TargetClusterID != "" {
		p.ClusterID = pointy.String(spec.TargetClusterID)
	}

	params := vm.NewCloneVMParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.VMCloneParams{p})

	resp, err := c.api.VM.CloneVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("clone vm %s: %w", srcID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- DestroyVM ----------

func (c *defaultClient) DestroyVM(ctx context.Context, id string, _ bool) (TaskRef, error) {
	where := &models.VMWhereInput{ID: pointy.String(id)}
	body := &models.VMDeleteParams{Where: where}

	params := vm.NewDeleteVMParams()
	params.SetContext(ctx)
	params.SetRequestBody(body)
	resp, err := c.api.VM.DeleteVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("destroy vm %s: %w", id, err)
	}
	return firstDeleteVMTaskRef(resp.Payload), nil
}

// ---------- MigrateVM ----------

func (c *defaultClient) MigrateVM(ctx context.Context, id, hostID string) (TaskRef, error) {
	params := vm.NewMigrateVMParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMMigrateParams{
		Where: &models.VMWhereInput{ID: pointy.String(id)},
		Data: &models.VMMigrateParamsData{
			HostID: pointy.String(hostID),
		},
	})
	resp, err := c.api.VM.MigrateVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("migrate vm %s: %w", id, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- ExportVM ----------

func (c *defaultClient) ExportVM(ctx context.Context, id string, spec VMExportSpec) (TaskRef, error) {
	ft := models.VMExportFileTypeOVF
	params := vm.NewExportVMParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMExportParams{
		Where: &models.VMWhereInput{ID: pointy.String(id)},
		Data: &models.VMExportParamsData{
			Type:    &ft,
			KeepMac: pointy.Bool(spec.KeepMAC),
		},
	})
	resp, err := c.api.VM.ExportVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("export vm %s: %w", id, err)
	}
	return firstExportTaskRef(resp.Payload), nil
}

// ---------- PowerVM ----------

func (c *defaultClient) PowerVM(ctx context.Context, id string, action PowerAction, force bool) (TaskRef, error) {
	where := &models.VMWhereInput{ID: pointy.String(id)}
	switch action {
	case PowerOn:
		params := vm.NewStartVMParams()
		params.SetContext(ctx)
		params.SetRequestBody(&models.VMStartParams{Where: where})
		resp, err := c.api.VM.StartVM(params)
		if err != nil {
			return TaskRef{}, fmt.Errorf("start vm %s: %w", id, err)
		}
		return firstVMTaskRef(resp.Payload), nil

	case PowerOff:
		if force {
			params := vm.NewPoweroffVMParams()
			params.SetContext(ctx)
			params.SetRequestBody(&models.VMOperateParams{Where: where})
			resp, err := c.api.VM.PoweroffVM(params)
			if err != nil {
				return TaskRef{}, fmt.Errorf("poweroff vm %s: %w", id, err)
			}
			return firstVMTaskRef(resp.Payload), nil
		}
		params := vm.NewShutDownVMParams()
		params.SetContext(ctx)
		params.SetRequestBody(&models.VMOperateParams{Where: where})
		resp, err := c.api.VM.ShutDownVM(params)
		if err != nil {
			return TaskRef{}, fmt.Errorf("shutdown vm %s: %w", id, err)
		}
		return firstVMTaskRef(resp.Payload), nil

	case PowerReset:
		if force {
			params := vm.NewForceRestartVMParams()
			params.SetContext(ctx)
			params.SetRequestBody(&models.VMOperateParams{Where: where})
			resp, err := c.api.VM.ForceRestartVM(params)
			if err != nil {
				return TaskRef{}, fmt.Errorf("force restart vm %s: %w", id, err)
			}
			return firstVMTaskRef(resp.Payload), nil
		}
		params := vm.NewRestartVMParams()
		params.SetContext(ctx)
		params.SetRequestBody(&models.VMOperateParams{Where: where})
		resp, err := c.api.VM.RestartVM(params)
		if err != nil {
			return TaskRef{}, fmt.Errorf("restart vm %s: %w", id, err)
		}
		return firstVMTaskRef(resp.Payload), nil

	case PowerSuspend:
		params := vm.NewSuspendVMParams()
		params.SetContext(ctx)
		params.SetRequestBody(&models.VMOperateParams{Where: where})
		resp, err := c.api.VM.SuspendVM(params)
		if err != nil {
			return TaskRef{}, fmt.Errorf("suspend vm %s: %w", id, err)
		}
		return firstVMTaskRef(resp.Payload), nil

	case PowerResume:
		params := vm.NewResumeVMParams()
		params.SetContext(ctx)
		params.SetRequestBody(&models.VMOperateParams{Where: where})
		resp, err := c.api.VM.ResumeVM(params)
		if err != nil {
			return TaskRef{}, fmt.Errorf("resume vm %s: %w", id, err)
		}
		return firstVMTaskRef(resp.Payload), nil

	default:
		return TaskRef{}, fmt.Errorf("unsupported power action: %s", action)
	}
}

// ---------- UpdateVM ----------

func (c *defaultClient) UpdateVM(ctx context.Context, id string, spec VMUpdateSpec) (TaskRef, error) {
	where := &models.VMWhereInput{ID: pointy.String(id)}
	data := &models.VMUpdateParamsData{}
	if spec.Name != "" {
		data.Name = pointy.String(spec.Name)
	}
	if spec.Description != "" {
		data.Description = pointy.String(spec.Description)
	}
	params := vm.NewUpdateVMParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMUpdateParams{Where: where, Data: data})
	resp, err := c.api.VM.UpdateVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("update vm %s: %w", id, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- MoveToRecycle ----------

func (c *defaultClient) MoveToRecycle(ctx context.Context, id string) (TaskRef, error) {
	where := &models.VMWhereInput{ID: pointy.String(id)}
	params := vm.NewMoveVMToRecycleBinParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMOperateParams{Where: where})
	resp, err := c.api.VM.MoveVMToRecycleBin(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("move vm %s to recycle: %w", id, err)
	}
	return firstDeleteVMTaskRef(resp.Payload), nil
}

// ---------- RecoverFromRecycle ----------

func (c *defaultClient) RecoverFromRecycle(ctx context.Context, id string) (TaskRef, error) {
	where := &models.VMWhereInput{ID: pointy.String(id)}
	params := vm.NewRecoverVMFromRecycleBinParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMOperateParams{Where: where})
	resp, err := c.api.VM.RecoverVMFromRecycleBin(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("recover vm %s from recycle: %w", id, err)
	}
	return firstDeleteVMTaskRef(resp.Payload), nil
}

// ---------- ShutDownVM ----------

func (c *defaultClient) ShutDownVM(ctx context.Context, id string) (TaskRef, error) {
	where := &models.VMWhereInput{ID: pointy.String(id)}
	params := vm.NewShutDownVMParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMOperateParams{Where: where})
	resp, err := c.api.VM.ShutDownVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("shutdown vm %s: %w", id, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- AddDisk ----------

func (c *defaultClient) AddDisk(ctx context.Context, vmID string, spec DiskAddSpec) (TaskRef, error) {
	bus := models.BusSCSI
	switch spec.Bus {
	case "IDE":
		bus = models.BusIDE
	case "VIRTIO", "NVMe", "NVME":
		bus = models.BusVIRTIO
	}

	bootVal := int32(0)
	if spec.Boot > 0 {
		bootVal = spec.Boot
	}

	data := &models.VMAddDiskParamsData{
		VMDisks: &models.VMAddDiskParamsDataVMDisks{
			MountNewCreateDisks: []*models.MountNewCreateDisksParams{
				{
					Boot:  pointy.Int32(bootVal),
					Bus:   &bus,
					Index: pointy.Int32(spec.Index),
					VMVolume: &models.MountNewCreateDisksParamsVMVolume{
						Name: pointy.String(spec.Name),
						Size: pointy.Int64(spec.SizeBytes),
					},
					MaxIops: pointy.Int64(spec.IOPSMax),
				},
			},
		},
	}

	params := vm.NewAddVMDiskParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMAddDiskParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data:  data,
	})
	resp, err := c.api.VM.AddVMDisk(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("add disk to vm %s: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- ExpandDisk ----------

func (c *defaultClient) ExpandDisk(ctx context.Context, vmID, diskID string, sizeBytes int64) (TaskRef, error) {
	params := vm.NewExpandVMDiskParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMExpandVMDiskParams{
		Where: &models.VMDiskWhereInput{ID: pointy.String(diskID), VM: &models.VMWhereInput{ID: pointy.String(vmID)}},
		Size:  pointy.Int64(sizeBytes),
	})
	resp, err := c.api.VM.ExpandVMDisk(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("expand disk %s: %w", diskID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- RemoveDisk ----------

func (c *defaultClient) RemoveDisk(ctx context.Context, vmID, diskID string) (TaskRef, error) {
	params := vm.NewRemoveVMDiskParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMRemoveDiskParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data: &models.VMRemoveDiskParamsData{
			DiskIds: []string{diskID},
		},
	})
	resp, err := c.api.VM.RemoveVMDisk(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("remove disk %s: %w", diskID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- AddCdRom ----------

func (c *defaultClient) AddCdRom(ctx context.Context, vmID string, isoPath string) (TaskRef, error) {
	params := vm.NewAddVMCdRomParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMAddCdRomParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data: &models.VMAddCdRomParamsData{
			VMCdRoms: []*models.VMCdRomParams{
				{
					Boot:        pointy.Int32(0),
					ElfImageID: pointy.String(isoPath),
				},
			},
		},
	})
	resp, err := c.api.VM.AddVMCdRom(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("add cdrom to vm %s: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- EjectCdRom ----------

func (c *defaultClient) EjectCdRom(ctx context.Context, cdromID string) (TaskRef, error) {
	params := vm.NewEjectIsoFromVMCdRomParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMEjectCdRomParams{
		Where: &models.VMDiskWhereInput{ID: pointy.String(cdromID)},
	})
	resp, err := c.api.VM.EjectIsoFromVMCdRom(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("eject cdrom %s: %w", cdromID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- RemoveCdRom ----------

func (c *defaultClient) RemoveCdRom(ctx context.Context, vmID, cdromID string) (TaskRef, error) {
	params := vm.NewRemoveVMCdRomParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMRemoveCdRomParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data: &models.VMRemoveCdRomParamsData{
			CdRomIds: []string{cdromID},
		},
	})
	resp, err := c.api.VM.RemoveVMCdRom(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("remove cdrom %s: %w", cdromID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- AddNic ----------

func (c *defaultClient) AddNic(ctx context.Context, vmID string, spec NicAddSpec) (TaskRef, error) {
	nicType := models.VMNicTypeVLAN
	if spec.Type == "VPC" {
		nicType = models.VMNicTypeVPC
	}
	model := models.VMNicModel(spec.Model)
	if spec.Model == "" {
		model = models.VMNicModelVIRTIO
	}
	params := vm.NewAddVMNicParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMAddNicParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data: &models.VMAddNicParamsData{
			VMNics: []*models.VMNicParams{
				{
					Type:  &nicType,
					Model: &model,
				},
			},
		},
	})
	resp, err := c.api.VM.AddVMNic(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("add nic to vm %s: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- RemoveNic ----------

func (c *defaultClient) RemoveNic(ctx context.Context, vmID string, nicIndex int32) (TaskRef, error) {
	params := vm.NewRemoveVMNicParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMRemoveNicParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data: &models.VMRemoveNicParamsData{
			NicIndex: []int32{nicIndex},
		},
	})
	resp, err := c.api.VM.RemoveVMNic(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("remove nic from vm %s: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- AddGpuDevice ----------

func (c *defaultClient) AddGpuDevice(ctx context.Context, vmID, gpuDeviceID string) (TaskRef, error) {
	params := vm.NewAddVMGpuDeviceParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMAddGpuDeviceParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data: []*models.VMGpuOperationParams{
			{GpuID: pointy.String(gpuDeviceID)},
		},
	})
	resp, err := c.api.VM.AddVMGpuDevice(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("add gpu device to vm %s: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- RemoveGpuDevice ----------

func (c *defaultClient) RemoveGpuDevice(ctx context.Context, vmID, gpuDeviceID string) (TaskRef, error) {
	params := vm.NewRemoveVMGpuDeviceParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMRemoveGpuDeviceParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data: []*models.VMGpuOperationParams{
			{GpuID: pointy.String(gpuDeviceID)},
		},
	})
	resp, err := c.api.VM.RemoveVMGpuDevice(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("remove gpu device from vm %s: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- InstallVmtools ----------

func (c *defaultClient) InstallVmtools(ctx context.Context, vmID string) (TaskRef, error) {
	params := vm.NewInstallVmtoolsParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.InstallVmtoolsParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
	})
	resp, err := c.api.VM.InstallVmtools(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("install vmtools on vm %s: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- GetVNCInfo ----------

func (c *defaultClient) GetVNCInfo(ctx context.Context, vmID string) (*VNCInfo, error) {
	params := vm.NewGetVMVncInfoParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetVMVncInfoParams{
		VM: &models.VMWhereUniqueInput{ID: pointy.String(vmID)},
	})
	resp, err := c.api.VM.GetVMVncInfo(params)
	if err != nil {
		return nil, fmt.Errorf("get vnc info for vm %s: %w", vmID, err)
	}
	info := &VNCInfo{}
	if resp.Payload.ClusterIP != nil {
		info.ClusterIP = *resp.Payload.ClusterIP
	}
	if resp.Payload.Redirect != nil {
		info.Redirect = *resp.Payload.Redirect
	}
	if resp.Payload.Terminal != nil {
		info.Terminal = *resp.Payload.Terminal
	}
	if resp.Payload.Direct != nil {
		info.Direct = *resp.Payload.Direct
	}
	return info, nil
}

// ---------- CreateVMFromTemplate ----------

func (c *defaultClient) CreateVMFromTemplate(ctx context.Context, spec VMCreateFromTemplateSpec) (TaskRef, error) {
	fw := models.VMFirmwareBIOS
	if strings.EqualFold(spec.Firmware, "UEFI") {
		fw = models.VMFirmwareUEFI
	}
	isFullCopy := spec.IsFullCopy
	p := &models.VMCreateVMFromContentLibraryTemplateParams{
		TemplateID: pointy.String(spec.TemplateID),
		Name:       pointy.String(spec.Name),
		ClusterID:  pointy.String(spec.ClusterID),
		IsFullCopy: &isFullCopy,
		Firmware:  &fw,
	}
	if spec.HostID != "" {
		p.HostID = pointy.String(spec.HostID)
	}
	if spec.VCPU > 0 {
		p.Vcpu = pointy.Int32(spec.VCPU)
	}
	if spec.MemoryBytes > 0 {
		p.Memory = pointy.Int64(spec.MemoryBytes)
	}
	if spec.Description != "" {
		p.Description = pointy.String(spec.Description)
	}

	// NIC 配置
	if spec.NIC.Type != "" || spec.NIC.Model != "" || spec.NIC.VlanID != "" {
		nicType := models.VMNicTypeVLAN
		if spec.NIC.Type == "VPC" {
			nicType = models.VMNicTypeVPC
		}
		model := models.VMNicModel(strings.ToUpper(spec.NIC.Model))
		if model == "" {
			model = models.VMNicModelVIRTIO
		}
		enabled := true
		vmNics := []*models.VMNicParams{
			{
				Enabled: &enabled,
				Type:   &nicType,
				Model:  &model,
			},
		}
		if spec.NIC.VlanID != "" {
			vmNics[0].ConnectVlanID = pointy.String(spec.NIC.VlanID)
		}
		p.VMNics = vmNics
	}

	params := vm.NewCreateVMFromContentLibraryTemplateParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.VMCreateVMFromContentLibraryTemplateParams{p})

	resp, err := c.api.VM.CreateVMFromContentLibraryTemplate(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("create vm from template %s: %w", spec.TemplateID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- MigrateAcrossCluster ----------

func (c *defaultClient) MigrateAcrossCluster(ctx context.Context, vmID, targetClusterID, hostID string) (TaskRef, error) {
	params := vm.NewMigrateVMAcrossClusterParams()
	params.SetContext(ctx)
	data := &models.VMMigrateAcrossClusterParamsData{
		ClusterID: pointy.String(targetClusterID),
		VMConfig:  &models.MigrateVMConfig{},
	}
	if hostID != "" {
		data.HostID = pointy.String(hostID)
	}
	params.SetRequestBody(&models.VMMigrateAcrossClusterParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data:  data,
	})
	resp, err := c.api.VM.MigrateVMAcrossCluster(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("migrate vm %s across cluster: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- 内部辅助 ----------

// toVM 把 SDK models.VM 转成内部 VM 模型。
func toVM(v *models.VM) VM {
	out := VM{}
	if v.ID != nil {
		out.ID = *v.ID
	}
	if v.Name != nil {
		out.Name = *v.Name
	}
	if v.Status != nil {
		out.Status = string(*v.Status)
	}
	if v.Description != nil {
		out.Description = *v.Description
	}
	if v.Vcpu != nil {
		out.VCPU = *v.Vcpu
	}
	if v.Memory != nil {
		out.MemoryBytes = uint64(*v.Memory)
	}
	if v.Ips != nil && *v.Ips != "" {
		out.IPs = strings.Split(*v.Ips, ",")
		for i := range out.IPs {
			out.IPs[i] = strings.TrimSpace(out.IPs[i])
		}
	}
	if v.Cluster != nil && v.Cluster.ID != nil {
		out.ClusterID = *v.Cluster.ID
	}
	if v.Cluster != nil && v.Cluster.Name != nil {
		out.ClusterName = *v.Cluster.Name
	}
	if v.Host != nil && v.Host.ID != nil {
		out.HostID = *v.Host.ID
	}
	if v.Host != nil && v.Host.Name != nil {
		out.HostName = *v.Host.Name
	}
	if v.Host != nil && v.Host.ManagementIP != nil {
		out.HostIP = *v.Host.ManagementIP
	}
	if v.Firmware != nil {
		out.Firmware = string(*v.Firmware)
	}
	if v.Ha != nil {
		out.Ha = *v.Ha
	}
	if v.GuestOsType != nil {
		out.GuestOS = string(*v.GuestOsType)
	}
	if v.VMToolsStatus != nil {
		out.VMToolsStatus = string(*v.VMToolsStatus)
	}
	if v.VMToolsVersion != nil {
		out.VMToolsVersion = *v.VMToolsVersion
	}
	if v.CPUModel != nil {
		out.CPUModel = *v.CPUModel
	}
	if v.DNSServers != nil {
		out.DNSServers = *v.DNSServers
	}
	if v.Hostname != nil {
		out.Hostname = *v.Hostname
	}
	if v.VMDisks != nil {
		out.DiskCount = len(v.VMDisks)
	}
	if v.VMNics != nil {
		out.NicCount = len(v.VMNics)
	}
	if v.ProvisionedSize != nil {
		out.ProvisionedBytes = uint64(*v.ProvisionedSize)
	}
	if v.UsedSize != nil {
		out.UsedBytes = uint64(*v.UsedSize)
	}
	if v.InRecycleBin != nil {
		out.InRecycleBin = *v.InRecycleBin
	}
	if v.Protected != nil {
		out.Protected = *v.Protected
	}
	if v.LocalCreatedAt != nil {
		out.CreatedAt = *v.LocalCreatedAt
	}
	if v.BiosUUID != nil {
		out.BiosUUID = *v.BiosUUID
	}
	if v.CPUUsage != nil {
		out.CPUUsage = *v.CPUUsage
	}
	if v.MemoryUsage != nil {
		out.MemoryUsage = *v.MemoryUsage
	}
	if v.GuestSizeUsage != nil {
		out.GuestSizeUsage = *v.GuestSizeUsage
	}
	if v.GuestUsedSize != nil {
		out.GuestUsedSize = *v.GuestUsedSize
	}
	if v.LogicalSizeBytes != nil {
		out.LogicalSizeBytes = *v.LogicalSizeBytes
	}
	if v.VideoType != nil {
		out.VideoType = string(*v.VideoType)
	}
	if v.NestedVirtualization != nil {
		out.NestedVirt = *v.NestedVirtualization
	}
	if v.HaPriority != nil {
		out.HaPriority = string(*v.HaPriority)
	}
	if v.CloudInitSupported != nil {
		out.CloudInit = *v.CloudInitSupported
	}
	if v.Labels != nil {
		for _, l := range v.Labels {
			if l != nil {
				label := *l.Key
				if l.Value != nil {
					label = *l.Key + "=" + *l.Value
				}
				out.Labels = append(out.Labels, label)
			}
		}
	}
	if v.UsbDevices != nil {
		for _, d := range v.UsbDevices {
			if d != nil {
				dev := UsbDevice{}
				if d.ID != nil {
					dev.ID = *d.ID
				}
				if d.Name != nil {
					dev.Name = *d.Name
				}
				out.UsbDevices = append(out.UsbDevices, dev)
			}
		}
	}
	if v.GpuDevices != nil {
		for _, d := range v.GpuDevices {
			if d != nil {
				dev := GpuDevice{}
				if d.ID != nil {
					dev.ID = *d.ID
				}
				if d.Name != nil {
					dev.Name = *d.Name
				}
				out.GpuDevices = append(out.GpuDevices, dev)
			}
		}
	}
	return out
}

// firstVMTaskRef 从 WithTaskVM 数组提取第一个 TaskRef。
func firstVMTaskRef(items []*models.WithTaskVM) TaskRef {
	if len(items) == 0 {
		return TaskRef{}
	}
	ref := TaskRef{}
	it := items[0]
	if it.TaskID != nil {
		ref.ID = *it.TaskID
	}
	if it.Data != nil && it.Data.ID != nil {
		ref.EntityID = *it.Data.ID
		ref.EntityKind = "VM"
	}
	return ref
}

// firstDeleteVMTaskRef 从 WithTaskDeleteVM 数组提取 TaskRef。
func firstDeleteVMTaskRef(items []*models.WithTaskDeleteVM) TaskRef {
	if len(items) == 0 {
		return TaskRef{}
	}
	ref := TaskRef{}
	it := items[0]
	if it.TaskID != nil {
		ref.ID = *it.TaskID
	}
	if it.Data != nil && it.Data.ID != nil {
		ref.EntityID = *it.Data.ID
		ref.EntityKind = "VM"
	}
	return ref
}

// firstExportTaskRef 从 WithTaskVMExportFile 数组提取 TaskRef。
func firstExportTaskRef(items []*models.WithTaskVMExportFile) TaskRef {
	if len(items) == 0 {
		return TaskRef{}
	}
	ref := TaskRef{}
	it := items[0]
	if it.TaskID != nil {
		ref.ID = *it.TaskID
	}
	return ref
}

// ---------- ListVMDisks ----------

func (c *defaultClient) ListVMDisks(ctx context.Context, vmID string) ([]VMDisk, error) {
	params := vm_disk.NewGetVMDisksParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetVMDisksRequestBody{
		Where: &models.VMDiskWhereInput{
			VM: &models.VMWhereInput{ID: pointy.String(vmID)},
		},
	})
	resp, err := c.api.VMDisk.GetVMDisks(params)
	if err != nil {
		return nil, fmt.Errorf("list vm disks: %w", err)
	}
	out := make([]VMDisk, 0, len(resp.Payload))
	for _, d := range resp.Payload {
		out = append(out, toVMDisk(d))
	}
	return out, nil
}

func toVMDisk(d *models.VMDisk) VMDisk {
	out := VMDisk{}
	if d.ID != nil {
		out.ID = *d.ID
	}
	if d.Boot != nil {
		out.Boot = *d.Boot
	}
	if d.Bus != nil {
		out.Bus = string(*d.Bus)
	}
	if d.Key != nil {
		out.Key = *d.Key
	}
	if d.MaxBandwidth != nil {
		out.MaxBandwidth = pointy.Int64(*d.MaxBandwidth)
	}
	if d.MaxIops != nil {
		out.MaxIops = pointy.Int64(int64(*d.MaxIops))
	}
	if d.Type != nil {
		out.Type = string(*d.Type)
	}
	if d.VM != nil && d.VM.ID != nil {
		out.VMID = *d.VM.ID
	}
	if d.VMVolume != nil {
		if d.VMVolume.ID != nil {
			out.VolumeID = *d.VMVolume.ID
		}
		if d.VMVolume.Name != nil {
			out.VolumeName = *d.VMVolume.Name
		}
	}
	if d.ElfImage != nil {
		if d.ElfImage.ID != nil {
			out.ElfImageID = *d.ElfImage.ID
		}
		if d.ElfImage.Name != nil {
			out.ElfImageName = *d.ElfImage.Name
		}
	}
	return out
}

// ---------- ListVMNics ----------

func (c *defaultClient) ListVMNics(ctx context.Context, vmID string) ([]VMNic, error) {
	params := vm_nic.NewGetVMNicsParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetVMNicsRequestBody{
		Where: &models.VMNicWhereInput{
			VM: &models.VMWhereInput{ID: pointy.String(vmID)},
		},
	})
	resp, err := c.api.VMNic.GetVMNics(params)
	if err != nil {
		return nil, fmt.Errorf("list vm nics: %w", err)
	}
	out := make([]VMNic, 0, len(resp.Payload))
	for _, n := range resp.Payload {
		out = append(out, toVMNic(n))
	}
	return out, nil
}

func toVMNic(n *models.VMNic) VMNic {
	out := VMNic{}
	if n.ID != nil {
		out.ID = *n.ID
	}
	if n.LocalID != nil {
		out.LocalID = *n.LocalID
	}
	if n.MacAddress != nil {
		out.MacAddress = *n.MacAddress
	}
	if n.Model != nil {
		out.Model = string(*n.Model)
	}
	if n.Type != nil {
		out.Type = string(*n.Type)
	}
	if n.Gateway != nil {
		out.Gateway = *n.Gateway
	}
	if n.SubnetMask != nil {
		out.SubnetMask = *n.SubnetMask
	}
	if n.IPAddress != nil {
		out.IPAddress = *n.IPAddress
	}
	if n.Enabled != nil {
		out.Enabled = *n.Enabled
	}
	if n.VM != nil && n.VM.ID != nil {
		out.VMID = *n.VM.ID
	}
	if n.Vlan != nil {
		if n.Vlan.ID != nil {
			out.VlanID = *n.Vlan.ID
		}
		if n.Vlan.Name != nil {
			out.VlanName = *n.Vlan.Name
		}
	}
	if n.IngressRateLimitMaxRateInBitps != nil {
		rate := int64(*n.IngressRateLimitMaxRateInBitps)
		out.IngressRateLimit = &rate
	}
	if n.EgressRateLimitMaxRateInBitps != nil {
		rate := int64(*n.EgressRateLimitMaxRateInBitps)
		out.EgressRateLimit = &rate
	}
	return out
}

// ---------- UpdateNic ----------

func (c *defaultClient) UpdateNic(ctx context.Context, nicID string, spec VMNicUpdateSpec) (TaskRef, error) {
	data := &models.VMUpdateNicParamsData{
		NicID: pointy.String(nicID),
	}
	if spec.Enabled != nil {
		data.Enabled = spec.Enabled
	}
	if spec.Gateway != "" {
		data.Gateway = pointy.String(spec.Gateway)
	}
	if spec.IPAddress != "" {
		data.IPAddress = pointy.String(spec.IPAddress)
	}
	if spec.MacAddress != "" {
		data.MacAddress = pointy.String(spec.MacAddress)
	}
	if spec.Model != "" {
		model := models.VMNicModel(spec.Model)
		data.Model = &model
	}
	if spec.SubnetMask != "" {
		data.SubnetMask = pointy.String(spec.SubnetMask)
	}

	params := vm.NewUpdateVMNicParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMUpdateNicParams{
		Data: data,
	})
	resp, err := c.api.VM.UpdateVMNic(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("update nic %s: %w", nicID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- UpdateDisk ----------

func (c *defaultClient) UpdateDisk(ctx context.Context, diskID string, spec DiskUpdateSpec) (TaskRef, error) {
	data := &models.VMUpdateDiskParamsData{
		VMDiskID: pointy.String(diskID),
	}

	params := vm.NewUpdateVMDiskParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMUpdateDiskParams{
		Where: &models.VMWhereInput{ID: pointy.String(diskID)},
		Data:  data,
	})
	resp, err := c.api.VM.UpdateVMDisk(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("update disk %s: %w", diskID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- ToggleCdRom ----------

func (c *defaultClient) ToggleCdRom(ctx context.Context, cdromID string, spec CdRomToggleSpec) (TaskRef, error) {
	params := vm.NewToggleVMCdRomDisableParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMToggleCdRomDisableParams{
		Where:    &models.VMDiskWhereInput{ID: pointy.String(cdromID)},
		Disabled: pointy.Bool(spec.Disabled),
	})
	resp, err := c.api.VM.ToggleVMCdRomDisable(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("toggle cdrom %s: %w", cdromID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- ResetPassword ----------

func (c *defaultClient) ResetPassword(ctx context.Context, vmID string, spec ResetPasswordSpec) (TaskRef, error) {
	params := vm.NewResetVMGuestOsPasswordParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMResetGuestOsPasswordParams{
		Where: &models.VMWhereInput{ID: pointy.String(vmID)},
		Data: &models.VMResetGuestOsPasswordParamsData{
			Username: pointy.String(spec.Username),
			Password: pointy.String(spec.Password),
		},
	})
	resp, err := c.api.VM.ResetVMGuestOsPassword(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("reset password for vm %s: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- RebuildVM ----------

func (c *defaultClient) RebuildVM(ctx context.Context, vmID string, spec RebuildVMSpec) (TaskRef, error) {
	data := &models.VMRebuildParams{
		Name: pointy.String(spec.Name),
		RebuildFromSnapshotID: pointy.String(spec.SnapshotID),
	}
	if spec.ClusterID != "" {
		data.ClusterID = pointy.String(spec.ClusterID)
	}
	if spec.HostID != "" {
		data.HostID = pointy.String(spec.HostID)
	}

	params := vm.NewRebuildVMParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.VMRebuildParams{data})
	resp, err := c.api.VM.RebuildVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("rebuild vm %s: %w", vmID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// ---------- AbortMigrateAcrossCluster ----------

func (c *defaultClient) AbortMigrateAcrossCluster(ctx context.Context, vmID string) (TaskRef, error) {
	params := vm.NewAbortMigrateVMAcrossClusterParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.AbortMigrateVMAcrossClusterParams{
		Tasks: &models.TaskWhereInput{},
	})
	resp, err := c.api.VM.AbortMigrateVMAcrossCluster(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("abort migrate vm %s: %w", vmID, err)
	}
	if len(resp.Payload) == 0 {
		return TaskRef{}, nil
	}
	ref := TaskRef{EntityKind: "Task"}
	if resp.Payload[0].ID != nil {
		ref.ID = *resp.Payload[0].ID
	}
	return ref, nil
}

// ---------- ConvertToVM ----------

func (c *defaultClient) ConvertToVM(ctx context.Context, templateID string) (TaskRef, error) {
	params := vm.NewConvertVMTemplateToVMParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.ConvertVMTemplateToVMParams{
		{
			ConvertedFromTemplateID: pointy.String(templateID),
			Name: pointy.String(""),
		},
	})
	resp, err := c.api.VM.ConvertVMTemplateToVM(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("convert template %s to vm: %w", templateID, err)
	}
	return firstVMTaskRef(resp.Payload), nil
}

// modelVMStatusStopped 返回 VM 创建时的初始状态。
func modelVMStatusStopped() *models.VMStatus {
	s := models.VMStatusSTOPPED
	return &s
}
