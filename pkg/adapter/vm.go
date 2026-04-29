package adapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm"
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

	p := &models.VMCreationParams{
		ClusterID:  pointy.String(spec.ClusterID),
		Name:       pointy.String(spec.Name),
		Ha:        pointy.Bool(true),
		CPUSockets: pointy.Int32(sockets),
		CPUCores:   pointy.Int32(cores),
		Memory:     pointy.Int64(spec.MemoryBytes),
		Firmware:   &fw,
		Status:     modelVMStatusStopped(),
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
	if v.Host != nil && v.Host.Name != nil {
		out.HostName = *v.Host.Name
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

// modelVMStatusStopped 返回 VM 创建时的初始状态。
func modelVMStatusStopped() *models.VMStatus {
	s := models.VMStatusSTOPPED
	return &s
}
