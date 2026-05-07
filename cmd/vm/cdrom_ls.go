package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newCdRomLs() *cobra.Command {
	var vmID string
	c := &cobra.Command{
		Use:   "cdrom.ls",
		Short: "List VM CD-ROMs",
		Long: `List all CD-ROM devices attached to a VM.
Use --vm to specify the VM by name or ID.

Examples:
  goct vm cdrom.ls --vm my-vm
  goct vm cdrom.ls vm-uuid`,
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			if len(args) > 0 && vmID == "" {
				vmID = args[0]
			}
			if vmID == "" {
				return fmt.Errorf("VM name or ID required (use --vm or positional arg)")
			}
			disks, err := service.NewVM(cli).ListDisks(c.Context(), vmID)
			if err != nil {
				return err
			}
			// Filter to only CD-ROMs
			var cdroms []adapter.VMDisk
			for _, d := range disks {
				if d.Type == "CD_ROM" {
					cdroms = append(cdroms, d)
				}
			}
			out := make([]any, len(cdroms))
			for i := range cdroms {
				out[i] = cdroms[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, output.VMDiskListColumns)
		},
	}
	c.Flags().StringVar(&vmID, "vm", "", "VM name or ID")
	return c
}

func newCdRomToggle() *cobra.Command {
	var disabled bool
	c := &cobra.Command{
		Use:   "cdrom.toggle [cdrom-id]",
		Short: "Enable or disable a CD-ROM",
		Long: `Enable or disable a CD-ROM device.

Examples:
  goct vm cdrom.toggle cdrom-uuid
  goct vm cdrom.toggle cdrom-uuid --disabled`,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewVM(cli).ToggleCdRom(c.Context(), args[0], disabled)
			if err != nil {
				return err
			}
			if ref.IsSync() {
				_, _ = fmt.Fprintln(c.OutOrStdout(), "cdrom toggled")
			} else {
				_, _ = fmt.Fprintf(c.OutOrStdout(), "task: %s\n", ref.ID)
			}
			return nil
		},
	}
	c.Flags().BoolVar(&disabled, "disabled", false, "Disable the CD-ROM (default: enable)")
	return c
}