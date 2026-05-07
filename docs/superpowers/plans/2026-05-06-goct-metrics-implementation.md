# goct Metrics 子系统实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现 goct 的 metrics 查询子命令，支持 VM/Host/Volume/Cluster/SFS 的监控指标，覆盖官方所有指标，支持 table/json/chart 三种输出格式。

**Architecture:** 在 goct 现有三层架构（cmd → service → adapter）基础上，新增 metrics 相关命令和数据结构。指标定义使用内置 JSON，API 调用直接透传给 CloudTower GetMetrics 系列接口。

**Tech Stack:** Go, Cobra, CloudTower Go SDK, 内置 JSON 指标定义

---

## 文件结构

```
cmd/metrics/                    # 新建：metrics 子命令
├── root.go                     # 父命令，定义 --list/--range/--latest/--format flags
├── vm_metrics.go               # vm.metrics（elf_*）
├── vm_volume.go                # vm.volume（zbs_volume_*）
├── host_metrics.go             # host.metrics（host_*）
├── cluster_metrics.go          # cluster.metrics（zbs_cluster_*）
├── volume_metrics.go           # volume.metrics（独立 ZBS 卷）
└── sfs_metrics.go              # sfs.metrics（sfs_*）

pkg/metrics/                    # 新建：metrics 核心逻辑
├── definitions/                # 内置指标 JSON
│   ├── vm_metrics.json         # elf_* 指标（~80条）
│   ├── host_metrics.json       # host_* 指标（~96条）
│   ├── volume_metrics.json     # zbs_volume_* 指标（~256条）
│   ├── cluster_metrics.json    # zbs_cluster_* 指标
│   └── sfs_metrics.json        # sfs_* 指标（~82条）
├── adapter.go                  # GetMetrics API 封装
├── output.go                   # 渲染（table/json/chart）
└── types.go                    # MetricQuery, MetricResult 等类型
```

---

## Task 1: 创建 pkg/metrics/types.go

**Files:**
- Create: `pkg/metrics/types.go`

- [ ] **Step 1: 写入 types.go**

```go
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
    VMMetrics     []MetricDefinition `json:"vm_metrics"`
    HostMetrics   []MetricDefinition `json:"host_metrics"`
    VolumeMetrics []MetricDefinition `json:"volume_metrics"`
    ClusterMetrics []MetricDefinition `json:"cluster_metrics"`
    SFSMetrics    []MetricDefinition `json:"sfs_metrics"`
}
```

- [ ] **Step 2: 验证语法**

Run: `cd /Users/liguoqiang/project/goct && go build ./pkg/metrics/...`
Expected: 无输出（编译通过）

---

## Task 2: 创建 pkg/metrics/output.go

**Files:**
- Create: `pkg/metrics/output.go`

- [ ] **Step 1: 写入 output.go**

```go
package metrics

import (
    "encoding/json"
    "fmt"
    "io"
    "strings"
    "time"

    "github.com/olekukonko/tablewriter"
)

// OutputFormat 支持的三种输出格式。
const (
    FormatTable = "table"
    FormatJSON  = "json"
    FormatChart = "chart"
)

// RenderResult 将 MetricResult 列表渲染为指定格式。
func RenderResult(w io.Writer, results []MetricResult, format string, latest bool) error {
    switch format {
    case FormatJSON:
        return renderJSON(w, results)
    case FormatChart:
        return renderChart(w, results)
    default:
        return renderTable(w, results, latest)
    }
}

func renderJSON(w io.Writer, results []MetricResult) error {
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    return enc.Encode(results)
}

func renderTable(w io.Writer, results []MetricResult, latest bool) error {
    tw := tablewriter.NewWriter(w)
    headers := []string{"TIME", "METRIC", "VALUE", "UNIT"}
    tw.Header(headers...)

    for _, r := range results {
        if latest && r.Latest != nil {
            tw.Append([]string{
                r.Latest.Timestamp,
                r.Metric,
                fmt.Sprintf("%.2f", r.Latest.Value),
                r.Latest.Unit,
            })
        } else {
            for _, s := range r.Samples {
                tw.Append([]string{
                    s.Timestamp,
                    r.Metric,
                    fmt.Sprintf("%.2f", s.Value),
                    s.Unit,
                })
            }
        }
    }
    tw.Render()
    return nil
}

func renderChart(w io.Writer, results []MetricResult) error {
    for _, r := range results {
        if len(r.Samples) == 0 {
            continue
        }
        fmt.Fprintf(w, "\n%s (%s)\n", r.Metric, r.Target)

        // 找最大最小值
        var min, max float64
        for i, s := range r.Samples {
            if i == 0 || s.Value < min {
                min = s.Value
            }
            if s.Value > max {
                max = s.Value
            }
        }

        // 渲染 ASCII 图表
        barWidth := 50
        range_ := max - min
        for _, s := range r.Samples {
            var barLen int
            if range_ > 0 {
                barLen = int((s.Value - min) / range_ * float64(barWidth))
            }
            bar := strings.Repeat("█", barLen)
            fmt.Fprintf(w, "%5.1f|%s%s %.2f%%\n", s.Value, bar,
                strings.Repeat(" ", barWidth-barLen), s.Value)
        }
        fmt.Fprintln(w)
    }
    return nil
}
```

