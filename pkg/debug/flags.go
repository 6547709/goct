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