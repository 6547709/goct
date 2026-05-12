package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newDiskLs() *cobra.Command {
	var vmID string
	var idOnly bool
	c := &cobra.Command{
		Use:     "vm.disk.ls",
		Short:   "List VM disks",
		GroupID: "vm",
		Long: `List all disks (including CD-ROMs) attached to a VM.
Use --vm to specify the VM by name or ID.

Examples:
  goct vm disk.ls --vm my-vm
  goct vm disk.ls vm-uuid
  goct vm disk.ls --vm my-vm --format json`,
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			if len(args) > 0 && vmID == "" {
				vmID = args[0]
			}
			if vmID == "" {
				return fmt.Errorf("VM name or ID required (use --vm or positional arg)")
			}
			items, err := service.NewVM(cli).ListDisks(c.Context(), vmID)
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
			return output.Render(c.OutOrStdout(), out, format, output.VMDiskListColumns)
		},
	}
	c.Flags().StringVar(&vmID, "vm", "", "VM name or ID")
	c.Flags().BoolVar(&idOnly, "id-only", false, "Output only IDs, one per line (for scripting)")
	return c
}

func newDiskUpdate() *cobra.Command {
	var diskID string
	c := &cobra.Command{
		Use:     "vm.disk.update",
		Short:   "Update VM disk settings",
		GroupID: "vm",
		Long: `Update a VM disk configuration.

Examples:
  goct vm disk.update --disk disk-uuid
  goct vm disk.update disk-uuid`,
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			if len(args) > 0 && diskID == "" {
				diskID = args[0]
			}
			if diskID == "" {
				return fmt.Errorf("disk ID required (use --disk or positional arg)")
			}
			// Note: MaxBandwidth and MaxIops not supported by update-vm-disk API
			ref, err := service.NewVM(cli).UpdateDisk(c.Context(), diskID, adapter.DiskUpdateSpec{})
			if err != nil {
				return err
			}
			if ref.IsSync() {
				_, _ = fmt.Fprintln(c.OutOrStdout(), "disk updated")
			} else {
				_, _ = fmt.Fprintf(c.OutOrStdout(), "task: %s\n", ref.ID)
			}
			return nil
		},
	}
	c.Flags().StringVar(&diskID, "disk", "", "Disk ID")
	return c
}