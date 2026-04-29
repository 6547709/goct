package vlan
import "github.com/spf13/cobra"
const groupID = "network"
func Register(root *cobra.Command) { root.AddCommand(newLs(), newInfo(), newCreate(), newDestroy()) }
