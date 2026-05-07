package vm

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newShutDown() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.shutdown [name|id]",
		Short: "Gracefully shut down a VM (guest OS shutdown)",
		Long: `Send a graceful shutdown request to the VM's guest OS.
Unlike vm.power.off which cuts power, this signals the OS to shut down cleanly.
Requires VMtools to be installed and running in the guest.`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			ref, err := service.NewVM(cli).ShutDown(c.Context(), id)
			if err != nil {
				return err
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
