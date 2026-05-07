package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newNicLs() *cobra.Command {
	var vmID string
	var idOnly bool
	c := &cobra.Command{
		Use:   "nic.ls",
		Short: "List VM NICs",
		Long: `List all network interfaces attached to a VM.
Use --vm to specify the VM by name or ID.

Examples:
  goct vm nic.ls --vm my-vm
  goct vm nic.ls vm-uuid
  goct vm nic.ls --vm my-vm --format json`,
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			if len(args) > 0 && vmID == "" {
				vmID = args[0]
			}
			if vmID == "" {
				return fmt.Errorf("VM name or ID required (use --vm or positional arg)")
			}
			items, err := service.NewVM(cli).ListNics(c.Context(), vmID)
			if err != nil {
				return err
			}
			if idOnly {
				for _, it := range items {
					_, _ = fmt.Fprintln(c.OutOrStdout(), it.ID)
				}
				return nil
			}
			out := make([]any, len(items))
			for i := range items {
				out[i] = items[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, output.VMNicListColumns)
		},
	}
	c.Flags().StringVar(&vmID, "vm", "", "VM name or ID")
	c.Flags().BoolVar(&idOnly, "id-only", false, "Output only IDs, one per line (for scripting)")
	return c
}