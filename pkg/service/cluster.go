package service

import (
	"context"
	"github.com/6547709/goct/pkg/adapter"
)

type ClusterService struct{ c adapter.ClusterOps }
func NewCluster(c adapter.ClusterOps) *ClusterService { return &ClusterService{c: c} }

func (s *ClusterService) List(ctx context.Context, opts adapter.ListOpts) ([]adapter.Cluster, error) {
	return s.c.ListClusters(ctx, opts)
}

func (s *ClusterService) Resolve(ctx context.Context, idOrName string) (*adapter.Cluster, error) {
	return Resolve(ctx, s.c.ListClusters, s.c.GetCluster,
		func(c adapter.Cluster) (string, string) { return c.ID, c.Name },
		idOrName)
}
