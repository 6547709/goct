package metrics

import (
	"context"

	"github.com/smartxworks/cloudtower-go-sdk/v2/models"

	"github.com/6547709/goct/pkg/adapter"
)

// MetricsOps 定义指标查询接口。
type MetricsOps interface {
	GetVMMetrics(ctx context.Context, input *models.GetVMMetricInput) ([]*models.WithTaskMetric, error)
	GetHostMetrics(ctx context.Context, input *models.GetHostMetricInput) ([]*models.WithTaskMetric, error)
	GetVmVolumeMetrics(ctx context.Context, input *models.GetVMVolumeMetricInput) ([]*models.WithTaskMetric, error)
	GetClusterMetrics(ctx context.Context, input *models.GetClusterMetricInput) ([]*models.WithTaskMetric, error)
}

// NewMetricsClient 创建 metrics 操作客户端。
func NewMetricsClient(c adapter.Client) MetricsOps {
	return &metricsClient{client: c}
}

type metricsClient struct {
	client adapter.Client
}

func (m *metricsClient) GetVMMetrics(ctx context.Context, input *models.GetVMMetricInput) ([]*models.WithTaskMetric, error) {
	return m.client.GetVMMetrics(ctx, input)
}

func (m *metricsClient) GetHostMetrics(ctx context.Context, input *models.GetHostMetricInput) ([]*models.WithTaskMetric, error) {
	return m.client.GetHostMetrics(ctx, input)
}

func (m *metricsClient) GetVmVolumeMetrics(ctx context.Context, input *models.GetVMVolumeMetricInput) ([]*models.WithTaskMetric, error) {
	return m.client.GetVmVolumeMetrics(ctx, input)
}

func (m *metricsClient) GetClusterMetrics(ctx context.Context, input *models.GetClusterMetricInput) ([]*models.WithTaskMetric, error) {
	return m.client.GetClusterMetrics(ctx, input)
}