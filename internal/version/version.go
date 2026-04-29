// Package version 存放 goct 自身版本信息，值由 ldflags 在编译时注入。
package version

import "fmt"

// Version 是当前 git tag 或 "dev"。
var Version = "dev"

// Commit 是当前 git commit SHA（前 8 位）。
var Commit = ""

// Date 是编译时间（YYYY-MM-DD）。
var Date = ""

// String 返回格式化版本字符串。
func String() string {
	s := "goct " + Version
	if Commit != "" {
		s += " (" + Commit + ")"
	}
	if Date != "" {
		s += ", built " + Date
	}
	return s
}

// Full 返回包含全部信息的完整版本字符串（用于 goct about 输出头）。
func Full() string {
	if Version == "dev" && Commit == "" {
		return "goct " + Version + " (no build info)"
	}
	return String()
}

// Sprint 是 fmt.Sprintf 的简写别名。
func Sprint(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}