- [ ] **Step 2: 验证语法**

Run: `cd /Users/liguoqiang/project/goct && go build ./pkg/metrics/...`
Expected: 无输出（编译通过）

---

## Task 3: 创建 pkg/metrics/adapter.go

**Files:**
- Create: `pkg/metrics/adapter.go`

- [ ] **Step 1: 检查 SDK 中的 GetMetrics 相关模型**

先检查 SDK 中的模型定义，了解 GetVmMetricInput 等结构。

Run: `ls /Users/liguoqiang/project/goct/cloudtower-go-sdk/models/ | grep -i metric`
Expected: 列出 get_vm_metric_input.go 等文件

- [ ] **Step 2: 写入 adapter.go**

```go
package metrics

import (
    "context"

    "github.com/smartxworks/cloudtower-go-sdk/v2/models"

    "github.com/6547709/goct/pkg/adapter"
)

// MetricsOps 定义指标查询接口。
type MetricsOps interface {
    GetVMMetrics(ctx context.Context, input *models.GetVmMetricInput) ([]models.WithTask_Metric_, error)
    GetHostMetrics(ctx context.Context, input *models.GetHostMetricInput) ([]models.WithTask_Metric_, error)
    GetVmVolumeMetrics(ctx context.Context, input *models.GetVmVolumeMetricInput) ([]models.WithTask_Metric_, error)
    GetClusterMetrics(ctx context.Context, input *models.GetClusterMetricInput) ([]models.WithTask_Metric_, error)
}

// NewMetricsClient 创建 metrics 操作客户端。
func NewMetricsClient(c adapter.Client) MetricsOps {
    return &metricsClient{client: c}
}

type metricsClient struct {
    client adapter.Client
}

// GetVMMetrics 调用 GetVmMetrics API。
func (m *metricsClient) GetVMMetrics(ctx context.Context, input *models.GetVmMetricInput) ([]models.WithTask_Metric_, error) {
    return m.client.GetVMMetrics(ctx, input)
}

// GetHostMetrics 调用 GetHostMetrics API。
func (m *metricsClient) GetHostMetrics(ctx context.Context, input *models.GetHostMetricInput) ([]models.WithTask_Metric_, error) {
    return m.client.GetHostMetrics(ctx, input)
}

// GetVmVolumeMetrics 调用 GetVmVolumeMetrics API。
func (m *metricsClient) GetVmVolumeMetrics(ctx context.Context, input *models.GetVmVolumeMetricInput) ([]models.WithTask_Metric_, error) {
    return m.client.GetVmVolumeMetrics(ctx, input)
}

// GetClusterMetrics 调用 GetClusterMetrics API。
func (m *metricsClient) GetClusterMetrics(ctx context.Context, input *models.GetClusterMetricInput) ([]models.WithTask_Metric_, error) {
    return m.client.GetClusterMetrics(ctx, input)
}
```

- [ ] **Step 3: 检查 adapter.Client 是否已有这些方法**

Run: `grep -n "GetVMMetrics\|GetHostMetrics\|GetVmVolumeMetrics\|GetClusterMetrics" /Users/liguoqiang/project/goct/pkg/adapter/*.go`
Expected: 可能需要添加这些方法到 adapter

---

## Task 4: 添加 GetMetrics 方法到 adapter.Client

**Files:**
- Modify: `pkg/adapter/client.go`（添加方法签名）
- Modify: `pkg/adapter/vm.go`（实现 GetVMMetrics）
- Create: `pkg/adapter/metrics.go`（实现其他 GetMetrics）

- [ ] **Step 1: 检查 adapter.Client 接口定义**

Run: `grep -n "type Client interface" -A 30 /Users/liguoqiang/project/goct/pkg/adapter/client.go`
Expected: Client 接口定义

- [ ] **Step 2: 添加 MetricsOps 到 Client 接口**

