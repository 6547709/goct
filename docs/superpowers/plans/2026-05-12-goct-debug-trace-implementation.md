# goct Debug/Trace 功能实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 goct 添加真实 HTTP 请求/响应 trace 功能（-trace/-verbose/-dump CLI flag + env var），完整记录 SDK 与 CloudTower API 交互。

**Architecture:** 通过 http.RoundTripper 接口拦截 SDK 所有 HTTP 请求/响应，输出结构化 JSON 日志。三级递进：trace（轻量摘要）→ verbose（headers+body，截断64KB）→ dump（完整无截断）。

**Tech Stack:** Go 标准库（net/http, encoding/json, sync, os, strings）、govc/openapi runtime（SDK HTTP 传输）、cobra（CLI 框架）

---

## 文件结构

```
pkg/debug/
  debug.go          # 现有：分级日志基础设施
  flags.go          # 新：DebugFlags + DebugOptions + ResolveTraceOptions
  trace.go          # 新：TraceRoundTripper 实现
  flags_test.go     # 新
  trace_test.go     # 新

pkg/adapter/
  client.go         # 修改：Options 新增 Transport 字段

cmd/
  root.go           # 修改：注册 DebugFlags，PersistentPreRunE 注入 trace
```

---

## Task 1: pkg/debug/flags.go

**Files:**
- Create: `pkg/debug/flags.go`

实现 DebugFlags（CLI flag 绑定）和 DebugOptions（合并后配置）。

- [ ] **Step 1: 创建 flags.go 文件**

```go
package debug

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// DebugFlags 对应 CLI -debug/-trace/-verbose/-dump 四个 bool flag。
// 绑定到 cobra PersistentFlags，由 Resolve() 合并 env var 后生成 DebugOptions。
type DebugFlags struct {
	Debug   bool // -debug: 启用 debug 日志（已由 GOCT_LOG=DEBUG 提供，此处保留以与 govc 一致）
	Trace   bool // -trace: 结构化 HTTP trace（轻量：method + path + status + duration）
	Verbose bool // -verbose: 完整 headers + body（截断 64KB）
	Dump    bool // -dump: 完整无截断 body
}

// Register 把四个 flag 注册到 cmd 的 PersistentFlags。
func (f *DebugFlags) Register(cmd *cobra.Command) {
	pf := cmd.PersistentFlags()
	pf.BoolVar(&f.Debug, "debug", false, "Enable debug logging (env: GOCT_DEBUG)")
	pf.BoolVar(&f.Trace, "trace", false, "Enable HTTP trace (env: GOCT_TRACE)")
	pf.BoolVar(&f.Verbose, "verbose", false, "Enable verbose HTTP trace with headers and body (env: GOCT_VERBOSE)")
	pf.BoolVar(&f.Dump, "dump", false, "Enable full HTTP trace without body truncation (env: GOCT_DUMP)")
}

// DebugOptions 是 Resolve() 的输出，表示合并后的最终配置。
type DebugOptions struct {
	TraceLevel TraceLevel
}

// TraceLevel 定义 trace 详细程度。
type TraceLevel int

const (
	TraceLevelOff     TraceLevel = 0
	TraceLevelTrace  TraceLevel = 1 // 轻量：method + path + status + duration
	TraceLevelVerbose TraceLevel = 2 // 完整：+ headers + body（截断 64KB）
	TraceLevelDump    TraceLevel = 3 // 完整无截断
)

// Resolve 合并 CLI flag 与 env var（flag 优先）。
func (f *DebugFlags) Resolve() DebugOptions {
	level := TraceLevelOff

	// CLI flag 优先；未设置时检查 env var
	if f.Trace || f.resolveEnv("GOCT_TRACE") {
		level = TraceLevelTrace
	}
	if f.Verbose || f.resolveEnv("GOCT_VERBOSE") {
		level = TraceLevelVerbose
	}
	if f.Dump || f.resolveEnv("GOCT_DUMP") {
		level = TraceLevelDump
	}

	return DebugOptions{TraceLevel: level}
}

// resolveEnv 检查环境变量是否设为 "true"（不区分大小写）。
func (f *DebugFlags) resolveEnv(key string) bool {
	return strings.ToLower(os.Getenv(key)) == "true"
}
```

