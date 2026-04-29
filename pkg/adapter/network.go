package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	sdkvds "github.com/smartxworks/cloudtower-go-sdk/v2/client/vds"
	sdkvlan "github.com/smartxworks/cloudtower-go-sdk/v2/client/vlan"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// NetworkOps 定义虚拟交换机（VDS）操作。
type NetworkOps interface {
	ListNetworks(ctx context.Context, opts ListOpts) ([]Network, error)
	GetNetwork(ctx context.Context, id string) (*Network, error)
}

// VLANOps 定义 VLAN 操作。
type VLANOps interface {
	ListVLANs(ctx context.Context, opts ListOpts) ([]VLAN, error)
	GetVLAN(ctx context.Context, id string) (*VLAN, error)
	CreateVLAN(ctx context.Context, spec VLANCreateSpec) (TaskRef, error)
	DeleteVLAN(ctx context.Context, id string) (TaskRef, error)
}

// --- Network (VDS) ---

func (c *defaultClient) ListNetworks(ctx context.Context, opts ListOpts) ([]Network, error) {
	params := sdkvds.NewGetVdsesParams()
	params.SetContext(ctx)
	body := &models.GetVdsesRequestBody{}
	where := &models.VdsWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if hasWhere {
		body.Where = where
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	params.SetRequestBody(body)
	resp, err := c.api.Vds.GetVdses(params)
	if err != nil {
		return nil, fmt.Errorf("list networks: %w", err)
	}
	out := make([]Network, 0, len(resp.Payload))
	for _, v := range resp.Payload {
		out = append(out, toNetwork(v))
	}
	return out, nil
}

func (c *defaultClient) GetNetwork(ctx context.Context, id string) (*Network, error) {
	params := sdkvds.NewGetVdsesParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetVdsesRequestBody{
		Where: &models.VdsWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.Vds.GetVdses(params)
	if err != nil {
		return nil, fmt.Errorf("get network %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get network %s: %w", id, ErrNotFound)
	}
	n := toNetwork(resp.Payload[0])
	return &n, nil
}

// --- VLAN ---

func (c *defaultClient) ListVLANs(ctx context.Context, opts ListOpts) ([]VLAN, error) {
	params := sdkvlan.NewGetVlansParams()
	params.SetContext(ctx)
	body := &models.GetVlansRequestBody{}
	where := &models.VlanWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if hasWhere {
		body.Where = where
	}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	params.SetRequestBody(body)
	resp, err := c.api.Vlan.GetVlans(params)
	if err != nil {
		return nil, fmt.Errorf("list vlans: %w", err)
	}
	out := make([]VLAN, 0, len(resp.Payload))
	for _, v := range resp.Payload {
		out = append(out, toVLAN(v))
	}
	return out, nil
}

func (c *defaultClient) GetVLAN(ctx context.Context, id string) (*VLAN, error) {
	params := sdkvlan.NewGetVlansParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetVlansRequestBody{
		Where: &models.VlanWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.Vlan.GetVlans(params)
	if err != nil {
		return nil, fmt.Errorf("get vlan %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get vlan %s: %w", id, ErrNotFound)
	}
	v := toVLAN(resp.Payload[0])
	return &v, nil
}

func (c *defaultClient) CreateVLAN(ctx context.Context, spec VLANCreateSpec) (TaskRef, error) {
	params := sdkvlan.NewCreateVMVlanParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.VMVlanCreationParams{
		{
			Name:  pointy.String(spec.Name),
			VdsID: pointy.String(spec.VdsID),
		},
	})
	resp, err := c.api.Vlan.CreateVMVlan(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("create vlan %s: %w", spec.Name, err)
	}
	return firstVlanTaskRef(resp.Payload), nil
}

func (c *defaultClient) DeleteVLAN(ctx context.Context, id string) (TaskRef, error) {
	params := sdkvlan.NewDeleteVlanParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VlanDeletionParams{
		Where: &models.VlanWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.Vlan.DeleteVlan(params)
	if err != nil {
		return TaskRef{}, fmt.Errorf("delete vlan %s: %w", id, err)
	}
	return firstDeleteVlanTaskRef(resp.Payload), nil
}

// --- helpers ---

func toNetwork(v *models.Vds) Network {
	out := Network{}
	if v.ID != nil { out.ID = *v.ID }
	if v.Name != nil { out.Name = *v.Name }
	if v.Type != nil { out.Type = string(*v.Type) }
	if v.Cluster != nil && v.Cluster.ID != nil { out.ClusterID = *v.Cluster.ID }
	return out
}

func toVLAN(v *models.Vlan) VLAN {
	out := VLAN{}
	if v.ID != nil { out.ID = *v.ID }
	if v.Name != nil { out.Name = *v.Name }
	if v.VlanID != nil { out.VlanTag = int32(*v.VlanID) }
	if v.Type != nil { out.Type = string(*v.Type) }
	if v.Vds != nil && v.Vds.ID != nil { out.VdsID = *v.Vds.ID }
	return out
}

func firstVlanTaskRef(items []*models.WithTaskVlan) TaskRef {
	if len(items) == 0 { return TaskRef{} }
	ref := TaskRef{EntityKind: "VLAN"}
	if items[0].TaskID != nil { ref.ID = *items[0].TaskID }
	if items[0].Data != nil && items[0].Data.ID != nil { ref.EntityID = *items[0].Data.ID }
	return ref
}

func firstDeleteVlanTaskRef(items []*models.WithTaskDeleteVlan) TaskRef {
	if len(items) == 0 { return TaskRef{} }
	ref := TaskRef{EntityKind: "VLAN"}
	if items[0].TaskID != nil { ref.ID = *items[0].TaskID }
	if items[0].Data != nil && items[0].Data.ID != nil { ref.EntityID = *items[0].Data.ID }
	return ref
}