修改 client.go，在 Client 接口中添加 MetricsOps 内嵌或显式方法。

- [ ] **Step 3: 在 vm.go 中添加 GetVMMetrics 实现**

Run: `grep -n "func.*GetVM\|GetVMMetrics" /Users/liguoqiang/project/goct/pkg/adapter/vm.go`
Expected: 查看现有 VM 相关方法

- [ ] **Step 4: 创建 pkg/adapter/metrics.go 实现 GetMetrics 系列**

```go
package adapter

import (
    "context"

    "github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func (c *defaultClient) GetVMMetrics(ctx context.Context, input *models.GetVmMetricInput) ([]models.WithTask_Metric_, error) {
    resp, err := c.api.Metrics.GetVmMetrics(nil, input, c.transport)
    if err != nil {
        return nil, err
    }
    return resp.Payload, nil
}

func (c *defaultClient) GetHostMetrics(ctx context.Context, input *models.GetHostMetricInput) ([]models.WithTask_Metric_, error) {
    resp, err := c.api.Metrics.GetHostMetrics(nil, input, c.transport)
    if err != nil {
        return nil, err
    }
    return resp.Payload, nil
}

func (c *defaultClient) GetVmVolumeMetrics(ctx context.Context, input *models.GetVmVolumeMetricInput) ([]models.WithTask_Metric_, error) {
    resp, err := c.api.Metrics.GetVmVolumeMetrics(nil, input, c.transport)
    if err != nil {
        return nil, err
    }
    return resp.Payload, nil
}

func (c *defaultClient) GetClusterMetrics(ctx context.Context, input *models.GetClusterMetricInput) ([]models.WithTask_Metric_, error) {
    resp, err := c.api.Metrics.GetClusterMetrics(nil, input, c.transport)
    if err != nil {
        return nil, err
    }
    return resp.Payload, nil
}
```

- [ ] **Step 5: 验证编译**

Run: `cd /Users/liguoqiang/project/goct && go build ./pkg/adapter/...`
Expected: 无输出或错误（根据 SDK 方法名可能需要调整）

---

## Task 5: 创建指标定义 JSON 文件

**Files:**
- Create: `pkg/metrics/definitions/vm_metrics.json`
- Create: `pkg/metrics/definitions/host_metrics.json`
- Create: `pkg/metrics/definitions/volume_metrics.json`
- Create: `pkg/metrics/definitions/cluster_metrics.json`
- Create: `pkg/metrics/definitions/sfs_metrics.json`

- [ ] **Step 1: 分析 metrics-lookup 的 markdown 文件，提取指标**

转换 metrics_host.md 中的表格为 JSON 格式。

- [ ] **Step 2: 创建 vm_metrics.json（~80 条 elf_* 指标）**

```json
{
  "vm_metrics": [
    {"name": "elf_vm_cpu_overall_usage_percent", "description": "VM CPU 使用率", "version": "v5.1.0+"},
    {"name": "elf_vm_memory_usage_percent", "description": "VM 内存使用率", "version": "v5.1.0+"}
  ]
}
```

- [ ] **Step 3: 创建 host_metrics.json（~96 条 host_* 指标）**

- [ ] **Step 4: 创建 volume_metrics.json（~256 条 zbs_volume_* 指标）**

- [ ] **Step 5: 创建 cluster_metrics.json（zbs_cluster_* 指标）**

- [ ] **Step 6: 创建 sfs_metrics.json（~82 条 sfs_* 指标）**

---

## Task 6: 创建 cmd/metrics/root.go

**Files:**
- Create: `cmd/metrics/root.go`

- [ ] **Step 1: 创建目录并写入 root.go**

```go
package metrics

import (
    "github.com/spf13/cobra"
)

var (
    listFlag   bool
    latestFlag bool
    rangeFlag  string
    formatFlag string
)

var rootCmd = &cobra.Command{
    Use:   "metrics",
    Short: "Query CloudTower metrics",
    Long:  "Query VM, Host, Volume, Cluster and SFS metrics",
}

func Register(root *cobra.Command) {
    root.AddCommand(rootCmd)

    rootCmd.PersistentFlags().BoolVar(&listFlag, "list", false, "List available metrics")
    rootCmd.PersistentFlags().StringVar(&rangeFlag, "range", "5m", "Time range: 5m, 1h, 1d, 7d")
    rootCmd.PersistentFlags().BoolVar(&latestFlag, "latest", false, "Show only latest value")
    rootCmd.PersistentFlags().StringVar(&formatFlag, "format", "table", "Output format: table, json, chart")
}
```

---

## Task 7: 创建 cmd/metrics/vm_metrics.go

