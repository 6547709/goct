package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/vm_placement_group"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// VMPlacementGroupOps 定义 VM 放置组操作。
type VMPlacementGroupOps interface {
	ListVMPlacementGroups(ctx context.Context, opts ListOpts) ([]VMPlacementGroup, error)
	CreateVMPlacementGroup(ctx context.Context, spec VMPlacementGroupCreateSpec) ([]VMPlacementGroup, error)
	DeleteVMPlacementGroup(ctx context.Context, id string) error
}

// ---------- ListVMPlacementGroups ----------

func (c *defaultClient) ListVMPlacementGroups(ctx context.Context, opts ListOpts) ([]VMPlacementGroup, error) {
	params := vm_placement_group.NewGetVMPlacementGroupsParams()
	params.SetContext(ctx)

	where := &models.VMPlacementGroupWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if opts.ClusterID != "" {
		where.Cluster = &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}
		hasWhere = true
	}

	body := &models.GetVMPlacementGroupsRequestBody{}
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

	resp, err := c.api.VMPlacementGroup.GetVMPlacementGroups(params)
	if err != nil {
		return nil, fmt.Errorf("list vm placement groups: %w", err)
	}
	out := make([]VMPlacementGroup, 0, len(resp.Payload))
	for _, g := range resp.Payload {
		out = append(out, toVMPlacementGroup(g))
	}
	return out, nil
}

// ---------- CreateVMPlacementGroup ----------

func (c *defaultClient) CreateVMPlacementGroup(ctx context.Context, spec VMPlacementGroupCreateSpec) ([]VMPlacementGroup, error) {
	params := vm_placement_group.NewCreateVMPlacementGroupParams()
	params.SetContext(ctx)
	params.SetRequestBody([]*models.VMPlacementGroupCreationParams{
		{
			ClusterID: pointy.String(spec.ClusterID),
			Name:      pointy.String(spec.Name),
		},
	})

	resp, err := c.api.VMPlacementGroup.CreateVMPlacementGroup(params)
	if err != nil {
		return nil, fmt.Errorf("create vm placement group: %w", err)
	}
	out := make([]VMPlacementGroup, 0, len(resp.Payload))
	for _, g := range resp.Payload {
		if g.Data != nil {
			out = append(out, toVMPlacementGroup(g.Data))
		}
	}
	return out, nil
}

// ---------- DeleteVMPlacementGroup ----------

func (c *defaultClient) DeleteVMPlacementGroup(ctx context.Context, id string) error {
	params := vm_placement_group.NewDeleteVMPlacementGroupParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.VMPlacementGroupDeletionParams{
		Where: &models.VMPlacementGroupWhereInput{ID: pointy.String(id)},
	})

	_, err := c.api.VMPlacementGroup.DeleteVMPlacementGroup(params)
	if err != nil {
		return fmt.Errorf("delete vm placement group: %w", err)
	}
	return nil
}

// toVMPlacementGroup 把 SDK models.VMPlacementGroup 转成内部 VMPlacementGroup 模型。
func toVMPlacementGroup(g *models.VMPlacementGroup) VMPlacementGroup {
	out := VMPlacementGroup{}
	if g.ID != nil {
		out.ID = *g.ID
	}
	if g.Name != nil {
		out.Name = *g.Name
	}
	if g.Cluster != nil && g.Cluster.ID != nil {
		out.ClusterID = *g.Cluster.ID
	}
	return out
}