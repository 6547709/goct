package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newIp() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.ip [name|id]",
		Short: "Show VM IP addresses",
		Long: `Output the IP address(es) of a VM, one per line.
Designed for scripting: pipe VM list IDs to get IPs.

Examples:
  goct vm.ip myvm
  goct vm.ls --name web --id-only | xargs -I{} goct vm.ip {}`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			v, err := service.NewVM(cli).Resolve(c.Context(), id)
			if err != nil {
				return err
			}
			if len(v.IPs) == 0 {
				return nil
			}
			for _, ip := range v.IPs {
				fmt.Fprintln(c.OutOrStdout(), ip)
			}
			return nil
		},
	}
	return c
}