**Files:**
- Create: `cmd/metrics/vm_metrics.go`

- [ ] **Step 1: 写入 vm_metrics.go**

```go
package metrics

import (
    "github.com/6547709/goct/pkg/client"
    "github.com/6547709/goct/pkg/metrics"
    "github.com/spf13/cobra"
)

func newVMMetrics() *cobra.Command {
    var args struct {
        metric  string
        vmName  string
    }
    c := &cobra.Command{
        Use: "vm.metrics <metric> [vm-name]",
        Short: "Query VM metrics (elf_*)",
        Args: cobra.RangeArgs(1, 2),
        RunE: func(cmd *cobra.Command, args []string) error {
            cli := client.From(cmd.Context())
            metricClient := metrics.NewMetricsClient(cli)

            target := args[1] // vm name or id
            query := metrics.MetricQuery{
                Target:     target,
                TargetType: "vm",
                Metrics:    []string{args[0]},
                Range:      rangeFlag,
                Latest:     latestFlag,
                Format:     formatFlag,
            }

            results, err := metricClient.QueryVMMetrics(cmd.Context(), query)
            if err != nil {
                return err
            }

            return metrics.RenderResult(cmd.OutOrStdout(), results, formatFlag, latestFlag)
        },
    }
    return c
}
```

- [ ] **Step 2: 实现 pkg/metrics/client.go 的 QueryVMMetrics 方法**

在 pkg/metrics/adapter.go 或新建 client.go 中实现查询逻辑。

---

## Task 8: 创建 host.metrics 和其他子命令

**Files:**
- Create: `cmd/metrics/host_metrics.go`
- Create: `cmd/metrics/vm_volume.go`
- Create: `cmd/metrics/cluster_metrics.go`
- Create: `cmd/metrics/volume_metrics.go`
- Create: `cmd/metrics/sfs_metrics.go`

按 vm_metrics.go 相同的模式实现。

---

## Task 9: 注册命令到 root

**Files:**
- Modify: `cmd/root.go`

- [ ] **Step 1: 添加 metrics 包导入和注册**

```go
import "github.com/6547709/goct/cmd/metrics"

func init() {
    // ... existing commands ...
    metrics.Register(rootCmd)
}
```

- [ ] **Step 2: 验证编译**

Run: `cd /Users/liguoqiang/project/goct && go build -o goct .`
Expected: 编译成功

---

## Task 10: 测试接口

**Files:**
- Test: 实际调用 CloudTower API

- [ ] **Step 1: 设置环境变量**

```bash
export GOCT_URL=https://your-tower.example.com
export GOCT_USERNAME=admin
export GOCT_PASSWORD=your-password
export GOCT_INSECURE=true
```

- [ ] **Step 2: 测试 --list**

Run: `./goct vm.metrics --list | head -20`
Expected: 列出 VM 可用指标

- [ ] **Step 3: 测试实际查询**

Run: `./goct vm.metrics elf_vm_cpu_overall_usage_percent your-vm --range 5m`
Expected: 返回时序表格

- [ ] **Step 4: 测试 JSON 输出**

Run: `./goct vm.metrics elf_vm_cpu_overall_usage_percent your-vm --range 5m --format json`
Expected: JSON 格式输出

- [ ] **Step 5: 测试 ASCII chart**

Run: `./goct vm.metrics elf_vm_cpu_overall_usage_percent your-vm --range 1h --format chart`
Expected: ASCII 图表渲染

---

## 注意事项

1. **SDK 方法名**：CloudTower Go SDK 的 GetMetrics 方法名可能与设计文档中的略有不同，需要根据实际 SDK 代码调整。
2. **指标过滤**：vm.volume 需要通过 VM 找到关联的 Volume 再查询 ZBS 指标。
3. **SFS API**：需要确认 CloudTower 是否有 SFS 专用的 Metrics API，或者需要其他方式获取。
4. **时间范围格式**：CloudTower API 可能使用不同的 range 格式（如 "5m" vs "300" 秒），需要验证。

---

## 提交

- [ ] 提交所有更改

```bash
git add -A
git commit -m "feat: add metrics subsystem for VM/Host/Volume/Cluster/SFS

- cmd/metrics: add vm.metrics, host.metrics, vm.volume, cluster.metrics, volume.metrics, sfs.metrics
- pkg/metrics: add types, output (table/json/chart), adapter
- pkg/metrics/definitions: add embedded metric definitions JSON
- pkg/adapter: add GetVMMetrics, GetHostMetrics, GetVmVolumeMetrics, GetClusterMetrics"
```