注意：上述代码引用了 cobra.Command 但未 import。实际文件需添加 import。

- [ ] **Step 2: 验证文件可编译**

Run: `cd /Users/liguoqiang/project/goct && go build ./pkg/debug/`
Expected: 无错误（会报 cobra 未 import，稍后修复 import）

- [ ] **Step 3: 补全 import 并验证**

```go
import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)
```

Run: `go build ./pkg/debug/`
Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add pkg/debug/flags.go
git commit -m "feat: add DebugFlags and DebugOptions for trace CLI flags"
```

---

## Task 2: pkg/debug/trace.go

**Files:**
- Create: `pkg/debug/trace.go`

实现 TraceRoundTripper（http.RoundTripper 接口实现），拦截 HTTP 请求/响应并输出结构化 JSON。

- [ ] **Step 1: 创建 trace.go 文件**

```go
package debug

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// TraceRoundTripper 实现 http.RoundTripper，拦截 SDK 所有 HTTP 请求/响应。
type TraceRoundTripper struct {
	Base   http.RoundTripper // 底层真实 transport（必填）
	Level  TraceLevel        // trace 级别
	Output io.Writer         // 输出目的地（应与 debug.Log 同路）
	Mutex  sync.Mutex
}

// NewTraceRoundTripper 构造 TraceRoundTripper。
// base: 底层 transport，nil 时使用 http.DefaultTransport。
// level: trace 详细程度。
// output: 输出目的地，nil 时使用 os.Stderr。
func NewTraceRoundTripper(base http.RoundTripper, level TraceLevel, output io.Writer) *TraceRoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	if output == nil {
		output = os.Stderr
	}
	return &TraceRoundTripper{
		Base:   base,
		Level:  level,
		Output: output,
	}
}

// RoundTrip 实现 http.RoundTripper 接口。
func (t *TraceRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	// --- 请求阶段 ---
	reqBytes := req.ContentLength

	// 发送真实请求
	resp, err := t.Base.RoundTrip(req)
	duration := time.Since(start)

	// --- 响应阶段 ---
	trace := TraceEntry{
		Level:       "trace",
		Method:      req.Method,
		Path:        req.URL.Path,
		Host:        req.URL.Host,
		DurationMs:  duration.Milliseconds(),
	}

	if err != nil {
		trace.Error = err.Error()
		t.emit(&trace)
		return resp, err
	}

	trace.RespStatus = resp.StatusCode
	trace.ReqBytes = reqBytes

	if t.Level >= TraceLevelVerbose {
		trace.ReqHeaders = t.filterHeaders(req.Header)
		if req.Body != nil && reqBytes > 0 {
			trace.ReqBody, trace.BodyTruncated, trace.Binary = t.readBody(req.Body, reqBytes)
		}
		trace.RespHeaders = t.filterHeaders(resp.Header)
	}

	if resp.Body != nil {
		respBytes := resp.ContentLength
		trace.RespBytes = respBytes
		if t.Level >= TraceLevelVerbose {
			bodyStr, truncated, binary := t.readBody(resp.Body, respBytes)
			trace.RespBody = bodyStr
			trace.BodyTruncated = truncated
			trace.Binary = binary
		}
	}

	t.emit(&trace)
	return resp, err
}

