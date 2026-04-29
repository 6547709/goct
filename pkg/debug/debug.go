// Package debug 提供 goct 的分级日志功能（Packer 风格）。
//
// 通过环境变量 GOCT_LOG 控制：
//
//	GOCT_LOG=TRACE   最详细：SDK 请求/响应细节
//	GOCT_LOG=DEBUG   调试：config 解析、session 命中/miss、adapter 调用
//	GOCT_LOG=INFO    一般：登录成功、命令执行
//	GOCT_LOG=WARN    警告：session 过期、token 刷新
//	GOCT_LOG=ERROR   仅错误
//	GOCT_LOG=（空/未设）关闭所有日志
//
// 可选 GOCT_LOG_FILE 把日志写入文件（默认 stderr）。
// 所有日志走 stderr / 文件，不污染 stdout 业务输出。
package debug

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Level 定义日志级别。
type Level int

const (
	LevelOff   Level = 0
	LevelError Level = 1
	LevelWarn  Level = 2
	LevelInfo  Level = 3
	LevelDebug Level = 4
	LevelTrace Level = 5
)

var (
	currentLevel Level     = LevelOff
	out          io.Writer = io.Discard
	logFile      *os.File
)

// Init 根据 GOCT_LOG / GOCT_LOG_FILE 初始化日志。
// 应在 main 或 PersistentPreRunE 最早期调用一次。
func Init() {
	raw := strings.ToUpper(strings.TrimSpace(os.Getenv("GOCT_LOG")))
	if raw == "" {
		return
	}
	currentLevel = parseLevel(raw)
	if currentLevel == LevelOff {
		return
	}

	out = os.Stderr
	if path := os.Getenv("GOCT_LOG_FILE"); path != "" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			out = f
			logFile = f
		}
	}

	Infof("logging enabled: level=%s", levelName(currentLevel))
}

// Close 关闭日志文件（如果打开过）。在程序退出前调用。
func Close() {
	if logFile != nil {
		_ = logFile.Close()
	}
}

// Enabled 报告日志是否开启。
func Enabled() bool { return currentLevel > LevelOff }

// IsLevel 报告指定级别是否会被输出。
func IsLevel(l Level) bool { return currentLevel >= l }

// --- 各级别日志函数 ---

func Trace(args ...any)                 { log(LevelTrace, args...) }
func Tracef(format string, args ...any) { logf(LevelTrace, format, args...) }

func Debug(args ...any)                 { log(LevelDebug, args...) }
func Debugf(format string, args ...any) { logf(LevelDebug, format, args...) }

func Info(args ...any)                  { log(LevelInfo, args...) }
func Infof(format string, args ...any)  { logf(LevelInfo, format, args...) }

func Warn(args ...any)                  { log(LevelWarn, args...) }
func Warnf(format string, args ...any)  { logf(LevelWarn, format, args...) }

func Error(args ...any)                 { log(LevelError, args...) }
func Errorf(format string, args ...any) { logf(LevelError, format, args...) }

// --- 内部实现 ---

func log(level Level, args ...any) {
	if currentLevel < level {
		return
	}
	fmt.Fprintf(out, "%s [%-5s] %s\n", timestamp(), levelName(level), fmt.Sprint(args...))
}

func logf(level Level, format string, args ...any) {
	if currentLevel < level {
		return
	}
	fmt.Fprintf(out, "%s [%-5s] %s\n", timestamp(), levelName(level), fmt.Sprintf(format, args...))
}

func timestamp() string {
	return time.Now().Format("2006-01-02T15:04:05.000")
}

func parseLevel(s string) Level {
	switch s {
	case "TRACE":
		return LevelTrace
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	default:
		return LevelOff
	}
}

func levelName(l Level) string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "OFF"
	}
}
