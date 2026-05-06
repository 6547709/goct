package adapter

import (
	"context"

	"github.com/smartxworks/cloudtower-go-sdk/v2/client/metrics"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// MetricsOps 定义 metrics 相关的 SDK 操作。
type MetricsOps interface {
	GetVMMetrics(ctx context.Context, input *models.GetVMMetricInput) ([]*models.WithTaskMetric, error)
	GetHostMetrics(ctx context.Context, input *models.GetHostMetricInput) ([]*models.WithTaskMetric, error)
	GetVmVolumeMetrics(ctx context.Context, input *models.GetVMVolumeMetricInput) ([]*models.WithTaskMetric, error)
	GetClusterMetrics(ctx context.Context, input *models.GetClusterMetricInput) ([]*models.WithTaskMetric, error)
}

func (c *defaultClient) GetVMMetrics(ctx context.Context, input *models.GetVMMetricInput) ([]*models.WithTaskMetric, error) {
	params := metrics.NewGetVMMetricsParams()
	params.SetContext(ctx)
	params.SetRequestBody(input)
	resp, err := c.api.Metrics.GetVMMetrics(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *defaultClient) GetHostMetrics(ctx context.Context, input *models.GetHostMetricInput) ([]*models.WithTaskMetric, error) {
	params := metrics.NewGetHostMetricsParams()
	params.SetContext(ctx)
	params.SetRequestBody(input)
	resp, err := c.api.Metrics.GetHostMetrics(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *defaultClient) GetVmVolumeMetrics(ctx context.Context, input *models.GetVMVolumeMetricInput) ([]*models.WithTaskMetric, error) {
	params := metrics.NewGetVMVolumeMetricsParams()
	params.SetContext(ctx)
	params.SetRequestBody(input)
	resp, err := c.api.Metrics.GetVMVolumeMetrics(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (c *defaultClient) GetClusterMetrics(ctx context.Context, input *models.GetClusterMetricInput) ([]*models.WithTaskMetric, error) {
	params := metrics.NewGetClusterMetricsParams()
	params.SetContext(ctx)
	params.SetRequestBody(input)
	resp, err := c.api.Metrics.GetClusterMetrics(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}