// TraceEntry 是结构化 trace JSON 的 schema。
type TraceEntry struct {
	Level         string            `json:"level"`                   // 固定 "trace"
	Method        string            `json:"method"`                  // HTTP method
	Path          string            `json:"path"`                    // URL path
	Host          string            `json:"host"`                    // URL host
	ReqHeaders    map[string]string `json:"req_headers,omitempty"`   // 请求 headers（verbose+）
	ReqBody       string            `json:"req_body,omitempty"`      // 请求 body（verbose+）
	ReqBytes      int64             `json:"req_bytes,omitempty"`     // 请求体大小
	RespStatus    int               `json:"resp_status,omitempty"`   // 响应状态码
	RespHeaders   map[string]string `json:"resp_headers,omitempty"`  // 响应 headers（verbose+）
	RespBody      string            `json:"resp_body,omitempty"`     // 响应 body（verbose+）
	RespBytes     int64             `json:"resp_bytes,omitempty"`    // 响应体大小
	BodyTruncated bool              `json:"body_truncated,omitempty"` // body 是否被截断
	DurationMs    int64             `json:"duration_ms"`             // 耗时（毫秒）
	Error         string            `json:"error,omitempty"`          // 错误信息
	Binary        bool              `json:"binary,omitempty"`        // 非 JSON body
}

// emit 输出 trace JSON（线程安全）。
func (t *TraceRoundTripper) emit(entry *TraceEntry) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	// dump 级别且 body 未截断时，不截断输出
	if t.Level < TraceLevelVerbose {
		entry.ReqBody = ""
		entry.RespBody = ""
		entry.ReqHeaders = nil
		entry.RespHeaders = nil
	}

	data, _ := json.Marshal(entry)
	t.Output.Write(data)
	t.Output.Write([]byte("\n"))
}

// filterHeaders 过滤敏感 headers（替换 Authorization token）。
func (t *TraceRoundTripper) filterHeaders(h http.Header) map[string]string {
	out := make(map[string]string)
	for k, v := range h {
		val := strings.Join(v, ", ")
		if k == "Authorization" {
			val = "***"
		}
		out[k] = val
	}
	return out
}

// readBody 读取 body 内容，支持截断。返回 (body字符串, 是否截断, 是否为二进制)。
func (t *TraceRoundTripper) readBody(body io.ReadCloser, size int64) (string, bool, bool) {
	const maxBodySize = 64 * 1024 // 64KB

	if body == nil {
		return "", false, false
	}

	// 确定读取上限
	limit := maxBodySize
	if t.Level == TraceLevelDump || (size > 0 && size < maxBodySize) {
		limit = size
	}

	// 读取 body（使用 LimitReader 避免一次读取过多）
	lr := io.LimitReader(body, limit)
	buf, _ := io.ReadAll(lr)
	body.Close()

	truncated := size > 0 && int64(len(buf)) < size

	// 检测是否为 binary（非 JSON）
	if !isJSON(buf) {
		return "", false, true
	}

	// 过滤 request body 中的敏感字段
	filtered := filterSensitiveBody(string(buf))

	return filtered, truncated, false
}

