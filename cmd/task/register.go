package task
import "github.com/spf13/cobra"
const groupID = "task"
func Register(root *cobra.Command) { root.AddCommand(newLs(), newInfo(), newCancel(), newWait()) }
