// Package vm 包含 vm.* 系列命令。
package vm

import "github.com/spf13/cobra"

const groupID = "vm"

func Register(root *cobra.Command) {
	root.AddCommand(
		newLs(), newInfo(), newIp(),
		newCreate(), newClone(), newDestroy(),
		newUpdate(),
		newMigrate(), newMigrateAcrossCluster(), newMigrateAbort(), newExport(),
		newPowerOn(), newPowerOff(), newPowerReset(),
		newPowerSuspend(), newPowerResume(),
		newShutDown(),
		newRecycle(), newRecover(),
		// vm.disk.*
		newDiskLs(), newDiskAdd(), newDiskExpand(), newDiskRm(), newDiskUpdate(),
		// vm.cdrom.*
		newCdRomLs(), newCdRomAdd(), newCdRomEject(), newCdRomRm(), newCdRomToggle(),
		// vm.nic.*
		newNicLs(), newNicAdd(), newNicRm(), newNicUpdate(),
		// vm.gpu.*
		newGpuLs(), newGpuAdd(), newGpuRm(),
		// vm.vnc & vm.tools
		newVNC(), newToolsInstall(),
		// vm.snapshot.*
		newSnapshotLs(), newSnapshotCreate(),
		newSnapshotRevert(), newSnapshotRm(),
		// vm.rebuild & misc
		newRebuild(), newResetPassword(), newConvertToVM(),
	)
}
