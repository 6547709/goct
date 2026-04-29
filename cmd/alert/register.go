package alert
import "github.com/spf13/cobra"
const groupID = "alert"
func Register(root *cobra.Command) { root.AddCommand(newLs(), newInfo(), newAck()) }