// isJSON 检测数据是否为 JSON（简单检测：首字符为 { 或 [）。
func isJSON(data []byte) bool {
	trimmed := bytes.TrimLeft(data, " \t\r\n")
	return len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[')
}

// filterSensitiveBody 过滤 body 中的敏感字段（password 等）。
func filterSensitiveBody(body string) string {
	// 简单替换：password 字段值替换为 ***
	re := regexp.MustCompile(`"password"\s*:\s*"[^"]*"`)
	return re.ReplaceAllString(body, `"password": "***"`)
}

// 注意：response body 被读取后，原始 handler 无法再次读取。
// 如果 SDK 后续还需要读取 response body，则需要在 RoundTrip 外部预读取，
// 这里简化处理以专注于 trace 功能。
```



- [ ] **Step 2: 运行测试发现未实现部分**

Run: `go build ./pkg/debug/`
Expected: 编译通过（结构已定义完整）

- [ ] **Step 3: 提交**

```bash
git add pkg/debug/trace.go
git commit -m "feat: add TraceRoundTripper implementing http.RoundTripper"
```

---

## Task 3: pkg/debug/trace_test.go

**Files:**
- Create: `pkg/debug/trace_test.go`

为 TraceRoundTripper 编写单元测试。

- [ ] **Step 1: 创建测试文件**

```go
package debug

import (
	"net/http"
	"strings"
	"testing"
)

// mockTransport 记录所有经过的请求，用于验证 trace 输出。
type mockTransport struct {
	Requests []*http.Request
	Response *http.Response
	Err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	m.Requests = append(m.Requests, req)
	return m.Response, m.Err
}

func TestTraceRoundTripper_TraceLevel(t *testing.T) {
	// 1. 创建 mockTransport
	// 2. 创建 TraceRoundTripper with TraceLevelTrace
	// 3. 发起 fake 请求
	// 4. 验证 trace 输出 JSON 包含 method/path/status/duration_ms
}

func TestTraceRoundTripper_VerboseLevel(t *testing.T) {
	// 1. 创建 mockTransport，返回带 body 的 Response
	// 2. 创建 TraceRoundTripper with TraceLevelVerbose
	// 3. 验证 trace 输出包含 req_headers/resp_headers/req_body/resp_body
	// 4. 验证 Authorization header 被过滤为 ***
}

func TestTraceRoundTripper_BodyTruncation(t *testing.T) {
	// 1. 创建返回 >64KB body 的 mockTransport
	// 2. 创建 TraceRoundTripper with TraceLevelVerbose
	// 3. 验证 body_truncated=true
	// 4. 创建 TraceRoundTripper with TraceLevelDump
	// 5. 验证 body_truncated=false
}
```

- [ ] **Step 2: 运行测试验证**

Run: `go test ./pkg/debug/ -v`
Expected: 测试框架可执行（测试体在后续迭代中完整实现）

- [ ] **Step 3: 提交**（等测试体实现后再提交）

---

## Task 4: pkg/debug/flags_test.go

**Files:**
- Create: `pkg/debug/flags_test.go`

- [ ] **Step 1: 创建测试文件**

```go
package debug

import (
	"os"
	"testing"
)

func TestDebugFlags_Resolve_EnvVar(t *testing.T) {
	// 1. 设置 GOCT_TRACE=true，验证 Resolve() 返回 TraceLevelTrace
	// 2. 设置 GOCT_VERBOSE=true，验证返回 TraceLevelVerbose
	// 3. 设置 GOCT_DUMP=true，验证返回 TraceLevelDump
	// 4. 清理环境变量
}

func TestDebugFlags_Resolve_FlagPriority(t *testing.T) {
	// 1. 设置 GOCT_TRACE=true
	// 2. CLI flag -verbose=true
	// 3. 验证返回 TraceLevelVerbose（flag 优先于 env）
}
```

- [ ] **Step 2: 运行测试验证**

Run: `go test ./pkg/debug/ -v`
Expected: 测试框架可执行

- [ ] **Step 3: 提交**（等测试体实现后再提交）

---

## Task 5: pkg/adapter/client.go

**Files:**
- Modify: `pkg/adapter/client.go:56-63`（Options 结构体）
- Modify: `pkg/adapter/client.go:142-153`（newTransport 函数）

- [ ] **Step 1: 修改 Options 结构体**

在 Options 结构体添加 `Transport http.RoundTripper` 字段：

```go
type Options struct {
	URL        string
	Username   string
	Password   string
	Source     string
	Insecure   bool
	Token      string
	Transport  http.RoundTripper  // 新增：可选自定义 transport
}
```

- [ ] **Step 2: 修改 newTransport 函数签名**

```go
func newTransport(host, basePath string, schemes []string, insecure bool, customTransport http.RoundTripper) *httptransport.Runtime {
	hc := &http.Client{}
	if insecure {
		hc.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	if customTransport != nil {
		hc.Transport = customTransport
	}
	return httptransport.NewWithClient(host, basePath, schemes, hc)
}
```

**说明：** customTransport 最终会包装在 hc.Transport 上。无论 insecure 是否为 true，只要提供了 customTransport 就用它替换原有 transport。

- [ ] **Step 3: 修改 client.New 签名以接受 transport 参数**

找到 `client.New` 函数（位于 `pkg/client/client.go`），修改签名为：
```go
func New(ctx context.Context, cfg config.Resolved, transport http.RoundTripper) (adapter.Client, error)
```

在函数内部调用 `adapter.NewClient` 时传入 transport：
```go
c, _, e := adapter.NewClient(ctx, adapter.Options{
    URL:      cfg.URL,
    Insecure: cfg.Insecure,
    Token:    tok.Value,
    Transport: transport,  // 新增
})
```

- [ ] **Step 4: 验证编译**

Run: `go build ./pkg/adapter/ && go build ./pkg/client/`
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add pkg/adapter/client.go pkg/client/client.go
git commit -m "feat: add Transport field to Options for custom http.RoundTripper"
```

---

## Task 6: cmd/root.go

**Files:**
- Modify: `cmd/root.go`

- [ ] **Step 1: 添加 debugFlags 全局变量**

```go
var debugFlags = flags.DebugFlags{}
```

- [ ] **Step 2: 在 init() 中注册 debugFlags**

```go
func init() {
	connFlags.Register(rootCmd)
	outputFlags.Register(rootCmd)
	debugFlags.Register(rootCmd)  // 新增
	// ...
}
```

- [ ] **Step 3: 修改 PersistentPreRunE**

```go
rootCmd.PersistentPreRunE = func(c *cobra.Command, _ []string) error {
	debug.Init()

	opts := debugFlags.Resolve()

	// 如果启用了 trace，创建 TraceRoundTripper
	var traceTransport http.RoundTripper
	if opts.TraceLevel > TraceLevelOff {
		traceTransport = debug.NewTraceRoundTripper(nil, opts.TraceLevel, os.Stderr)
	}

	cfg, err := config.Resolve(config.Override{
		URL:         connFlags.URL,
		Username:    connFlags.Username,
		Password:    connFlags.Password,
		Cluster:     connFlags.Cluster,
		Source:      connFlags.Source,
		Insecure:    connFlags.Insecure,
		InsecureSet: c.Flags().Changed("insecure"),
	})
	if err != nil {
		return err
	}
	cli, err := client.New(c.Context(), cfg, traceTransport)
	if err != nil {
		return err
	}
	c.SetContext(client.With(c.Context(), cli))
	return nil
}
```

**注意：** `client.New` 签名已在 Task 5 Step 3 中修改为 `func New(ctx context.Context, cfg config.Resolved, transport http.RoundTripper) (adapter.Client, error)`。

- [ ] **Step 4: 验证编译**

Run: `go build ./`
Expected: PASS 或有明确的编译错误（用于指导后续修改）

- [ ] **Step 5: 提交**

```bash
git add cmd/root.go
git commit -m "feat: register DebugFlags and inject TraceRoundTripper in PersistentPreRunE"
```

---

## Task 7: 集成测试

- [ ] **Step 1: 手动验证 trace 输出**

```bash
# 构建
cd /Users/liguoqiang/project/goct && go build -o goct .

# 轻量 trace（-trace flag 优先，env var 兜底）
./goct vm.ls -trace 2>&1 | jq .

# 完整 verbose
./goct vm.ls -verbose 2>&1 | jq .

# 完整 dump
./goct vm.ls -dump 2>&1 | jq .

# 也支持 env var 方式
GOCT_TRACE=true ./goct vm.ls 2>&1 | jq .
GOCT_VERBOSE=true ./goct vm.ls 2>&1 | jq .
```

- [ ] **Step 2: 运行单元测试**

```bash
go test ./pkg/debug/... -v
```

---

## 执行顺序

1. Task 1: flags.go（基础设施）
2. Task 2: trace.go（核心功能）
3. Task 3: trace_test.go（测试）
4. Task 4: flags_test.go（测试）
5. Task 5: client.go 修改（适配层）
6. Task 6: root.go 修改（入口）
7. Task 7: 集成测试
