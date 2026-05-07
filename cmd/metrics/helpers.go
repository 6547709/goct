package metrics

import (
	"fmt"
	"io"
	"time"

	"github.com/smartxworks/cloudtower-go-sdk/v2/models"

	"github.com/6547709/goct/pkg/metrics"
)

// renderMetricsResults 将 API 返回结果转换为 MetricResult 列表并渲染
func renderMetricsResults(w io.Writer, apiResults []*models.WithTaskMetric, targetName, metricName, targetType string) error {
	results := make([]metrics.MetricResult, 0)

	for _, r := range apiResults {
		if r.Data == nil {
			continue
		}

		mr := metrics.MetricResult{
			Target:     targetName,
			TargetType: targetType,
			Metric:     metricName,
		}

		// 提取 sample_streams 或 samples
		if len(r.Data.SampleStreams) > 0 {
			for _, stream := range r.Data.SampleStreams {
				if stream.Points == nil {
					continue
				}
				for _, p := range stream.Points {
					if p == nil || p.T == nil {
						continue
					}
					var value float64
					if p.V != nil {
						value = *p.V
					}
					mr.Samples = append(mr.Samples, metrics.MetricSample{
						Timestamp: time.Unix(*p.T/1000, 0).Format(time.RFC3339),
						Value:     value,
						Unit:      string(*r.Data.Unit),
					})
				}
			}
		}

		// 处理 samples 字段
		if len(r.Data.Samples) > 0 {
			for _, s := range r.Data.Samples {
				if s == nil || s.Point == nil || s.Point.T == nil {
					continue
				}
				var value float64
				if s.Point.V != nil {
					value = *s.Point.V
				}
				mr.Samples = append(mr.Samples, metrics.MetricSample{
					Timestamp: time.Unix(*s.Point.T/1000, 0).Format(time.RFC3339),
					Value:     value,
					Unit:      string(*r.Data.Unit),
				})
			}
		}

		if latestFlag && len(mr.Samples) > 0 {
			mr.Latest = &mr.Samples[len(mr.Samples)-1]
		}

		results = append(results, mr)
	}

	return metrics.RenderResult(w, results, formatFlag, latestFlag)
}

// resolveTargetArg 解析目标名称或 ID
func resolveTargetArg(args []string, envVar string) (string, error) {
	if len(args) > 0 && args[0] != "" {
		return args[0], nil
	}
	// envVar would be checked by caller if needed
	return "", fmt.Errorf("target not specified")
}
