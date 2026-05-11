package metrics

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

//go:embed definitions/*.json
var definitionsFS embed.FS

// listMetric 定义列表项
type listMetric struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version,omitempty"`
}

type listMetricsData struct {
	VMMetrics      []listMetric `json:"vm_metrics"`
	HostMetrics    []listMetric `json:"host_metrics"`
	VolumeMetrics  []listMetric `json:"volume_metrics"`
	ClusterMetrics []listMetric `json:"cluster_metrics"`
	SFSMetrics     []listMetric `json:"sfs_metrics"`
}

// ListMetrics lists all available metrics from embedded JSON
func ListMetrics(w io.Writer, metricType string) error {
	var data []listMetric
	var err error

	switch metricType {
	case "vm":
		data, err = loadMetrics("vm_metrics")
	case "host":
		data, err = loadMetrics("host_metrics")
	case "volume":
		data, err = loadMetrics("volume_metrics")
	case "cluster":
		data, err = loadMetrics("cluster_metrics")
	case "sfs":
		data, err = loadMetrics("sfs_metrics")
	default:
		data, err = loadMetrics(metricType)
	}

	if err != nil {
		return err
	}

	// 输出表格
	tw := tablewriter.NewWriter(w)
	tw.Header([]string{"NAME", "DESCRIPTION", "VERSION"})

	for _, m := range data {
		version := m.Version
		if version == "" {
			version = "-"
		}
		tw.Append([]string{m.Name, m.Description, version})
	}
	tw.Render()
	return nil
}

func loadMetrics(name string) ([]listMetric, error) {
	filename := fmt.Sprintf("definitions/%s.json", name)

	// 尝试从 embed.FS 读取
	data, err := definitionsFS.ReadFile(filename)
	if err != nil {
		// 如果 embed 失败，尝试从文件系统读取
		exePath, _ := os.Executable()
		fullPath := filepath.Join(filepath.Dir(exePath), "pkg", "metrics", filename)
		data, err = os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load %s: %w", filename, err)
		}
	}

	var result map[string][]listMetric
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	// 找到对应的 key
	for key, value := range result {
		if strings.Contains(key, name) || strings.HasPrefix(key, name) {
			return value, nil
		}
	}

	return nil, fmt.Errorf("metric type %s not found", name)
}

// GetMetricNames returns only metric names for completion
func GetMetricNames(metricType string) ([]string, error) {
	data, err := loadMetrics(metricType)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(data))
	for i, m := range data {
		names[i] = m.Name
	}
	return names, nil
}