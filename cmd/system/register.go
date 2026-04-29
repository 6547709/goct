// Package system 包含 goct 系统类子命令：about / version / session.*。
package system

import "github.com/spf13/cobra"

const groupID = "system"

// Register 把所有 system 子命令挂到 rootCmd。
func Register(root *cobra.Command) {
	root.AddCommand(newAbout(), newVersion(),
		newSessionLogin(), newSessionLogout(), newSessionLs())
}
