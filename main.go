// Command goct 是 SmartX CloudTower 的 govc 风格命令行工具。
// 程序入口仅负责调用 cmd.Execute() 并处理顶层错误。
package main

import (
	"fmt"
	"os"

	"github.com/6547709/goct/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
