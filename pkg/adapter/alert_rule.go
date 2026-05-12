package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/alert_rule"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// AlertRuleOps 定义告警规则操作。
type AlertRuleOps interface {
	ListAlertRules(ctx context.Context, opts ListOpts) ([]AlertRule, error)
}

// ---------- ListAlertRules ----------

func (c *defaultClient) ListAlertRules(ctx context.Context, opts ListOpts) ([]AlertRule, error) {
	params := alert_rule.NewGetAlertRulesParams()
	params.SetContext(ctx)

	where := &models.AlertRuleWhereInput{}
	hasWhere := false
	if opts.ClusterID != "" {
		where.Cluster = &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}
		hasWhere = true
	}

	body := &models.GetAlertRulesRequestBody{}
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

	resp, err := c.api.AlertRule.GetAlertRules(params)
	if err != nil {
		return nil, fmt.Errorf("list alert rules: %w", err)
	}
	out := make([]AlertRule, 0, len(resp.Payload))
	for _, r := range resp.Payload {
		out = append(out, toAlertRule(r))
	}
	return out, nil
}

// toAlertRule 把 SDK models.AlertRule 转成内部 AlertRule 模型。
func toAlertRule(r *models.AlertRule) AlertRule {
	out := AlertRule{}
	if r.ID != nil {
		out.ID = *r.ID
	}
	if r.LocalID != nil {
		out.ID = *r.LocalID
	}
	if r.Disabled != nil {
		out.Enabled = !*r.Disabled
	}
	if r.Cluster != nil {
		if r.Cluster.ID != nil {
			out.TargetID = *r.Cluster.ID
		}
		if r.Cluster.Name != nil {
			out.TargetKind = *r.Cluster.Name
		}
	}
	return out
}
