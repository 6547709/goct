// Package version 存放 goct 自身版本信息，值由 ldflags 在编译时注入。
package version

import (
	"fmt"
	"os/exec"
	"strings"
)

// Version 是当前 git tag 或 "dev"。
// 可通过 ldflags 覆盖：-ldflags "-X github.com/6547709/goct/internal/version.Version=v1.0.0"
var Version = "dev"

// Commit 是当前 git commit SHA（前 8 位）。
var Commit = ""

// Date 是编译时间（YYYY-MM-DD）。
var Date = ""

func init() {
	// 如果 Version 已被 ldflags 覆盖（非 "dev"），跳过自动检测
	if Version != "dev" {
		return
	}

	// 自动从 git tag 检测版本
	if tag := gitDescribe(); tag != "" {
		Version = tag
	}
}

// gitDescribe 执行 git describe --tags --always，返回当前版本 tag。
func gitDescribe() string {
	cmd := exec.Command("git", "describe", "--tags", "--always", "--abbrev=0")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

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
