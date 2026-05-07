package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	sdkcluster "github.com/smartxworks/cloudtower-go-sdk/v2/client/cluster"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// ClusterOps 定义集群操作。
type ClusterOps interface {
	ListClusters(ctx context.Context, opts ListOpts) ([]Cluster, error)
	GetCluster(ctx context.Context, id string) (*Cluster, error)
	GetClusterByName(ctx context.Context, name string) (*Cluster, error)
}

func (c *defaultClient) ListClusters(ctx context.Context, opts ListOpts) ([]Cluster, error) {
	params := sdkcluster.NewGetClustersParams()
	params.SetContext(ctx)
	body := &models.GetClustersRequestBody{}
	where := &models.ClusterWhereInput{}
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
	resp, err := c.api.Cluster.GetClusters(params)
	if err != nil {
		return nil, fmt.Errorf("list clusters: %w", err)
	}
	out := make([]Cluster, 0, len(resp.Payload))
	for _, cl := range resp.Payload {
		out = append(out, toCluster(cl))
	}
	return out, nil
}

func (c *defaultClient) GetCluster(ctx context.Context, id string) (*Cluster, error) {
	params := sdkcluster.NewGetClustersParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetClustersRequestBody{
		Where: &models.ClusterWhereInput{ID: pointy.String(id)},
	})
	resp, err := c.api.Cluster.GetClusters(params)
	if err != nil {
		return nil, fmt.Errorf("get cluster %s: %w", id, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get cluster %s: %w", id, ErrNotFound)
	}
	cl := toCluster(resp.Payload[0])
	return &cl, nil
}

func (c *defaultClient) GetClusterByName(ctx context.Context, name string) (*Cluster, error) {
	params := sdkcluster.NewGetClustersParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetClustersRequestBody{
		Where: &models.ClusterWhereInput{Name: &name},
	})
	resp, err := c.api.Cluster.GetClusters(params)
	if err != nil {
		return nil, fmt.Errorf("get cluster by name %s: %w", name, err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("get cluster by name %s: %w", name, ErrNotFound)
	}
	cl := toCluster(resp.Payload[0])
	return &cl, nil
}

func toCluster(cl *models.Cluster) Cluster {
	out := Cluster{}
	if cl.ID != nil {
		out.ID = *cl.ID
	}
	if cl.Name != nil {
		out.Name = *cl.Name
	}
	if cl.TotalMemoryBytes != nil {
		out.TotalMemoryBytes = uint64(*cl.TotalMemoryBytes)
	}
	if cl.TotalDataCapacity != nil {
		out.TotalDataCapacity = uint64(*cl.TotalDataCapacity)
	}
	if cl.UsedDataSpace != nil {
		out.UsedDataSpace = uint64(*cl.UsedDataSpace)
	}
	if cl.TotalCPUCores != nil {
		out.TotalCPUCores = *cl.TotalCPUCores
	}
	if cl.RunningVMNum != nil {
		out.RunningVMs = *cl.RunningVMNum
	}
	return out
}
