# goct Metrics 子系统设计

日期: 2026-05-06
状态: draft

## 概述

goct 增加 metrics 查询能力，支持 VM、Host、Volume (ZBS)、Cluster、 SFS 的监控指标查询。
- AI 调度用 JSON 输出
- 人类查看支持 table 和 ASCII chart

## 设计原则

1. **最大覆盖** - 覆盖官方所有 VM/Host/Volume 指标，不过滤
2. **AI 优先** - JSON 输出，结构清晰，便于 AI 解析
3. **人类友好** - table 默认输出，chart 模式 ASCII 可视化
4. **自包含** - 内置指标定义 JSON，不依赖外部文件

## 命令结构

```
goct vm.metrics <metric> [vm-name]        # VM CPU/内存/磁盘 (elf_*)
goct vm.volume <metric> [vm-name]          # VM 关联的 ZBS 卷 (zbs_volume_*)
goct host.metrics <metric> [host-name]     # Host (host_*)
goct cluster.metrics <metric>              # Cluster 存储容量 (zbs_cluster_*)
goct volume.metrics <metric> [volume-name] # 独立 ZBS 卷
goct sfs.metrics <metric> [sfs-name]      # SFS 文件存储 (sfs_*)
```

### 全局 Flags

| Flag | 说明 | 默认值 |
|------|------|--------|
| `--list` | 列出所有可用指标（带中文说明） | - |
| `--latest` | 只输出最新值 | false（输出时序表格） |
| `--range` | 时间范围：5m, 1h, 1d, 7d | 5m |
| `--format` | 输出格式: table, json, chart | table |

## 指标定义

内置 JSON 文件，位于 `pkg/metrics/definitions/`:

```json
{
  "vm_metrics": [
    {"name": "elf_vm_cpu_overall_usage_percent", "description": "VM CPU 使用率", "version": "v5.1.0+"}
  ],
  "host_metrics": [...],
  "volume_metrics": [...],
  "cluster_metrics": [...],
  "sfs_metrics": [...]
}
```

数据来源：转换自 cloudtower-skills/skills/metrics-lookup/references/

## API 调用

直接透传给 CloudTower GetMetrics 系列 API:

- `GetVmMetrics` → elf_* 指标
- `GetHostMetrics` → host_* 指标
- `GetVmVolumeMetrics` → zbs_volume_* 指标
- `GetClusterMetrics` → zbs_cluster_* 指标
- `GetVmNetWorkMetrics` → elf_* 网络指标

### 输入参数

```go
type MetricQuery struct {
    Target  string   // VM/Host/Cluster name 或 ID（位置参数）
    Metrics []string // 指标名列表（位置参数，多个用逗号分隔）
    Range   string   // --range，默认 "5m"
    Latest  bool     // --latest
    Format  string   // --format: table|json|chart
}
```

## 输出格式

### table（默认）

```
TIME                METRIC                              VALUE    UNIT
2026-05-06 10:00:00 elf_vm_cpu_overall_usage_percent   45.2     PERCENT
2026-05-06 10:01:00 elf_vm_cpu_overall_usage_percent   43.8     PERCENT
2026-05-06 10:02:00 elf_vm_cpu_overall_usage_percent   44.1     PERCENT
```

### --latest 单行

```
elf_vm_cpu_overall_usage_percent  45.2  PERCENT  2026-05-06 10:02:00
```

### --format json

```json
{
  "target": "my-vm",
  "target_type": "vm",
  "metrics": ["elf_vm_cpu_overall_usage_percent"],
  "range": "5m",
  "samples": [
    {"timestamp": "2026-05-06T10:00:00Z", "value": 45.2, "unit": "PERCENT"},
    {"timestamp": "2026-05-06T10:01:00Z", "value": 43.8, "unit": "PERCENT"}
  ]
}
```

### --format chart (ASCII)

```
elf_vm_cpu_overall_usage_percent (my-vm)
50|████████████████████▒▒▒▒▒▒▒▒▒▒▒▒
45|███████████████████████▒▒▒▒▒▒▒▒▒▒ 45.2%
40|███████████████████████████████▒▒
35|███████████████████████████████▒▒
   └──────────────────────────────────
   10:00    10:30    11:00    11:30
```

## 实现计划

### Phase 1: 核心实现

1. 创建 `cmd/metrics/` 目录结构
2. 实现 `pkg/metrics/definitions/` 内置指标 JSON
3. 实现各子命令 (vm.metrics, host.metrics, vm.volume, cluster.metrics, volume.metrics, sfs.metrics)
4. 实现 GetMetrics API 适配

### Phase 2: 输出格式

1. table 输出（默认）
2. json 输出
3. chart 输出（ASCII）

### Phase 3: --list 功能

1. 读取内置 JSON
2. 按产品/子系统过滤
3. 展示中文说明

## 文件结构

```
cmd/
├── metrics/
│   ├── root.go           # 父命令 metrics
│   ├── vm_metrics.go     # vm.metrics
│   ├── vm_volume.go      # vm.volume
│   ├── host_metrics.go   # host.metrics
│   ├── cluster_metrics.go # cluster.metrics
│   ├── volume_metrics.go  # volume.metrics
│   └── sfs_metrics.go    # sfs.metrics
pkg/
├── metrics/
│   ├── definitions/      # 内置指标 JSON
│   │   ├── vm_metrics.json
│   │   ├── host_metrics.json
│   │   ├── volume_metrics.json
│   │   ├── cluster_metrics.json
│   │   └── sfs_metrics.json
│   ├── adapter/          # SDK 调用封装
│   ├── output/           # 渲染 (table/json/chart)
│   └── types.go         # MetricQuery, MetricResult
```

## API 对应关系

| 命令 | API | 指标前缀 |
|------|-----|----------|
| vm.metrics | GetVmMetrics | elf_* |
| vm.volume | GetVmVolumeMetrics | zbs_volume_* |
| host.metrics | GetHostMetrics | host_* |
| cluster.metrics | GetClusterMetrics | zbs_cluster_* |
| volume.metrics | GetVmVolumeMetrics | zbs_volume_* (按 volume name 查询) |
| sfs.metrics | (待确认) | sfs_* |

## 测试场景

### AI 调度场景

```bash
# 发现 VM 异常，查 CPU 和内存
goct vm.metrics elf_vm_cpu_overall_usage_percent,elf_vm_memory_usage_percent my-vm --range 10m --format json

# 解释告警，查磁盘延迟
goct vm.volume zbs_volume_read_latency_ns,zbs_volume_write_latency_ns my-vm --latest --format json

# 查 SFS 存储
goct sfs.metrics sfs_bytes_received_by_export_total my-sfs --range 1h --format json
```

### 人类巡检场景

```bash
# 列出所有可用 VM 指标
goct vm.metrics --list

# 图表查看
goct vm.metrics elf_vm_cpu_overall_usage_percent my-vm --format chart

# 最新值单行
goct host.metrics host_memory_usage_percent,host_cpu_overall_usage_percent --latest
```

## 未来扩展

- 支持更多产品 (Everoute, etc.)
- 指标告警阈值配置
- 导出到 Prometheus 格式