package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	sdkhost "github.com/smartxworks/cloudtower-go-sdk/v2/client/host"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// HostOps 定义主机相关操作。
type HostOps interface {
	ListHosts(ctx context.Context, opts ListOpts) ([]Host, error)
	GetHost(ctx context.Context, id string) (*Host, error)
	GetHostByName(ctx context.Context, name string) (*Host, error)
	ListHostsByCluster(ctx context.Context, clusterID string) ([]Host, error)
	EnterMaintenanceMode(ctx context.Context, id string) (TaskRef, error)
	ExitMaintenanceMode(ctx context.Context, id string) (TaskRef, error)
	ShutdownHost(ctx context.Context, id string, force bool) (TaskRef, error)
	RebootHost(ctx context.Context, id string, force bool) (TaskRef, error)
}

func (c *defaultClient) ListHosts(ctx context.Context, opts ListOpts) ([]Host, error) {
	params := sdkhost.NewGetHostsParams()
	params.SetContext(ctx)

	body := &models.GetHostsRequestBody{}
	where := &models.HostWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if opts.ClusterID != "" {
		where.Cluster = &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}
		hasWhere = true
	}
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

	resp, err := c.api.Host.GetHosts(params)
	if err != nil {
		return nil, fmt.Errorf("list hosts: %w", err)
	}
	out := make([]Host, 0, len(resp.Payload))
	for _, h := range resp.Payload {
		out = append(out, toHost(h))
	}
	return out, nil
}

func (c *defaultClient) GetHost(ctx context.Context, id string) (*Host, error) {
	params := sdkhost.NewGetHostsParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetHostsRequestBody{
		Where: &models.HostWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.Host.GetHosts(params)
	if err != nil {
		return nil, fmt.Errorf("get host %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get host %s: %w", id, ErrNotFound)
	}
	h := toHost(resp.Payload[0])
	return &h, nil
}

func (c *defaultClient) GetHostByName(ctx context.Context, name string) (*Host, error) {
	params := sdkhost.NewGetHostsParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetHostsRequestBody{
		Where: &models.HostWhereInput{Name: &name},
	})
	resp, err := c.api.Host.GetHosts(params)
	if err != nil {
		return nil, fmt.Errorf("get host by name %s: %w", name, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get host by name %s: %w", name, ErrNotFound)
	}
	h := toHost(resp.Payload[0])
	return &h, nil
}

func (c *defaultClient) ListHostsByCluster(ctx context.Context, clusterID string) ([]Host, error) {
	params := sdkhost.NewGetHostsParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetHostsRequestBody{
		Where: &models.HostWhereInput{
			Cluster: &models.ClusterWhereInput{ID: pointy.String(clusterID)},
		},
	})
	resp, err := c.api.Host.GetHosts(params)
	if err != nil {
		return nil, fmt.Errorf("list hosts by cluster %s: %w", clusterID, err)
	}
	out := make([]Host, 0, len(resp.Payload))
	for _, h := range resp.Payload {
		out = append(out, toHost(h))
	}
	return out, nil
}

func (c *defaultClient) EnterMaintenanceMode(ctx context.Context, id string) (TaskRef, error) {
	params := sdkhost.NewEnterMaintenanceModeParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.EnterMaintenanceModeParams{
		Where: &models.EnterMaintenanceModeParamsWhere{HostID: pointy.String(id)},
		Data:  &models.EnterMaintenanceModeInput{},
	})
	resp, err := c.api.Host.EnterMaintenanceMode(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("enter maintenance %s: %w", id, err)
	}
	return singleHostTaskRef(resp.Payload), nil
}

func (c *defaultClient) ExitMaintenanceMode(ctx context.Context, id string) (TaskRef, error) {
	params := sdkhost.NewExitMaintenanceModeParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.ExitMaintenanceModeParams{
		Where: &models.ExitMaintenanceModeParamsWhere{HostID: pointy.String(id)},
		Data:  &models.ExitMaintenanceModeInput{},
	})
	resp, err := c.api.Host.ExitMaintenanceMode(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("exit maintenance %s: %w", id, err)
	}
	return singleHostTaskRef(resp.Payload), nil
}

func (c *defaultClient) ShutdownHost(ctx context.Context, id string, force bool) (TaskRef, error) {
	action := models.OperateActionEnumPoweroff
	params := sdkhost.NewPowerOffHostParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.OperateHostPowerParams{
		Where: &models.OperateHostPowerParamsWhere{HostID: pointy.String(id)},
		Data: &models.OperateHostPowerData{
			Action: &action,
			Force:  pointy.Bool(force),
		},
	})
	resp, err := c.api.Host.PowerOffHost(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("shutdown host %s: %w", id, err)
	}
	return singleHostTaskRef(resp.Payload), nil
}

func (c *defaultClient) RebootHost(ctx context.Context, id string, force bool) (TaskRef, error) {
	action := models.OperateActionEnumReboot
	params := sdkhost.NewPowerOffHostParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.OperateHostPowerParams{
		Where: &models.OperateHostPowerParamsWhere{HostID: pointy.String(id)},
		Data: &models.OperateHostPowerData{
			Action: &action,
			Force:  pointy.Bool(force),
		},
	})
	resp, err := c.api.Host.PowerOffHost(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("reboot host %s: %w", id, err)
	}
	return singleHostTaskRef(resp.Payload), nil
}

// ---------- helpers ----------

func toHost(h *models.Host) Host {
	out := Host{}
	if h.ID != nil {
		out.ID = *h.ID
	}
	if h.Name != nil {
		out.Name = *h.Name
	}
	if h.Status != nil {
		out.Status = string(*h.Status)
	}
	if h.ManagementIP != nil {
		out.ManagementIP = *h.ManagementIP
	}
	if h.DataIP != nil {
		out.DataIP = *h.DataIP
	}
	if h.CPUModel != nil {
		out.CPUModel = *h.CPUModel
	}
	if h.TotalMemoryBytes != nil {
		out.TotalMemoryBytes = uint64(*h.TotalMemoryBytes)
	}
	if h.RunningVMNum != nil {
		out.RunningVMs = *h.RunningVMNum
	}
	if h.Cluster != nil && h.Cluster.ID != nil {
		out.ClusterID = *h.Cluster.ID
	}
	return out
}

// singleHostTaskRef 处理 *WithTaskHost（单体，不是数组）。
func singleHostTaskRef(wt *models.WithTaskHost) TaskRef {
	if wt == nil {
		return TaskRef{}
	}
	ref := TaskRef{EntityKind: "Host"}
	if wt.TaskID != nil {
		ref.ID = *wt.TaskID
	}
	if wt.Data != nil && wt.Data.ID != nil {
		ref.EntityID = *wt.Data.ID
	}
	return ref
}
