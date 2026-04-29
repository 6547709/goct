package host

import "github.com/spf13/cobra"

const groupID = "host"

func Register(root *cobra.Command) {
	root.AddCommand(
		newLs(), newInfo(),
		newMaintenanceEnter(), newMaintenanceExit(),
		newShutdown(), newReboot(),
		newReconnect(), newDisconnect(),
	)
}
