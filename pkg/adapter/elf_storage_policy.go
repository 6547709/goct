package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/elf_storage_policy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// ElfStoragePolicyOps 定义存储策略操作。
type ElfStoragePolicyOps interface {
	ListElfStoragePolicies(ctx context.Context, opts ListOpts) ([]ElfStoragePolicy, error)
}

// ---------- ListElfStoragePolicies ----------

func (c *defaultClient) ListElfStoragePolicies(ctx context.Context, opts ListOpts) ([]ElfStoragePolicy, error) {
	params := elf_storage_policy.NewGetElfStoragePoliciesParams()
	params.SetContext(ctx)

	where := &models.ElfStoragePolicyWhereInput{}
	hasWhere := false
	if opts.NameContains != "" {
		where.NameContains = pointy.String(opts.NameContains)
		hasWhere = true
	}
	if opts.ClusterID != "" {
		where.Cluster = &models.ClusterWhereInput{ID: pointy.String(opts.ClusterID)}
		hasWhere = true
	}

	body := &models.GetElfStoragePoliciesRequestBody{}
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

	resp, err := c.api.ElfStoragePolicy.GetElfStoragePolicies(params)
	if err != nil {
		return nil, fmt.Errorf("list elf storage policies: %w", err)
	}
	out := make([]ElfStoragePolicy, 0, len(resp.Payload))
	for _, p := range resp.Payload {
		out = append(out, toElfStoragePolicy(p))
	}
	return out, nil
}

// toElfStoragePolicy 把 SDK models.ElfStoragePolicy 转成内部 ElfStoragePolicy 模型。
func toElfStoragePolicy(p *models.ElfStoragePolicy) ElfStoragePolicy {
	out := ElfStoragePolicy{}
	if p.ID != nil {
		out.ID = *p.ID
	}
	if p.Name != nil {
		out.Name = *p.Name
	}
	if p.Description != nil {
		out.Description = *p.Description
	}
	if p.Cluster != nil {
		if p.Cluster.ID != nil {
			out.ClusterID = *p.Cluster.ID
		}
		if p.Cluster.Name != nil {
			out.ClusterName = *p.Cluster.Name
		}
	}
	if p.LocalID != nil {
		out.LocalID = *p.LocalID
	}
	if p.ReplicaNum != nil {
		out.ReplicaNum = *p.ReplicaNum
	}
	if p.StripeNum != nil {
		out.StripeNum = *p.StripeNum
	}
	if p.StripeSize != nil {
		out.StripeSize = *p.StripeSize
	}
	if p.ThinProvision != nil {
		out.ThinProvision = *p.ThinProvision
	}
	return out
}
