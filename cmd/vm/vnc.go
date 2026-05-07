package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newVNC() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.vnc [vm-name|vm-id]",
		Short: "Get VM VNC connection info",
		Long: `Get VNC connection information for a virtual machine.

Examples:
  goct vm.vnc myvm`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			info, err := service.NewVM(cli).GetVNCInfo(c.Context(), id)
			if err != nil {
				return err
			}
			if info.ClusterIP != "" {
				fmt.Fprintf(c.OutOrStdout(), "ClusterIP: %s\n", info.ClusterIP)
			}
			if info.Redirect != "" {
				fmt.Fprintf(c.OutOrStdout(), "Redirect: %s\n", info.Redirect)
			}
			if info.Terminal != "" {
				fmt.Fprintf(c.OutOrStdout(), "Terminal: %s\n", info.Terminal)
			}
			if info.Direct != "" {
				fmt.Fprintf(c.OutOrStdout(), "Direct: %s\n", info.Direct)
			}
			return nil
		},
	}
	return c
}