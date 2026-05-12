package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/snapshot_plan"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// SnapshotPlanOps 定义快照计划操作。
type SnapshotPlanOps interface {
	ListSnapshotPlans(ctx context.Context, opts ListOpts) ([]SnapshotPlan, error)
	DeleteSnapshotPlan(ctx context.Context, id string) error
}

// ---------- ListSnapshotPlans ----------

func (c *defaultClient) ListSnapshotPlans(ctx context.Context, opts ListOpts) ([]SnapshotPlan, error) {
	params := snapshot_plan.NewGetSnapshotPlansParams()
	params.SetContext(ctx)

	where := &models.SnapshotPlanWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if opts.ClusterID != "" {
		where.Cluster = &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}
		hasWhere = true
	}

	body := &models.GetSnapshotPlansRequestBody{}
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

	resp, err := c.api.SnapshotPlan.GetSnapshotPlans(params)
	if err != nil {
		return nil, fmt.Errorf("list snapshot plans: %w", err)
	}
	out := make([]SnapshotPlan, 0, len(resp.Payload))
	for _, p := range resp.Payload {
		out = append(out, toSnapshotPlan(p))
	}
	return out, nil
}

// ---------- DeleteSnapshotPlan ----------

func (c *defaultClient) DeleteSnapshotPlan(ctx context.Context, id string) error {
	params := snapshot_plan.NewDeleteSnapshotPlanParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.SnapshotPlanDeletionParams{
		Where: &models.SnapshotPlanWhereInput{ID: pointy.String(id)},
	})

	_, err := c.api.SnapshotPlan.DeleteSnapshotPlan(params)
	if err != nil {
		return fmt.Errorf("delete snapshot plan: %w", err)
	}
	return nil
}

// toSnapshotPlan 把 SDK models.SnapshotPlan 转成内部 SnapshotPlan 模型。
func toSnapshotPlan(p *models.SnapshotPlan) SnapshotPlan {
	out := SnapshotPlan{}
	if p.ID != nil {
		out.ID = *p.ID
	}
	if p.Name != nil {
		out.Name = *p.Name
	}
	if p.Cluster != nil && p.Cluster.ID != nil {
		out.ClusterID = *p.Cluster.ID
	}
	if p.Status != nil {
		out.Status = string(*p.Status)
	}
	if p.ExecutePlanType != nil {
		out.PlanType = string(*p.ExecutePlanType)
	}
	if p.AutoDeleteNum != nil {
		out.Retention = *p.AutoDeleteNum
	}
	if p.StartTime != nil {
		out.StartTime = *p.StartTime
	}
	if p.EndTime != nil {
		out.EndTime = *p.EndTime
	}
	if p.Exechm != nil {
		out.ExecHM = fmt.Sprintf("%v", p.Exechm)
	}
	if p.ExecuteIntervals != nil {
		out.ExecuteIntervals = p.ExecuteIntervals
	}
	return out
}