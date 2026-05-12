package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/deploy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// DeployOps 定义部署操作。
type DeployOps interface {
	ListDeploys(ctx context.Context, opts ListOpts) ([]Deploy, error)
}

// ---------- ListDeploys ----------

func (c *defaultClient) ListDeploys(ctx context.Context, opts ListOpts) ([]Deploy, error) {
	params := deploy.NewGetDeploysParams()
	params.SetContext(ctx)

	body := &models.GetDeploysRequestBody{}
	if opts.Limit > 0 {
		body.First = pointy.Int32(opts.Limit)
	}
	if opts.Skip > 0 {
		body.Skip = pointy.Int32(opts.Skip)
	}
	params.SetRequestBody(body)

	resp, err := c.api.Deploy.GetDeploys(params)
	if err != nil {
		return nil, fmt.Errorf("list deploys: %w", err)
	}
	out := make([]Deploy, 0, len(resp.Payload))
	for _, d := range resp.Payload {
		out = append(out, toDeploy(d))
	}
	return out, nil
}

// toDeploy 把 SDK models.Deploy 转成内部 Deploy 模型。
func toDeploy(d *models.Deploy) Deploy {
	out := Deploy{}
	if d.ID != nil {
		out.ID = *d.ID
	}
	if d.Version != nil {
		out.VMID = *d.Version
	}
	return out
}
