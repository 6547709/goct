package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/application"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// ApplicationOps 定义应用操作。
type ApplicationOps interface {
	ListApplications(ctx context.Context, opts ListOpts) ([]Application, error)
}

// ---------- ListApplications ----------

func (c *defaultClient) ListApplications(ctx context.Context, opts ListOpts) ([]Application, error) {
	params := application.NewGetApplicationsParams()
	params.SetContext(ctx)

	where := &models.ApplicationWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.ImageNameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if opts.ClusterID != "" {
		where.Cluster = &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}
		hasWhere = true
	}

	body := &models.GetApplicationsRequestBody{}
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

	resp, err := c.api.Application.GetApplications(params)
	if err != nil {
		return nil, fmt.Errorf("list applications: %w", err)
	}
	out := make([]Application, 0, len(resp.Payload))
	for _, a := range resp.Payload {
		out = append(out, toApplication(a))
	}
	return out, nil
}

// toApplication 把 SDK models.Application 转成内部 Application 模型。
func toApplication(a *models.Application) Application {
	out := Application{}
	if a.ID != nil {
		out.ID = *a.ID
	}
	if a.LocalID != nil {
		out.LocalID = *a.LocalID
	}
	if a.ImageName != nil {
		out.ImageName = *a.ImageName
	}
	if a.Memory != nil {
		out.Memory = *a.Memory
	}
	if a.State != nil {
		out.State = string(*a.State)
	}
	if a.StorageIP != nil {
		out.StorageIP = *a.StorageIP
	}
	if a.ErrorMessage != nil {
		out.ErrorMessage = *a.ErrorMessage
	}
	if a.Cluster != nil {
		if a.Cluster.ID != nil {
			out.ClusterID = *a.Cluster.ID
		}
		if a.Cluster.Name != nil {
			out.ClusterName = *a.Cluster.Name
		}
	}
	return out
}
