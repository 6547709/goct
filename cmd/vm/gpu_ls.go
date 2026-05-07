package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newGpuLs() *cobra.Command {
	var vmID string
	c := &cobra.Command{
		Use:   "gpu.ls",
		Short: "List VM GPU devices",
		Long: `List all GPU devices attached to a VM.
Use --vm to specify the VM by name or ID.

Examples:
  goct vm gpu.ls --vm my-vm
  goct vm gpu.ls vm-uuid`,
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			if len(args) > 0 && vmID == "" {
				vmID = args[0]
			}
			if vmID == "" {
				return fmt.Errorf("VM name or ID required (use --vm or positional arg)")
			}
			v, err := service.NewVM(cli).Resolve(c.Context(), vmID)
			if err != nil {
				return err
			}
			out := make([]any, len(v.GpuDevices))
			for i := range v.GpuDevices {
				out[i] = v.GpuDevices[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, output.GpuDeviceListColumns)
		},
	}
	c.Flags().StringVar(&vmID, "vm", "", "VM name or ID")
	return c
}