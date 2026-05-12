# goct Debug/Trace 功能设计

> **状态:** 已批准

## 目标

为 goct 添加真实的 HTTP 请求/响应 trace 功能，完整记录 SDK 与 CloudTower API 的交互，便于调试和问题定位。

## 背景

goct 是 SmartX CloudTower 的 govc 风格 CLI 工具。当前 `pkg/debug/debug.go` 仅提供分级日志功能（GOCT_LOG env var），缺乏真实的 HTTP 观测能力。用户需要像 govc 那样的 `-debug/-trace/-verbose/-dump` CLI flag，以及完整的请求/响应追踪。

## 架构概览

```
用户执行命令
    │
    ▼
cmd/root.go PersistentPreRunE
    │ 读取 DebugFlags（CLI flag 值）
    │ 合并 env var（flag 优先）
    ▼
adapter.NewClient(ctx, opts)
    │ 传入 trace options
    ▼
newTransport()
    │ 创建 http.Client
    │ 包装 RoundTripper（如果 trace 开启）
    ▼
httptransport.Runtime (SDK 内部)
    │ 使用包装后的 http.Client
    ▼
RoundTripper.RoundTrip()
    │ 拦截每个 HTTP 请求/响应
    │ 输出结构化 JSON trace
    ▼
真实 HTTP 请求发出
```

## 组件设计

### 1. pkg/debug/flags.go（新文件）

**DebugFlags 结构体：**

```go
type DebugFlags struct {
    Debug   bool   // -debug: 启用 debug 日志
    Trace   bool   // -trace: 结构化 HTTP trace（轻量：method + path + status + duration）
    Verbose bool   // -verbose: 完整 headers + body（截断 64KB）
    Dump    bool   // -dump: 完整无截断 body
}

func (f *DebugFlags) Register(c *cobra.Command)
func (f *DebugFlags) Resolve() DebugOptions  // 合并 env var，返回最终配置
```

**对应环境变量：**

| CLI Flag | Env Var | 说明 |
|----------|---------|------|
| -debug | GOCT_DEBUG | 启用 debug 日志 |
| -trace | GOCT_TRACE | 启用 trace（轻量） |
| -verbose | GOCT_VERBOSE | 启用 trace（完整） |
| -dump | GOCT_DUMP | 启用 trace（无截断） |

**合并规则：** CLI flag 优先于 env var；env var 值为 "true"（不区分大小写）时视为启用。

### 2. pkg/debug/trace.go（新文件）

**TraceRoundTripper 实现 `http.RoundTripper` 接口：**

```go
type TraceRoundTripper struct {
    Base   http.RoundTripper  // 底层真实 transport
    Level  TraceLevel         // trace / verbose / dump
    Output io.Writer          // 输出目的地（与 debug.Log 同路）
    Mutex  sync.Mutex
}

type TraceLevel int
const (
    TraceLevelOff     TraceLevel = 0
    TraceLevelTrace   TraceLevel = 1  // 轻量：method + path + status + duration
    TraceLevelVerbose TraceLevel = 2  // 完整：+ headers + body（截断 64KB）
    TraceLevelDump    TraceLevel = 3  // 完整无截断
)
```

**Trace JSON 结构：**

```json
{
  "level": "trace",
  "method": "POST",
  "path": "/v2/api/vm/getVms",
  "host": "10.0.50.210",
  "req_headers": {"Content-Type": "application/json", "Authorization": "***"},
  "req_bytes": 120,
  "resp_status": 200,
  "resp_headers": {"Content-Type": "application/json"},
  "resp_bytes": 456,
  "duration_ms": 12
}
```

**verbose/dump 级别额外字段：**
```json
{
  "req_body": "{\"where\":{\"name_contains\":\"demo01\"},...}",
  "resp_body": "{\"data\":[{\"id\":\"...\",\"name\":\"demo01\"}]}",
  "body_truncated": false
}
```

**边界处理：**

1. **Body 截断：** verbose 级别 body > 64KB 时截断，设置 `"body_truncated": true`
2. **敏感信息过滤：** Authorization header token 值用 `***` 替换；request body 中 password 字段过滤
3. **非 JSON body：** 标注 `"binary": true`，不输出 body 内容
4. **网络错误：** `"error": "connection refused"`，无 resp_status/resp_body
5. **超时：** `"error": "context deadline exceeded"`，附 duration_ms

### 3. pkg/adapter/client.go（修改）

**Options 新增字段：**

```go
type Options struct {
    // 现有字段...
    Transport http.RoundTripper  // 可选：注入自定义 RoundTripper（用于 trace）
}
```

