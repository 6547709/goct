package cluster_settings

import "github.com/spf13/cobra"

const groupID = "system"

func Register(root *cobra.Command) {
	root.AddCommand(newGet())
}
