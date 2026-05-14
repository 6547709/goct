package adapter

import (
	"context"
	"fmt"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/client/cluster_settings"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// ClusterSettingsOps 定义集群设置操作。
type ClusterSettingsOps interface {
	GetClusterSettings(ctx context.Context, clusterID string) (*ClusterSettings, error)
}

// ---------- GetClusterSettings ----------

func (c *defaultClient) GetClusterSettings(ctx context.Context, clusterID string) (*ClusterSettings, error) {
	params := cluster_settings.NewGetClusterSettingsesParams()
	params.SetContext(ctx)
	params.SetRequestBody(&models.GetClusterSettingsesRequestBody{
		First: pointy.Int32(1),
		Where: &models.ClusterSettingsWhereInput{
			Cluster: &models.ClusterWhereInput{ID: pointy.String(clusterID)},
		},
	})

	resp, err := c.api.ClusterSettings.GetClusterSettingses(params)
	if err != nil {
		return nil, fmt.Errorf("get cluster settings: %w", err)
	}
	if len(resp.Payload) == 0 {
		return nil, fmt.Errorf("cluster settings not found")
	}
	return toClusterSettings(resp.Payload[0]), nil
}

// toClusterSettings 把 SDK models.ClusterSettings 转成内部 ClusterSettings 模型。
func toClusterSettings(s *models.ClusterSettings) *ClusterSettings {
	out := &ClusterSettings{}
	if s.ID != nil {
		out.ID = *s.ID
	}
	if s.Cluster != nil && s.Cluster.ID != nil {
		out.ClusterID = *s.Cluster.ID
	}
	if s.DefaultStoragePolicy != nil {
		out.DefaultStoragePolicy = string(*s.DefaultStoragePolicy)
	}
	return out
}
