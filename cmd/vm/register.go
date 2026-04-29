// Package vm 包含 vm.* 系列命令。
package vm

import "github.com/spf13/cobra"

const groupID = "vm"

// Register 把所有 vm + vm.snapshot 子命令挂到 rootCmd。
func Register(root *cobra.Command) {
	root.AddCommand(
		newLs(), newInfo(),
		newCreate(), newClone(), newDestroy(),
		newMigrate(), newExport(),
		newPowerOn(), newPowerOff(), newPowerReset(),
		newPowerSuspend(), newPowerResume(),
		// vm.snapshot.*
		newSnapshotLs(), newSnapshotCreate(),
		newSnapshotRevert(), newSnapshotRm(),
	)
}