**newTransport 改动：**

```go
func newTransport(host, basePath string, schemes []string, insecure bool, customTransport http.RoundTripper) *httptransport.Runtime {
    // 现有逻辑...
    if customTransport != nil {
        hc = &http.Client{Transport: customTransport}
    }
    return httptransport.NewWithClient(host, basePath, schemes, hc)
}
```

### 4. cmd/root.go（修改）

**PersistentPreRunE 改动：**

```go
var debugFlags flags.DebugFlags

func init() {
    debugFlags.Register(rootCmd)
}

rootCmd.PersistentPreRunE = func(c *cobra.Command, _ []string) error {
    debug.Init()  // 现有
    debugOpts := debugFlags.Resolve()

    // 如果启用了 trace/verbose/dump，创建 TraceRoundTripper
    if debugOpts.Trace || debugOpts.Verbose || debugOpts.Dump {
        traceOpts := debug.ResolveTraceOptions(debugOpts)
        traceTransport := trace.NewTraceRoundTripper(nil, traceOpts)
        // 注入到 adapter.Options
    }
    // ...
}
```

## 三层递进关系

| 级别 | Flag | 信息量 | JSON 字段 |
|------|------|--------|-----------|
| Off | - | 无 trace | - |
| Trace | -trace | 轻量摘要 | method, path, host, status, duration_ms, error |
| Verbose | -verbose | 完整 headers + body | + req_headers, req_body, resp_status, resp_headers, resp_body, body_truncated |
| Dump | -dump | 完整无截断 | 同 verbose，但 body 不截断 |

## 数据流

```
1. cmd/root.go 解析 CLI flags（-trace/-verbose/-dump）
2. 与 env var 合并（flag 优先）
3. 如果任一 trace 级别启用，创建 TraceRoundTripper
4. adapter.NewClient 接收自定义 Transport
5. newTransport() 用自定义 http.Client 包装 RoundTripper
6. SDK 所有 HTTP 请求经过 TraceRoundTripper.RoundTrip()
7. 输出结构化 JSON 到 debug 日志同一目的地（stderr 或 GOCT_LOG_FILE）
```

## 错误处理

| 场景 | JSON error 字段 | 其他字段 |
|------|-----------------|---------|
| 网络错误 | "connection refused" | 无 resp_status/resp_body |
| 超时 | "context deadline exceeded" | duration_ms |
| 非 JSON body | - | "binary": true |
| Body 过大（verbose） | - | "body_truncated": true |

## 并发安全

- `http.RoundTripper` 可能被多个 goroutine 并发调用
- `TraceRoundTripper` 使用 `sync.Mutex` 保护写入
- JSON 日志行一次性写入（避免多 goroutine 输出交叉）

## 测试策略

### 单元测试

1. **pkg/debug/flags_test.go**
   - 测试 CLI flag 解析
   - 测试 env var 合并逻辑
   - 测试 flag 优先于 env var

2. **pkg/debug/trace_test.go**
   - 测试 JSON 序列化格式
   - 测试 body 截断逻辑（64KB 阈值）
   - 测试敏感信息过滤（Authorization、password）
   - 测试边界情况（网络错误、超时、非 JSON body）

### 集成测试

- 使用 `fakeClient` 架构
- 构造 MockRoundTripper，记录所有 trace 输出
- 验证 trace/verbose/dump 三级输出内容差异

### 手动验证

```bash
# 轻量 trace
GOCT_LOG=TRACE ./goct vm.ls -trace 2>&1 | jq .

# 完整 verbose
GOCT_LOG=TRACE ./goct vm.ls -verbose 2>&1 | jq .

# 完整 dump（无截断）
GOCT_LOG=TRACE ./goct vm.ls -dump 2>&1 | jq .
```

## 文件清单

| 文件 | 操作 | 说明 |
|------|------|------|
| `pkg/debug/flags.go` | 创建 | DebugFlags 结构体 |
| `pkg/debug/trace.go` | 创建 | TraceRoundTripper 实现 |
| `pkg/debug/trace_test.go` | 创建 | trace 单元测试 |
| `pkg/debug/flags_test.go` | 创建 | flags 单元测试 |
| `pkg/adapter/client.go` | 修改 | Options 新增 Transport 字段 |
| `cmd/root.go` | 修改 | 注册 DebugFlags |

## 依赖关系

- 标准库：`net/http`, `encoding/json`, `sync`
- 现有包：`pkg/debug`（日志基础设施）
- 第三方：`go-openapi/runtime/client`（SDK HTTP 传输）
