package metrics

// MetricQuery 描述一次指标查询请求。
type MetricQuery struct {
    Target     string   // 查询对象名称或 ID（VM/Host/Volume 等）
    TargetType string   // vm|host|volume|cluster|sfs
    Metrics    []string // 指标名列表
    Range      string   // 时间范围，如 "5m", "1h", "1d"
    Latest     bool     // 是否只返回最新值
    Format     string   // table|json|chart
}

// MetricSample 表示一个数据点。
type MetricSample struct {
    Timestamp  string  `json:"timestamp"`
    Value     float64 `json:"value"`
    Unit      string  `json:"unit"`
}

// MetricResult 描述查询结果。
type MetricResult struct {
    Target     string         `json:"target"`
    TargetType string         `json:"target_type"`
    Metric     string         `json:"metric"`
    Samples    []MetricSample `json:"samples,omitempty"`
    Latest     *MetricSample  `json:"latest,omitempty"`
    Error      string         `json:"error,omitempty"`
}

// MetricDefinition 指标定义（用于 --list）。
type MetricDefinition struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Version     string `json:"version,omitempty"`
}

// MetricDefinitions 各类指标的集合。
type MetricDefinitions struct {
    VMMetrics      []MetricDefinition `json:"vm_metrics"`
    HostMetrics    []MetricDefinition `json:"host_metrics"`
    VolumeMetrics  []MetricDefinition `json:"volume_metrics"`
    ClusterMetrics []MetricDefinition `json:"cluster_metrics"`
    SFSMetrics     []MetricDefinition `json:"sfs_metrics"`
}