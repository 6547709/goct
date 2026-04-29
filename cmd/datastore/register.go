package datastore

import "github.com/spf13/cobra"

const groupID = "datastore"

func Register(root *cobra.Command) {
	root.AddCommand(newLs(), newInfo(), newDiskLs())
}
