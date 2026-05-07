package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newToolsInstall() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.tools.install [vm-name|vm-id]",
		Short: "Install VMware Tools on VM",
		Long: `Install VMware Tools on a virtual machine.

Examples:
  goct vm.tools.install myvm`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			ref, err := service.NewVM(cli).InstallVmtools(c.Context(), id)
			if err != nil {
				return err
			}
			if ref.ID == "already-installed" {
				fmt.Fprintln(c.OutOrStdout(), "VMware Tools already installed")
				return nil
			}
			if ref.IsSync() {
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	return c
}