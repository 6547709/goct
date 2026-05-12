package service

import (
	"context"

	"github.com/6547709/goct/pkg/adapter"
)

type ClusterSettingsService struct{ c adapter.ClusterSettingsOps }
func NewClusterSettings(c adapter.ClusterSettingsOps) *ClusterSettingsService { return &ClusterSettingsService{c: c} }

func (s *ClusterSettingsService) GetSettings(ctx context.Context, clusterID string) (*adapter.ClusterSettings, error) {
	return s.c.GetClusterSettings(ctx, clusterID)
}