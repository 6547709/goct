package cluster

import "github.com/spf13/cobra"

const groupID = "cluster"

func Register(root *cobra.Command) {
	root.AddCommand(newLs(), newInfo())
}
