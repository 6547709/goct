// Package cmd 定义 goct 的所有 CLI 命令。
// root.go 提供根命令与全局 flag 注册入口；
// 各资源命令在自己的子包通过 init() 挂载到 rootCmd。
package cmd

import "github.com/spf13/cobra"

// rootCmd 是 goct 的根命令，子命令通过各子包注册。
var rootCmd = &cobra.Command{
	Use:           "goct",
	Short:         "CloudTower CLI (govc-style)",
	Long:          "goct is a govc-style command-line client for SmartX CloudTower.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute 是 main.go 的唯一入口。
func Execute() error {
	return rootCmd.Execute()
}
