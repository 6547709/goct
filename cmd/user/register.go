package user
import "github.com/spf13/cobra"
const groupID = "user"
func Register(root *cobra.Command) { root.AddCommand(newLs(), newInfo(), newCreate(), newDestroy()) }
